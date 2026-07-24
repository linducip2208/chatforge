//go:build pro

package discord

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	Token     string
	AppID     string
	PublicKey string
}

func New(token, appID, publicKey string) *Client {
	return &Client{Token: token, AppID: appID, PublicKey: publicKey}
}

type Interaction struct {
	ID        string `json:"id"`
	Type      int    `json:"type"`
	ChannelID string `json:"channel_id"`
	Member    struct {
		User struct {
			ID       string `json:"id"`
			Username string `json:"username"`
		} `json:"user"`
	} `json:"member"`
	Data struct {
		Name    string `json:"name"`
		Options []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"options"`
	} `json:"data"`
	Token string `json:"token"`
}

func (c *Client) VerifySignature(r *http.Request, body []byte) bool {
	signature := r.Header.Get("X-Signature-Ed25519")
	timestamp := r.Header.Get("X-Signature-Timestamp")
	if signature == "" || timestamp == "" { return false }
	pubKey, _ := hex.DecodeString(c.PublicKey)
	msg := []byte(timestamp + string(body))
	sig, _ := hex.DecodeString(signature)
	return ed25519.Verify(pubKey, msg, sig)
}

func ParseInteraction(r *http.Request) (*Interaction, error) {
	body, _ := io.ReadAll(r.Body)
	var i Interaction
	if err := json.Unmarshal(body, &i); err != nil { return nil, err }
	return &i, nil
}

func (c *Client) Reply(interaction *Interaction, content string) error {
	url := fmt.Sprintf("https://discord.com/api/v10/interactions/%s/%s/callback", interaction.ID, interaction.Token)
	payload := map[string]interface{}{"type": 4, "data": map[string]string{"content": content}}
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil { return err }
	defer resp.Body.Close()
	return nil
}

func (c *Client) SendMessage(channelID, content string) error {
	url := fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages", channelID)
	payload := map[string]string{"content": content}
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(b))
	req.Header.Set("Authorization", "Bot "+c.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil { return err }
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("discord api error %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

func (c *Client) GetUser(userID int64) (string, error) {
	url := fmt.Sprintf("https://discord.com/api/v10/users/%d", userID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bot "+c.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil { return "", err }
	defer resp.Body.Close()
	var u struct {
		Username string `json:"username"`
	}
	json.NewDecoder(resp.Body).Decode(&u)
	return u.Username, nil
}

type GuildChannel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type int    `json:"type"`
}

func (c *Client) GetGuildChannels(guildID string) ([]GuildChannel, error) {
	url := fmt.Sprintf("https://discord.com/api/v10/guilds/%s/channels", guildID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bot "+c.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	var chs []GuildChannel
	json.NewDecoder(resp.Body).Decode(&chs)
	return chs, nil
}

func (c *Client) GetGuildMember(guildID, userID string) (string, error) {
	url := fmt.Sprintf("https://discord.com/api/v10/guilds/%s/members/%s", guildID, userID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bot "+c.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil { return "", err }
	defer resp.Body.Close()
	var m struct { User struct { Username string `json:"username"` } `json:"user"` }
	json.NewDecoder(resp.Body).Decode(&m)
	return m.User.Username, nil
}

func (c *Client) GetDMChannels() ([]map[string]string, error) {
	url := "https://discord.com/api/v10/users/@me/channels"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bot "+c.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	var chs []struct {
		ID        string   `json:"id"`
		Recipients []struct {
			ID       string `json:"id"`
			Username string `json:"username"`
		} `json:"recipients"`
	}
	json.NewDecoder(resp.Body).Decode(&chs)
	var result []map[string]string
	for _, ch := range chs {
		name := ch.ID
		if len(ch.Recipients) > 0 { name = ch.Recipients[0].Username }
		result = append(result, map[string]string{"id": ch.ID, "name": name})
	}
	return result, nil
}

func (c *Client) GetMessages(channelID string, limit int) ([]map[string]string, error) {
	url := fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages?limit=%d", channelID, limit)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bot "+c.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	var msgs []struct {
		Author struct { Username string `json:"username"` } `json:"author"`
		Content string `json:"content"`
		ID      string `json:"id"`
	}
	json.NewDecoder(resp.Body).Decode(&msgs)
	var result []map[string]string
	for _, m := range msgs {
		result = append(result, map[string]string{"author": m.Author.Username, "content": m.Content, "id": m.ID})
	}
	return result, nil
}

func SendInteractionResponse(w http.ResponseWriter, content string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"type": 4, "data": map[string]string{"content": content}})
}

func ParseSlashCommand(body []byte) (string, string, string, error) {
	var i Interaction
	if err := json.Unmarshal(body, &i); err != nil { return "", "", "", err }
	if i.Type == 1 { return "", "", "", fmt.Errorf("ping") } // PING
	if i.Type == 2 {
		cmdName := ""
		cmdArgs := ""
		if i.Data.Name != "" { cmdName = i.Data.Name }
		for _, opt := range i.Data.Options { cmdArgs += opt.Name + ":" + opt.Value + " " }
		return i.Member.User.ID, i.Member.User.Username, cmdName + " " + cmdArgs, nil
	}
	return "", "", "", fmt.Errorf("unknown interaction type: %d", i.Type)
}

func RespondPing(w http.ResponseWriter, publicKey string) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"type":1}`)
}

type ChannelInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	LastMsg  string `json:"last_msg"`
	MsgCount int    `json:"msg_count"`
}

func (c *Client) GetDMChannelsWithPreview() ([]ChannelInfo, error) {
	chs, err := c.GetDMChannels()
	if err != nil { return nil, err }
	var result []ChannelInfo
	for _, ch := range chs {
		msgs, _ := c.GetMessages(ch["id"], 1)
		last := ""
		if len(msgs) > 0 { last = msgs[0]["content"] }
		result = append(result, ChannelInfo{ID: ch["id"], Name: ch["name"], LastMsg: last, MsgCount: 1})
	}
	return result, nil
}

func ParseWebhook(r *http.Request, publicKey string) (string, string, string, string, string, error) {
	body, _ := io.ReadAll(r.Body)
	sig := r.Header.Get("X-Signature-Ed25519")
	ts := r.Header.Get("X-Signature-Timestamp")
	if sig == "" || ts == "" { return "", "", "", "", "", fmt.Errorf("no signature") }
	pubKey, _ := hex.DecodeString(publicKey)
	msg := []byte(ts + string(body))
	sigBytes, _ := hex.DecodeString(sig)
	if !ed25519.Verify(pubKey, msg, sigBytes) { return "", "", "", "", "", fmt.Errorf("invalid signature") }

	var i Interaction
	json.Unmarshal(body, &i)
	if i.Type == 1 { return "", "", "", "", "PING", nil }
	if i.Type == 2 {
		cmd := i.Data.Name
		args := ""
		for _, o := range i.Data.Options { args += o.Value + " " }
		text := cmd
		if args != "" { text += " " + args }
		return i.Member.User.ID, i.Member.User.Username, i.ChannelID, text, "command", nil
	}
	return "", "", "", "", "", nil
}
