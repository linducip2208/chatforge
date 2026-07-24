//go:build !pro

package shopee

import "net/http"

type Client struct{}

func New(partnerID int64, partnerKey string, shopID int64, accessToken string) *Client { return &Client{} }
func (c *Client) GetMessages(conversationID string, offset int) ([]map[string]interface{}, error) { return nil, nil }
func (c *Client) SendMessage(toShopID int64, text, msgType string) error { return nil }
func (c *Client) GetConversations(offset int) ([]map[string]interface{}, error) { return nil, nil }
func HandleWebhook(partnerKey string, r *http.Request) ([]map[string]interface{}, error) { return nil, nil }
