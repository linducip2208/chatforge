//go:build !pro

package xcom

import "net/http"

type Client struct{}
type DMEvent struct{ SenderID string; Text string; CreatedAt string }

func New(apiKey, apiSecret, accessToken, accessSecret, userID string) *Client { return &Client{} }
func (c *Client) SendDM(recipientID, message string) error { return nil }
func ParseWebhook(r *http.Request, consumerSecret string) ([]DMEvent, error) { return nil, nil }
func CRCResponse(consumerSecret, crcToken string) string { return "" }
