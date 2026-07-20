package pseo

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	indexNowKey     string
	indexNowKeyFile = "public/indexnow-key.txt"
	submittedURLs   = make(map[string]bool)
	submittedMu     sync.RWMutex
)

func InitIndexNow() {
	_ = os.MkdirAll("public", 0o755)

	data, err := os.ReadFile(indexNowKeyFile)
	if err != nil {
		key := generateIndexNowKey()
		if err := os.WriteFile(indexNowKeyFile, []byte(key), 0o644); err != nil {
			log.Printf("IndexNow: failed to write key file: %v", err)
			return
		}
		indexNowKey = key
		log.Printf("IndexNow: new key generated at %s", indexNowKeyFile)
	} else {
		indexNowKey = strings.TrimSpace(string(data))
	}
}

func generateIndexNowKey() string {
	h := sha256.Sum256([]byte(fmt.Sprintf("%s-%d", AppName, time.Now().UnixNano())))
	return hex.EncodeToString(h[:])[:32]
}

func GetIndexNowKey() string {
	return indexNowKey
}

func SubmitURL(url string) {
	submittedMu.RLock()
	if submittedURLs[url] {
		submittedMu.RUnlock()
		return
	}
	submittedMu.RUnlock()

	submittedMu.Lock()
	submittedURLs[url] = true
	submittedMu.Unlock()

	go submitToIndexNow(url)
}

func SubmitURLs(urls []string) {
	for _, url := range urls {
		go SubmitURL(url)
	}
}

func submitToIndexNow(url string) {
	if indexNowKey == "" {
		return
	}

	payload := map[string]interface{}{
		"host":        strings.TrimPrefix(strings.TrimPrefix(AppURL, "https://"), "http://"),
		"key":         indexNowKey,
		"keyLocation": AppURL + "/indexnow-key.txt",
		"urlList":     []string{url},
	}

	body, _ := json.Marshal(payload)

	endpoints := []string{
		"https://api.indexnow.org/indexnow",
		"https://www.bing.com/indexnow",
		"https://search.seznam.cz/indexnow",
		"https://search.naver.com/indexnow",
		"https://yandex.com/indexnow",
	}

	client := &http.Client{Timeout: 10 * time.Second}

	for _, endpoint := range endpoints {
		resp, err := client.Post(endpoint, "application/json", bytes.NewReader(body))
		if err != nil {
			log.Printf("IndexNow: failed to submit to %s: %v", endpoint, err)
			continue
		}
		resp.Body.Close()
		if resp.StatusCode == 200 || resp.StatusCode == 202 {
			log.Printf("IndexNow: submitted %s to %s", url, endpoint)
		} else {
			log.Printf("IndexNow: %s returned %d for %s", endpoint, resp.StatusCode, url)
		}
	}
}

func SubmitAllPSEOURLs() {
	urls := GetAllPSEOURLs()
	for _, url := range urls {
		go func(u string) {
			time.Sleep(100 * time.Millisecond)
			submitToIndexNow(u)
		}(url)
	}
	log.Printf("IndexNow: queued %d PSEO URLs for submission", len(urls))
}

func GetAllPSEOURLs() []string {
	var urls []string

	urls = append(urls, AppURL+"/")
	urls = append(urls, AppURL+"/docs")

	for _, bp := range BestPages {
		urls = append(urls, AppURL+"/best-"+bp.Slug)
		urls = append(urls, fmt.Sprintf("%s/best-%s-%d", AppURL, bp.Slug, time.Now().Year()))
	}

	for _, c := range Competitors {
		urls = append(urls, AppURL+"/alternatives-to-"+c)
	}

	for i := 0; i < len(Competitors); i++ {
		for j := i + 1; j < len(Competitors) && j < i+3; j++ {
			urls = append(urls, fmt.Sprintf("%s/compare/%s-vs-%s", AppURL, safeSlug(AppName), Competitors[j]))
		}
	}

	for _, ind := range Industries {
		urls = append(urls, AppURL+"/whatsapp-marketing-untuk-"+ind.Slug)
	}

	for _, sc := range SourceCodePages {
		urls = append(urls, AppURL+"/"+sc.Slug)
	}

	for _, city := range Cities {
		urls = append(urls, AppURL+"/jasa-whatsapp-marketing-"+city.Slug)
		urls = append(urls, AppURL+"/aplikasi-whatsapp-"+city.Slug)
	}

	for _, cp := range CaraPages {
		urls = append(urls, AppURL+"/cara-"+cp.Slug)
	}

	for _, cp := range ChatbotPages {
		urls = append(urls, AppURL+"/chatbot-"+cp.Slug)
	}

	for _, pp := range PanduanPages {
		urls = append(urls, AppURL+"/panduan-"+pp.Slug)
	}

	for _, feat := range Features {
		urls = append(urls, AppURL+"/best-whatsapp-"+feat.Slug)
	}

	for _, ind := range Industries {
		for _, city := range Cities {
			urls = append(urls, fmt.Sprintf("%s/whatsapp-marketing-untuk-%s-di-%s", AppURL, ind.Slug, city.Slug))
		}
	}

	return urls
}
