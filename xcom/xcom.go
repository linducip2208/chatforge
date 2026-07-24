//go:build pro

package xcom

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	APIKey       string
	APISecret    string
	AccessToken  string
	AccessSecret string
	UserID       string
}

func New(apiKey, apiSecret, accessToken, accessSecret, userID string) *Client {
	return &Client{APIKey: apiKey, APISecret: apiSecret, AccessToken: accessToken, AccessSecret: accessSecret, UserID: userID}
}

func (c *Client) nonce() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:32]
}

func (c *Client) oauth1(method, apiURL string, body url.Values) string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := c.nonce()
	all := url.Values{}
	all.Set("oauth_consumer_key", c.APIKey)
	all.Set("oauth_nonce", nonce)
	all.Set("oauth_signature_method", "HMAC-SHA1")
	all.Set("oauth_timestamp", timestamp)
	all.Set("oauth_token", c.AccessToken)
	all.Set("oauth_version", "1.0")
	for k, v := range body {
		all[k] = v
	}
	var keys []string
	for k := range all { keys = append(keys, k) }
	sort.Strings(keys)
	var paramStr string
	for i, k := range keys {
		if i > 0 { paramStr += "&" }
		paramStr += url.QueryEscape(k) + "=" + url.QueryEscape(strings.Join(all[k], ","))
	}
	sigBase := strings.ToUpper(method) + "&" + url.QueryEscape(apiURL) + "&" + url.QueryEscape(paramStr)
	signingKey := url.QueryEscape(c.APISecret) + "&" + url.QueryEscape(c.AccessSecret)
	mac := hmac.New(sha1.New, []byte(signingKey))
	mac.Write([]byte(sigBase))
	sig := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	all.Set("oauth_signature", sig)

	var hdr string
	for i, k := range keys {
		if i > 0 { hdr += ", " }
		hdr += k + `="` + url.QueryEscape(all.Get(k)) + `"`
	}
	return "OAuth " + hdr
}

type DMEvent struct {
	SenderID   string `json:"sender_id"`
	Text       string `json:"text"`
	CreatedAt  string `json:"created_at"`
}

func (c *Client) SendDM(recipientID, message string) error {
	apiURL := "https://api.twitter.com/1.1/direct_messages/events/new.json"
	body := map[string]interface{}{
		"event": map[string]interface{}{
			"type": "message_create",
			"message_create": map[string]interface{}{
				"target":       map[string]string{"recipient_id": recipientID},
				"message_data": map[string]string{"text": message},
			},
		},
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", apiURL, bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.oauth1("POST", apiURL, url.Values{}))
	resp, err := http.DefaultClient.Do(req)
	if err != nil { return err }
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("twitter api error %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}

func ParseWebhook(r *http.Request, consumerSecret string) ([]DMEvent, error) {
	body, _ := io.ReadAll(r.Body)
	token := r.Header.Get("X-Twitter-Webhooks-Signature")
	if token != "" {
		mac := hmac.New(sha256.New, []byte(consumerSecret))
		mac.Write(body)
		expected := "sha256=" + base64.StdEncoding.EncodeToString(mac.Sum(nil))
		if !hmac.Equal([]byte(token), []byte(expected)) {
			return nil, fmt.Errorf("invalid signature")
		}
	}
	var events []DMEvent
	var wrapper struct {
		DirectMessageEvents []struct {
			Type          string `json:"type"`
			MessageCreate struct {
				SenderID    string `json:"sender_id"`
				MessageData struct {
					Text string `json:"text"`
				} `json:"message_data"`
			} `json:"message_create"`
			CreatedAt string `json:"created_timestamp"`
		} `json:"direct_message_events"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, err
	}
	for _, evt := range wrapper.DirectMessageEvents {
		if evt.Type == "message_create" {
			events = append(events, DMEvent{
				SenderID: evt.MessageCreate.SenderID, Text: evt.MessageCreate.MessageData.Text, CreatedAt: evt.CreatedAt,
			})
		}
	}
	return events, nil
}

func CRCResponse(consumerSecret, crcToken string) string {
	mac := hmac.New(sha256.New, []byte(consumerSecret))
	mac.Write([]byte(crcToken))
	return "sha256=" + base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
