//go:build !pro

package discord

import (
	"fmt"
	"net/http"
)

type Client struct{}
func New(token, appID, publicKey string) *Client { return &Client{} }
func (c *Client) SendMessage(channelID, content string) error { return nil }
func (c *Client) GetDMChannelsWithPreview() ([]map[string]string, error) { return nil, nil }
func ParseWebhook(r *http.Request, publicKey string) (string, string, string, string, string, error) { return "", "", "", "", "", nil }
func RespondPing(w http.ResponseWriter, publicKey string) { fmt.Fprint(w, `{"type":1}`) }
