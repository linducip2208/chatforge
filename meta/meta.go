package meta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Client struct {
	PhoneNumberID string
	AccessToken   string
	VerifyToken   string
	HTTP          *http.Client
}

func New(phoneNumberID, accessToken, verifyToken string) *Client {
	return &Client{
		PhoneNumberID: phoneNumberID,
		AccessToken:   accessToken,
		VerifyToken:   verifyToken,
		HTTP:          &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) baseURL() string {
	return fmt.Sprintf("https://graph.facebook.com/v22.0/%s", c.PhoneNumberID)
}

func (c *Client) SendText(to, message string) (string, error) {
	body := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                to,
		"type":              "text",
		"text":              map[string]string{"body": message, "preview_url": "false"},
	}
	return c.doPost("/messages", body)
}

func (c *Client) SendImage(to, imageURL, caption string) (string, error) {
	body := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                to,
		"type":              "image",
		"image":             map[string]string{"link": imageURL, "caption": caption},
	}
	return c.doPost("/messages", body)
}

func (c *Client) SendDocument(to, docURL, filename, caption string) (string, error) {
	body := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                to,
		"type":              "document",
		"document":          map[string]string{"link": docURL, "filename": filename, "caption": caption},
	}
	return c.doPost("/messages", body)
}

func (c *Client) SendVideo(to, videoURL, caption string) (string, error) {
	body := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                to,
		"type":              "video",
		"video":             map[string]string{"link": videoURL, "caption": caption},
	}
	return c.doPost("/messages", body)
}

func (c *Client) SendAudio(to, audioURL string) (string, error) {
	body := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                to,
		"type":              "audio",
		"audio":             map[string]string{"link": audioURL},
	}
	return c.doPost("/messages", body)
}

// SendMedia smart dispatcher — picks the right method based on media type

func (c *Client) SendTemplate(to, templateName, language string, params []string) (string, error) {
	components := []map[string]interface{}{}
	if len(params) > 0 {
		bp := make([]map[string]string, 0)
		for _, p := range params {
			bp = append(bp, map[string]string{"type": "text", "text": p})
		}
		components = append(components, map[string]interface{}{
			"type":       "body",
			"parameters": bp,
		})
	}
	body := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                to,
		"type":              "template",
		"template": map[string]interface{}{
			"name":     templateName,
			"language": map[string]string{"code": language},
		},
	}
	if len(components) > 0 {
		body["template"].(map[string]interface{})["components"] = components
	}
	return c.doPost("/messages", body)
}

func (c *Client) SendInteractive(to string, interactive map[string]interface{}) (string, error) {
	body := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                to,
		"type":              "interactive",
		"interactive":       interactive,
	}
	return c.doPost("/messages", body)
}

func (c *Client) doPost(endpoint string, body interface{}) (string, error) {
	b, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", c.baseURL()+endpoint, bytes.NewReader(b))
	if err != nil { return "", err }
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTP.Do(req)
	if err != nil { return "", err }
	defer resp.Body.Close()
	rb, _ := io.ReadAll(resp.Body)
	var result struct {
		Messages []struct {
			ID string `json:"id"`
		} `json:"messages"`
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(rb, &result); err != nil {
		if resp.StatusCode >= 400 { return "", fmt.Errorf("meta api error: %s", string(rb)) }
		return "", nil
	}
	if result.Error.Message != "" {
		return "", fmt.Errorf("meta api error: %s", result.Error.Message)
	}
	if len(result.Messages) > 0 {
		return result.Messages[0].ID, nil
	}
	return "", nil
}

func (c *Client) VerifyWebhook(mode, challenge, verifyToken string) (string, bool) {
	if mode == "subscribe" && verifyToken == c.VerifyToken {
		return challenge, true
	}
	return "", false
}

type WebhookMessage struct {
	From      string `json:"from"`
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`
	Text      struct {
		Body string `json:"body"`
	} `json:"text"`
	Image    *MediaInfo `json:"image"`
	Video    *MediaInfo `json:"video"`
	Document *MediaInfo `json:"document"`
	Audio    *MediaInfo `json:"audio"`
	Interactive *struct {
		Type        string `json:"type"`
		ButtonReply *struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		} `json:"button_reply"`
		ListReply *struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		} `json:"list_reply"`
	} `json:"interactive"`
}

type MediaInfo struct {
	ID       string `json:"id"`
	Caption  string `json:"caption"`
	MimeType string `json:"mime_type"`
	Sha256   string `json:"sha256"`
}

func ParseWebhook(body []byte, verifyToken string) ([]WebhookMessage, string, bool) {
	var payload struct {
		Entry []struct {
			Changes []struct {
				Value struct {
					MessagingProduct string `json:"messaging_product"`
					Metadata         struct {
						DisplayPhoneNumber string `json:"display_phone_number"`
						PhoneNumberID      string `json:"phone_number_id"`
					} `json:"metadata"`
					Messages []WebhookMessage `json:"messages"`
					Statuses []struct {
						ID           string `json:"id"`
						Status       string `json:"status"`
						Timestamp    string `json:"timestamp"`
						RecipientID  string `json:"recipient_id"`
						Conversation struct {
							ID string `json:"id"`
						} `json:"conversation"`
					} `json:"statuses"`
				} `json:"value"`
			} `json:"changes"`
		} `json:"entry"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, "", false
	}
	var msgs []WebhookMessage
	phoneNumberID := ""
	for _, entry := range payload.Entry {
		for _, change := range entry.Changes {
			phoneNumberID = change.Value.Metadata.PhoneNumberID
			for _, m := range change.Value.Messages {
				phone := strings.TrimPrefix(m.From, "+")
				m.From = phone
				msgs = append(msgs, m)
			}
		}
	}
	return msgs, phoneNumberID, len(msgs) > 0
}

func (c *Client) MarkRead(messageID string) error {
	body := map[string]interface{}{
		"messaging_product": "whatsapp",
		"status":            "read",
		"message_id":        messageID,
	}
	_, err := c.doPost("/messages", body)
	return err
}

func (c *Client) FetchTemplates() ([]map[string]interface{}, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s/message_templates?limit=50", c.PhoneNumberID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	resp, err := c.HTTP.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	rb, _ := io.ReadAll(resp.Body)
	var result struct {
		Data []map[string]interface{} `json:"data"`
	}
	json.Unmarshal(rb, &result)
	return result.Data, nil
}

func (c *Client) SendMedia(to, mediaType, mediaURL, caption string) (string, error) {
	switch mediaType {
	case "image": return c.SendImage(to, mediaURL, caption)
	case "video": return c.SendVideo(to, mediaURL, caption)
	case "audio", "voice": return c.SendAudio(to, mediaURL)
	case "document": return c.SendDocument(to, mediaURL, "document", caption)
	case "sticker": return c.SendImage(to, mediaURL, caption)
	default: return c.SendImage(to, mediaURL, caption)
	}
}

// ── WhatsApp Flows (Meta native forms) ──

func (c *Client) SendFlow(to, flowID, flowToken string) (string, error) {
	body := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                to,
		"type":              "flow",
		"flow": map[string]interface{}{
			"id":           flowID,
			"token":        flowToken,
			"type":         "NAVIGATE",
			"mode":         "PUBLISHED",
		},
	}
	return c.doPost("/messages", body)
}

func (c *Client) CreateFlow(name string, categories []string, screens []map[string]interface{}) (string, error) {
	body := map[string]interface{}{
		"name":                name,
		"categories":          categories,
		"endpoint_uri":        fmt.Sprintf("https://%s/flow_callback", strings.TrimPrefix(c.baseURL(), "https://graph.facebook.com/v22.0/")),
		"screens":             screens,
	}
	resp, err := c.doPostWithURL(fmt.Sprintf("https://graph.facebook.com/v22.0/%s/whatsapp_flows", c.PhoneNumberID), body)
	return resp, err
}

func (c *Client) UpdateFlowStatus(flowID, status string) error {
	body := map[string]string{"status": status}
	_, err := c.doPostWithURL(fmt.Sprintf("https://graph.facebook.com/v22.0/%s/whatsapp_flows/%s", c.PhoneNumberID, flowID), body)
	return err
}

func (c *Client) ListFlows() ([]map[string]interface{}, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s/whatsapp_flows?fields=id,name,status,categories", c.PhoneNumberID)
	resp, err := c.HTTP.Get(url)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	rb, _ := io.ReadAll(resp.Body)
	var result struct{ Data []map[string]interface{} `json:"data"` }
	json.Unmarshal(rb, &result)
	return result.Data, nil
}

// ── Template Management ──

type Template struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Language string `json:"language"`
	Category string `json:"category"`
	Status   string `json:"status"`
}

func (c *Client) CreateTemplate(name, language, category string, components []map[string]interface{}) (string, error) {
	body := map[string]interface{}{
		"name":       name,
		"language":   language,
		"category":   category,
		"components": components,
	}
	resp, err := c.doPostWithURL(fmt.Sprintf("https://graph.facebook.com/v22.0/%s/message_templates", c.PhoneNumberID), body)
	return resp, err
}

func (c *Client) GetTemplateStatus(templateName string) (string, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s/message_templates?name=%s&fields=name,status", c.PhoneNumberID, templateName)
	resp, err := c.HTTP.Get(url)
	if err != nil { return "", err }
	defer resp.Body.Close()
	rb, _ := io.ReadAll(resp.Body)
	var result struct{ Data []Template `json:"data"` }
	json.Unmarshal(rb, &result)
	if len(result.Data) > 0 { return result.Data[0].Status, nil }
	return "not_found", nil
}

func (c *Client) DeleteTemplate(templateName string) error {
	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s/message_templates?name=%s", c.PhoneNumberID, templateName)
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	resp, err := c.HTTP.Do(req)
	if err != nil { return err }
	resp.Body.Close()
	return nil
}

// ── Helpers ──

func (c *Client) doPostWithURL(url string, body interface{}) (string, error) {
	b, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil { return "", err }
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTP.Do(req)
	if err != nil { return "", err }
	defer resp.Body.Close()
	rb, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("Meta API error %d: %s", resp.StatusCode, string(rb))
	}
	return string(rb), nil
}

func (c *Client) SendPoll(to, question string, options []string) (string, error) {
	rows := make([]map[string]interface{}, 0)
	for i, opt := range options {
		rows = append(rows, map[string]interface{}{
			"id":          fmt.Sprintf("opt_%d", i),
			"title":       opt,
			"description": "",
		})
	}
	interactive := map[string]interface{}{
		"type": "list",
		"header": map[string]string{
			"type": "text",
			"text": question,
		},
		"body": map[string]string{
			"text": "Pilih salah satu:",
		},
		"action": map[string]interface{}{
			"button":   "Pilih",
			"sections": []map[string]interface{}{{"title": "Pilihan", "rows": rows}},
		},
	}
	return c.SendInteractive(to, interactive)
}

// ── Click-to-Chat Ads ──

func (c *Client) GenerateAdLink(phone, prefillMessage string) string {
	link := fmt.Sprintf("https://wa.me/%s", phone)
	if prefillMessage != "" {
		link += "?text=" + strings.ReplaceAll(prefillMessage, " ", "%20")
	}
	return link
}

// ── WhatsApp Catalog ──

func (c *Client) SyncProduct(name, description string, price float64, imageURL, websiteURL string) (string, error) {
	body := map[string]interface{}{
		"name":         name,
		"description":  description,
		"price":        fmt.Sprintf("%.0f IDR", price),
		"images":       []string{imageURL},
		"website_url":  websiteURL,
	}
	resp, err := c.doPostWithURL(fmt.Sprintf("https://graph.facebook.com/v22.0/%s/products", c.PhoneNumberID), body)
	return resp, err
}

func (c *Client) ListProducts() ([]map[string]interface{}, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s/products?fields=id,name,price,image_url", c.PhoneNumberID)
	resp, err := c.HTTP.Get(url)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	rb, _ := io.ReadAll(resp.Body)
	var result struct{ Data []map[string]interface{} `json:"data"` }
	json.Unmarshal(rb, &result)
	return result.Data, nil
}

func (c *Client) UploadMedia(mediaURL string) (string, error) {
	// Download media from URL
	resp, err := c.HTTP.Get(mediaURL)
	if err != nil { return "", fmt.Errorf("download failed: %v", err) }
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil { return "", err }

	// Upload to Meta
	uploadURL := fmt.Sprintf("https://graph.facebook.com/v22.0/%s/media", c.PhoneNumberID)
	req, err := http.NewRequest("POST", uploadURL, bytes.NewReader(data))
	if err != nil { return "", err }
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("Content-Type", resp.Header.Get("Content-Type"))
	upResp, err := c.HTTP.Do(req)
	if err != nil { return "", err }
	defer upResp.Body.Close()
	rb, _ := io.ReadAll(upResp.Body)
	var result struct{ ID string `json:"id"` }
	json.Unmarshal(rb, &result)
	if result.ID == "" { return "", fmt.Errorf("upload failed: %s", string(rb)) }
	return result.ID, nil
}

func (c *Client) SendMediaByID(to, mediaType, mediaID, caption string) (string, error) {
	body := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                to,
		"type":              mediaType,
		mediaType: map[string]string{
			"id":      mediaID,
			"caption": caption,
		},
	}
	return c.doPost("/messages", body)
}

// ── WhatsApp Calling (WACall) ──

func (c *Client) EnableCalling() error {
	body := map[string]interface{}{
		"calling": map[string]string{"status": "ENABLED"},
	}
	_, err := c.doPostWithURL(fmt.Sprintf("https://graph.facebook.com/v22.0/%s/settings", c.PhoneNumberID), body)
	return err
}

func (c *Client) MakeCall(to, callType string, audioURL string) (string, error) {
	body := map[string]interface{}{
		"messaging_product": "whatsapp",
		"to":                to,
		"type":              "audio",
		"audio": map[string]string{
			"link": audioURL,
		},
	}
	return c.doPost("/messages", body)
}

func (c *Client) CheckCallStatus(callID string) (string, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s/calls/%s", c.PhoneNumberID, callID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	resp, err := c.HTTP.Do(req)
	if err != nil { return "", err }
	defer resp.Body.Close()
	rb, _ := io.ReadAll(resp.Body)
	var result struct{ Status string `json:"status"` }
	json.Unmarshal(rb, &result)
	return result.Status, nil
}

// ── ElevenLabs Voice ──

var ElevenLabsAPIKey string

func SetElevenLabsKey(key string) { ElevenLabsAPIKey = key }

func (c *Client) FetchElevenLabsVoices() ([]map[string]interface{}, error) {
	req, _ := http.NewRequest("GET", "https://api.elevenlabs.io/v1/voices", nil)
	req.Header.Set("xi-api-key", ElevenLabsAPIKey)
	resp, err := c.HTTP.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	rb, _ := io.ReadAll(resp.Body)
	var result struct{ Voices []map[string]interface{} `json:"voices"` }
	json.Unmarshal(rb, &result)
	return result.Voices, nil
}

func (c *Client) GenerateTTS(text, voiceID string) (string, error) {
	body := map[string]interface{}{
		"text":     text,
		"model_id": "eleven_multilingual_v2",
		"voice_settings": map[string]interface{}{
			"stability":        0.5,
			"similarity_boost": 0.75,
		},
	}
	b, _ := json.Marshal(body)
	url := fmt.Sprintf("https://api.elevenlabs.io/v1/text-to-speech/%s", voiceID)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(b))
	req.Header.Set("xi-api-key", ElevenLabsAPIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "audio/mpeg")
	resp, err := c.HTTP.Do(req)
	if err != nil { return "", err }
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	os.MkdirAll("public", 0o755)
	path := fmt.Sprintf("public/tts_%d.mp3", time.Now().UnixNano())
	os.WriteFile(path, data, 0o644)
	return "/" + path, nil
}

// ── WhatsApp Pay ──

func (c *Client) SendPaymentRequest(to, paymentToken, amount, currency string) (string, error) {
	body := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                to,
		"type":              "interactive",
		"interactive": map[string]interface{}{
			"type": "payment",
			"payment": map[string]string{
				"token":        paymentToken,
				"amount":       amount,
				"currency":     currency,
			},
		},
	}
	return c.doPost("/messages", body)
}

// ── Business Profile ──

func (c *Client) GetBusinessProfile() (map[string]interface{}, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s/settings/business/profile?fields=about,address,description,email,profile_picture_url,websites,vertical", c.PhoneNumberID)
	resp, err := c.HTTP.Get(url)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	rb, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(rb, &result)
	return result, nil
}

func (c *Client) UpdateBusinessProfile(about, email, website, description, vertical string) error {
	components := []map[string]string{}
	if about != "" { components = append(components, map[string]string{"type": "about", "text": about}) }
	if email != "" { components = append(components, map[string]string{"type": "email", "text": email}) }
	if website != "" { components = append(components, map[string]string{"type": "website", "text": website}) }
	if description != "" { components = append(components, map[string]string{"type": "description", "text": description}) }
	if vertical != "" { components = append(components, map[string]string{"type": "vertical", "text": vertical}) }
	body := map[string]interface{}{"components": components}
	_, err := c.doPostWithURL(fmt.Sprintf("https://graph.facebook.com/v22.0/%s/settings/business/profile", c.PhoneNumberID), body)
	return err
}

// ── QR Code ──

func (c *Client) GenerateQRCode(message string) (map[string]interface{}, error) {
	body := map[string]interface{}{
		"message":  message,
		"image_format": "PNG",
	}
	resp, err := c.doPostWithURL(fmt.Sprintf("https://graph.facebook.com/v22.0/%s/message_qrdls", c.PhoneNumberID), body)
	if err != nil { return nil, err }
	var result map[string]interface{}
	json.Unmarshal([]byte(resp), &result)
	return result, nil
}

// ── Health Status ──

func (c *Client) GetHealthStatus() (string, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s", c.PhoneNumberID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	resp, err := c.HTTP.Do(req)
	if err != nil { return "", err }
	defer resp.Body.Close()
	rb, _ := io.ReadAll(resp.Body)
	var result struct{ HealthStatus struct{ Entity string `json:"entity"` } `json:"health_status"` }
	json.Unmarshal(rb, &result)
	return result.HealthStatus.Entity, nil
}

// ── Message Insights ──

func (c *Client) GetMessageInsights(messageID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s/messages?fields=id,message_status,timestamp,type", messageID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	resp, err := c.HTTP.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	rb, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(rb, &result)
	return result, nil
}

// ── Rate Limit ──

func (c *Client) GetRateLimit() (map[string]interface{}, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s/settings/business/rate_limits", c.PhoneNumberID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	resp, err := c.HTTP.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	rb, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(rb, &result)
	return result, nil
}

// ── Reply with Context ──

func (c *Client) SendReply(to, replyToMessageID, message string) (string, error) {
	body := map[string]interface{}{
		"messaging_product": "whatsapp",
		"to":                to,
		"type":              "text",
		"text":              map[string]string{"body": message, "preview_url": "false"},
		"context":           map[string]string{"message_id": replyToMessageID},
	}
	return c.doPost("/messages", body)
}

// ── Reaction ──

func (c *Client) SendReaction(to, messageID, emoji string) (string, error) {
	body := map[string]interface{}{
		"messaging_product": "whatsapp",
		"to":                to,
		"type":              "reaction",
		"reaction":          map[string]string{"message_id": messageID, "emoji": emoji},
	}
	return c.doPost("/messages", body)
}

// ── Location Message ──

func (c *Client) SendLocation(to, name, address string, lat, lng float64) (string, error) {
	body := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                to,
		"type":              "location",
		"location":          map[string]interface{}{"latitude": fmt.Sprintf("%.6f", lat), "longitude": fmt.Sprintf("%.6f", lng), "name": name, "address": address},
	}
	return c.doPost("/messages", body)
}

// ── Contact Card ──

func (c *Client) SendContacts(to string, contacts []map[string]interface{}) (string, error) {
	body := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                to,
		"type":              "contacts",
		"contacts":          contacts,
	}
	return c.doPost("/messages", body)
}

// ── Webhook Subscribe ──

func (c *Client) SubscribeWebhook(fields []string) error {
	body := map[string]interface{}{"fields": fields}
	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s/subscribed_apps", c.PhoneNumberID)
	_, err := c.doPostWithURL(url, body)
	return err
}

// ── Business Verification ──

func (c *Client) GetVerificationStatus(businessID string) (string, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s?fields=verification_status", businessID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	resp, err := c.HTTP.Do(req)
	if err != nil { return "", err }
	defer resp.Body.Close()
	rb, _ := io.ReadAll(resp.Body)
	var result struct{ VerificationStatus string `json:"verification_status"` }
	json.Unmarshal(rb, &result)
	return result.VerificationStatus, nil
}

// ── Flows Data Exchange (callback endpoint) ──

func (c *Client) SetFlowEndpoint(flowID, endpointURL string) error {
	body := map[string]string{"endpoint_uri": endpointURL}
	_, err := c.doPostWithURL(fmt.Sprintf("https://graph.facebook.com/v22.0/%s/whatsapp_flows/%s", c.PhoneNumberID, flowID), body)
	return err
}

func (c *Client) GetFlowAnalytics(flowID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s/whatsapp_flows/%s?fields=id,name,status,categories,validation_errors,message_template,analytics", c.PhoneNumberID, flowID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	resp, err := c.HTTP.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	rb, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(rb, &result)
	return result, nil
}

// ── Embedded Signup (OAuth) ──

type EmbeddedSignupConfig struct {
	AppID     string
	AppSecret string
}

func (c *Client) GenerateOAuthURL(configID, redirectURI, state string) string {
	return fmt.Sprintf("https://www.facebook.com/v22.0/dialog/oauth?client_id=%s&redirect_uri=%s&state=%s&config_id=%s&response_type=code&override_default_response_type=true", c.PhoneNumberID, redirectURI, state, configID)
}

func (c *Client) ExchangeOAuthCode(code string, config EmbeddedSignupConfig) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v22.0/oauth/access_token?client_id=%s&client_secret=%s&code=%s&redirect_uri=https://%s/auth/callback", config.AppID, config.AppSecret, code, config.AppID)
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := c.HTTP.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	rb, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(rb, &result)
	return result, nil
}

func (c *Client) DebugToken(accessToken string) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v22.0/debug_token?input_token=%s&access_token=%s", accessToken, c.AccessToken)
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := c.HTTP.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	rb, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(rb, &result)
	return result, nil
}

// ── Carousel / Catalog Messages ──

func (c *Client) SendCarousel(to string, products []map[string]string) (string, error) {
	sections := make([]map[string]interface{}, 0)
	rows := make([]map[string]interface{}, 0)
	for i, p := range products {
		rows = append(rows, map[string]interface{}{
			"id":          fmt.Sprintf("p%d", i),
			"title":       p["title"],
			"description": p["description"],
		})
	}
	sections = append(sections, map[string]interface{}{
		"title": "Produk Kami",
		"rows":  rows,
	})
	interactive := map[string]interface{}{
		"type": "catalog_message",
		"body": map[string]string{"text": "Pilih produk:"},
		"action": map[string]interface{}{
			"name":     "catalog_message",
			"sections": sections,
		},
	}
	return c.SendInteractive(to, interactive)
}

// ── Quality Score ──

func (c *Client) GetQualityScore() (map[string]interface{}, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s?fields=name,quality_rating,message_limit,current_limit,health_status", c.PhoneNumberID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	resp, err := c.HTTP.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	rb, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(rb, &result)
	return result, nil
}

// ── WABA Migration ──

func (c *Client) RegisterPhoneNumber(pin string) error {
	body := map[string]string{
		"messaging_product": "whatsapp",
		"pin":               pin,
	}
	_, err := c.doPostWithURL(fmt.Sprintf("https://graph.facebook.com/v22.0/%s/register", c.PhoneNumberID), body)
	return err
}

func (c *Client) DeregisterPhoneNumber() error {
	body := map[string]string{"messaging_product": "whatsapp"}
	_, err := c.doPostWithURL(fmt.Sprintf("https://graph.facebook.com/v22.0/%s/deregister", c.PhoneNumberID), body)
	return err
}

func (c *Client) GetPhoneNumberStatus() (map[string]interface{}, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s?fields=id,display_phone_number,verified_name,code_verification_status,quality_rating,platform_type,throughput,health_status", c.PhoneNumberID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	resp, err := c.HTTP.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	rb, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(rb, &result)
	return result, nil
}

// ── Instagram Messaging ──

func (c *Client) SendIGMessage(to, message string) (string, error) {
	body := map[string]interface{}{
		"recipient": map[string]string{"id": to},
		"message":   map[string]string{"text": message},
		"messaging_type": "RESPONSE",
	}
	return c.doPostWithURL(fmt.Sprintf("https://graph.facebook.com/v22.0/%s/messages", c.PhoneNumberID), body)
}

func (c *Client) SendIGMedia(to, mediaURL, caption string) (string, error) {
	body := map[string]interface{}{
		"recipient": map[string]string{"id": to},
		"message": map[string]interface{}{
			"attachment": map[string]interface{}{
				"type": "image",
				"payload": map[string]string{"url": mediaURL},
			},
		},
		"messaging_type": "RESPONSE",
	}
	if caption != "" { body["message"].(map[string]interface{})["text"] = caption }
	return c.doPostWithURL(fmt.Sprintf("https://graph.facebook.com/v22.0/%s/messages", c.PhoneNumberID), body)
}

func (c *Client) GetIGConversations() ([]map[string]interface{}, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s/conversations?platform=instagram&fields=participants,messages.limit(10){message,from,created_time}", c.PhoneNumberID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	resp, err := c.HTTP.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	rb, _ := io.ReadAll(resp.Body)
	var result struct{ Data []map[string]interface{} `json:"data"` }
	json.Unmarshal(rb, &result)
	return result.Data, nil
}

func (c *Client) DeleteProduct(productID string) error {
	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s/products/%s", c.PhoneNumberID, productID)
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	resp, err := c.HTTP.Do(req)
	if err != nil { return err }
	resp.Body.Close()
	return nil
}

func ParseIGWebhook(body []byte) ([]WebhookMessage, bool) {
	var raw struct {
		Entry []struct {
			Messaging []struct {
				Sender  struct { ID string }
				Message struct {
					MID  string
					Text string
				}
			}
		}
	}
	json.Unmarshal(body, &raw)
	var msgs []WebhookMessage
	for _, entry := range raw.Entry {
		for _, m := range entry.Messaging {
			msgs = append(msgs, WebhookMessage{
				From: m.Sender.ID, ID: m.Message.MID,
			})
		}
	}
	return msgs, len(msgs) > 0
}

// -- Instagram Comment Automation --

func (c *Client) ReplyToComment(commentID, message string) (string, error) {
	body := map[string]string{"message": message}
	return c.doPostWithURL(fmt.Sprintf("https://graph.facebook.com/v22.0/%s/replies?access_token=%s", commentID, c.AccessToken), body)
}

func (c *Client) SendDMFromComment(igUserID, message string) (string, error) {
	return c.SendIGMessage(igUserID, message)
}

func (c *Client) HideComment(commentID string) error {
	body := map[string]bool{"hide": true}
	_, err := c.doPostWithURL(fmt.Sprintf("https://graph.facebook.com/v22.0/%s?access_token=%s", commentID, c.AccessToken), body)
	return err
}

// -- Instagram Story Mentions --

func (c *Client) GetStoryMentions(igUserID string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s/mentioned_media?fields=id,media_type,media_url,owner,timestamp", igUserID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	resp, err := c.HTTP.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	rb, _ := io.ReadAll(resp.Body)
	var result struct{ Data []map[string]interface{} `json:"data"` }
	json.Unmarshal(rb, &result)
	return result.Data, nil
}

// -- Instagram Post Publishing --

func (c *Client) CreateMediaContainer(igUserID, mediaURL, caption string, isVideo bool) (string, error) {
	body := map[string]interface{}{
		"image_url": mediaURL,
		"caption":   caption,
	}
	if isVideo {
		body = map[string]interface{}{
			"video_url": mediaURL,
			"caption":   caption,
			"media_type": "REELS",
		}
	}
	resp, err := c.doPostWithURL(fmt.Sprintf("https://graph.facebook.com/v22.0/%s/media", igUserID), body)
	if err != nil { return "", err }
	var result struct{ ID string `json:"id"` }
	json.Unmarshal([]byte(resp), &result)
	return result.ID, nil
}

func (c *Client) PublishMedia(igUserID, creationID string) (string, error) {
	body := map[string]string{"creation_id": creationID}
	resp, err := c.doPostWithURL(fmt.Sprintf("https://graph.facebook.com/v22.0/%s/media_publish", igUserID), body)
	if err != nil { return "", err }
	var result struct{ ID string `json:"id"` }
	json.Unmarshal([]byte(resp), &result)
	return result.ID, nil
}

// -- Instagram Quick Replies --

func (c *Client) SendIGQuickReplies(to, message string, replies []string) (string, error) {
	quickReplies := make([]map[string]string, 0)
	for _, r := range replies {
		quickReplies = append(quickReplies, map[string]string{
			"content_type": "text",
			"title":        r,
			"payload":      r,
		})
	}
	body := map[string]interface{}{
		"recipient": map[string]string{"id": to},
		"message": map[string]interface{}{
			"text":          message,
			"quick_replies": quickReplies,
		},
		"messaging_type": "RESPONSE",
	}
	return c.doPostWithURL(fmt.Sprintf("https://graph.facebook.com/v22.0/%s/messages", c.PhoneNumberID), body)
}

// -- Instagram Broadcast --

func (c *Client) SendIGBroadcast(igUserIDs []string, message string) ([]string, error) {
	var results []string
	for _, id := range igUserIDs {
		resp, err := c.SendIGMessage(id, message)
		if err == nil { results = append(results, resp) }
		time.Sleep(500 * time.Millisecond)
	}
	return results, nil
}

// -- Instagram Insights --

func (c *Client) GetIGAccountInsights(igUserID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s/insights?metric=impressions,reach,profile_views,website_clicks,follower_count,email_contacts,phone_call_clicks,text_message_clicks,direction_clicks&period=day&metric_type=total_value", igUserID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	resp, err := c.HTTP.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	rb, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(rb, &result)
	return result, nil
}

func (c *Client) GetIGMediaInsights(mediaID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s/insights?metric=engagement,impressions,reach,saved,video_views,likes,comments,shares", mediaID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	resp, err := c.HTTP.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	rb, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(rb, &result)
	return result, nil
}

func (c *Client) GetIGAudienceInsights(igUserID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s/insights?metric=audience_gender_age,audience_locale,audience_country,audience_city,online_followers&period=lifetime", igUserID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	resp, err := c.HTTP.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	rb, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(rb, &result)
	return result, nil
}

// -- Instagram Story + Reel Publisher --

func (c *Client) CreateStory(igUserID, mediaURL, caption string) (string, error) {
	body := map[string]interface{}{
		"media_type": "STORIES",
		"image_url":  mediaURL,
		"caption":    caption,
	}
	resp, err := c.doPostWithURL(fmt.Sprintf("https://graph.facebook.com/v22.0/%s/media", igUserID), body)
	if err != nil { return "", err }
	var result struct{ ID string `json:"id"` }
	json.Unmarshal([]byte(resp), &result)
	_, err = c.doPostWithURL(fmt.Sprintf("https://graph.facebook.com/v22.0/%s/media_publish?creation_id=%s", igUserID, result.ID),
		map[string]string{"creation_id": result.ID})
	return result.ID, err
}

func (c *Client) CreateReel(igUserID, videoURL, caption, shareToFeed string) (string, error) {
	body := map[string]interface{}{
		"media_type":  "REELS",
		"video_url":   videoURL,
		"caption":     caption,
		"share_to_feed": shareToFeed == "true",
	}
	resp, err := c.doPostWithURL(fmt.Sprintf("https://graph.facebook.com/v22.0/%s/media", igUserID), body)
	if err != nil { return "", err }
	var result struct{ ID string `json:"id"` }
	json.Unmarshal([]byte(resp), &result)
	time.Sleep(5 * time.Second)
	_, err = c.doPostWithURL(fmt.Sprintf("https://graph.facebook.com/v22.0/%s/media_publish?creation_id=%s", igUserID, result.ID),
		map[string]string{"creation_id": result.ID})
	return result.ID, err
}

// -- Instagram Auto-Responder + Handover --

func (c *Client) SetIGIceBreakers(questions []string) error {
	body := map[string]interface{}{
		"ice_breakers": questions,
	}
	_, err := c.doPostWithURL(fmt.Sprintf("https://graph.facebook.com/v22.0/%s/messenger_profile", c.PhoneNumberID), body)
	return err
}

func (c *Client) SetIGGreeting(text string) error {
	body := map[string]interface{}{
		"greeting": []map[string]string{{"locale": "default", "text": text}},
	}
	_, err := c.doPostWithURL(fmt.Sprintf("https://graph.facebook.com/v22.0/%s/messenger_profile", c.PhoneNumberID), body)
	return err
}

func (c *Client) HandoverToAgent(igUserID, agentID string) error {
	body := map[string]interface{}{
		"recipient": map[string]string{"id": igUserID},
		"target_app_id": agentID,
	}
	_, err := c.doPostWithURL(fmt.Sprintf("https://graph.facebook.com/v22.0/%s/pass_thread_control", c.PhoneNumberID), body)
	return err
}

func (c *Client) TakeThreadControl(igUserID string) error {
	body := map[string]interface{}{
		"recipient": map[string]string{"id": igUserID},
	}
	_, err := c.doPostWithURL(fmt.Sprintf("https://graph.facebook.com/v22.0/%s/take_thread_control", c.PhoneNumberID), body)
	return err
}

// -- Story Reply Detection --

type IGWebhookEntry struct {
	ID        string              `json:"id"`
	Messaging []IGMessagingEvent  `json:"messaging,omitempty"`
	Changes   []IGChangeEvent     `json:"changes,omitempty"`
}

type IGMessagingEvent struct {
	Sender    struct{ ID string } `json:"sender"`
	Recipient struct{ ID string } `json:"recipient"`
	Message   struct {
		MID      string `json:"mid"`
		Text     string `json:"text"`
		ReplyTo  struct {
			MID   string `json:"mid"`
			Story struct {
				ID    string `json:"id"`
				URL   string `json:"url"`
			} `json:"story"`
		} `json:"reply_to"`
	} `json:"message"`
}

func ParseIGStoryReply(body []byte) (storyID, fromUser, messageText string, isReply bool) {
	var raw struct {
		Entry []IGWebhookEntry `json:"entry"`
	}
	if err := json.Unmarshal(body, &raw); err != nil { return }
	for _, entry := range raw.Entry {
		for _, m := range entry.Messaging {
			if m.Message.ReplyTo.Story.ID != "" {
				return m.Message.ReplyTo.Story.ID, m.Sender.ID, m.Message.Text, true
			}
		}
	}
	return "", "", "", false
}

// -- Instagram Product Tagging --

func (c *Client) TagProductInMedia(igUserID, mediaID, productID string, x, y float64) error {
	body := map[string]interface{}{
		"product_tags": []map[string]interface{}{
			{"product_id": productID, "x": x, "y": y},
		},
	}
	_, err := c.doPostWithURL(fmt.Sprintf("https://graph.facebook.com/v22.0/%s?access_token=%s", mediaID, c.AccessToken), body)
	return err
}

// -- Instagram Broadcast Channel --

func (c *Client) CreateBroadcastChannel(igUserID, name, description string) (string, error) {
	body := map[string]interface{}{
		"name":        name,
		"description": description,
	}
	resp, err := c.doPostWithURL(fmt.Sprintf("https://graph.facebook.com/v22.0/%s/broadcast_channels", igUserID), body)
	if err != nil { return "", err }
	var result struct{ ID string `json:"id"` }
	json.Unmarshal([]byte(resp), &result)
	return result.ID, nil
}

func (c *Client) SendBroadcastMessage(channelID, message string) (string, error) {
	body := map[string]interface{}{
		"message": message,
	}
	return c.doPostWithURL(fmt.Sprintf("https://graph.facebook.com/v22.0/%s/messages", channelID), body)
}
