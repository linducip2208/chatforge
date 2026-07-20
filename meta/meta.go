package meta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
