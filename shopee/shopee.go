//go:build pro

package shopee

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	PartnerID  int64
	PartnerKey string
	ShopID     int64
	AccessToken string
	BaseURL    string
}

func New(partnerID int64, partnerKey string, shopID int64, accessToken string) *Client {
	return &Client{
		PartnerID: partnerID, PartnerKey: partnerKey,
		ShopID: shopID, AccessToken: accessToken,
		BaseURL: "https://partner.shopeemobile.com",
	}
}

type chatMessage struct {
	ID         string `json:"message_id"`
	FromShopID int64  `json:"from_shop_id"`
	ToShopID   int64  `json:"to_shop_id"`
	Content    struct {
		Text string `json:"text"`
	} `json:"content"`
	CreatedAt int64 `json:"created_timestamp"`
}

func (c *Client) sign(path string, body string) (string, string) {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	baseStr := fmt.Sprintf("%s%s%s%s%s", c.PartnerKey, path, timestamp, c.AccessToken, fmt.Sprintf("%d", c.ShopID))
	h := hmac.New(sha256.New, []byte(c.PartnerKey))
	h.Write([]byte(baseStr))
	return timestamp, hex.EncodeToString(h.Sum(nil))
}

func (c *Client) GetMessages(conversationID string, offset int) ([]chatMessage, error) {
	path := "/api/v2/sellerchat/get_message"
	ts, sign := c.sign(path, "")
	url := fmt.Sprintf("%s%s?partner_id=%d&timestamp=%s&sign=%s&access_token=%s&shop_id=%d&conversation_id=%s&page_size=50&offset=%d",
		c.BaseURL, path, c.PartnerID, ts, sign, c.AccessToken, c.ShopID, conversationID, offset)
	resp, err := http.Get(url)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Messages []chatMessage `json:"response"`
	}
	json.Unmarshal(body, &result)
	return result.Messages, nil
}

func (c *Client) SendMessage(toShopID int64, text, msgType string) error {
	path := "/api/v2/sellerchat/send_message"
	reqBody, _ := json.Marshal(map[string]interface{}{
		"to_shop_id": toShopID,
		"message_type": msgType,
		"content": map[string]string{"text": text},
	})
	ts, sign := c.sign(path, string(reqBody))
	url := fmt.Sprintf("%s%s?partner_id=%d&timestamp=%s&sign=%s&access_token=%s&shop_id=%d",
		c.BaseURL, path, c.PartnerID, ts, sign, c.AccessToken, c.ShopID)
	resp, err := http.Post(url, "application/json", bytes.NewReader(reqBody))
	if err != nil { return err }
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("shopee api error %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

func (c *Client) GetConversations(offset int) ([]map[string]interface{}, error) {
	path := "/api/v2/sellerchat/get_conversation_list"
	ts, sign := c.sign(path, "")
	url := fmt.Sprintf("%s%s?partner_id=%d&timestamp=%s&sign=%s&access_token=%s&shop_id=%d&page_size=50&offset=%d",
		c.BaseURL, path, c.PartnerID, ts, sign, c.AccessToken, c.ShopID, offset)
	resp, err := http.Get(url)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Conversations []map[string]interface{} `json:"response"`
	}
	json.Unmarshal(body, &result)
	return result.Conversations, nil
}

func HandleWebhook(partnerKey string, r *http.Request) ([]chatMessage, error) {
	body, _ := io.ReadAll(r.Body)
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		var authPrefix = "sha256="
		if len(authHeader) > len(authPrefix) {
			h := hmac.New(sha256.New, []byte(partnerKey))
			h.Write(body)
			expected := hex.EncodeToString(h.Sum(nil))
			if authHeader[len(authPrefix):] != expected {
				return nil, fmt.Errorf("invalid signature")
			}
		}
	}
	var msgs []chatMessage
	json.Unmarshal(body, &msgs)
	return msgs, nil
}
