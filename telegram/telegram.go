package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Bot struct {
	Token  string
	client *http.Client
}

type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	MessageID int    `json:"message_id"`
	Chat      Chat   `json:"chat"`
	Text      string `json:"text"`
	Photo     []PhotoSize `json:"photo"`
	Document  Document    `json:"document"`
	Video     Video       `json:"video"`
	Audio     Audio       `json:"audio"`
	Voice     Voice       `json:"voice"`
}

type Chat struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

type PhotoSize struct {
	FileID   string `json:"file_id"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

type Document struct {
	FileID   string `json:"file_id"`
	FileName string `json:"file_name"`
}

type Video struct {
	FileID string `json:"file_id"`
}

type Audio struct {
	FileID string `json:"file_id"`
}

type Voice struct {
	FileID string `json:"file_id"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data,omitempty"`
	URL          string `json:"url,omitempty"`
}

type ReplyKeyboardMarkup struct {
	Keyboard        [][]KeyboardButton `json:"keyboard"`
	ResizeKeyboard  bool               `json:"resize_keyboard"`
	OneTimeKeyboard bool               `json:"one_time_keyboard"`
}

type KeyboardButton struct {
	Text string `json:"text"`
}

func New(token string) *Bot {
	return &Bot{Token: token, client: &http.Client{}}
}

func (b *Bot) apiURL(method string) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/%s", b.Token, method)
}

func (b *Bot) doPost(method string, body interface{}) ([]byte, error) {
	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", b.apiURL(method), bytes.NewReader(jsonBody))
	if err != nil { return nil, err }
	req.Header.Set("Content-Type", "application/json")
	resp, err := b.client.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (b *Bot) SendMessage(chatID int64, text string) error {
	_, err := b.doPost("sendMessage", map[string]interface{}{
		"chat_id": chatID, "text": text, "parse_mode": "HTML",
	})
	return err
}

func (b *Bot) SendPhoto(chatID int64, photoURL, caption string) error {
	_, err := b.doPost("sendPhoto", map[string]interface{}{
		"chat_id": chatID, "photo": photoURL, "caption": caption,
	})
	return err
}

func (b *Bot) SendDocument(chatID int64, docURL, caption string) error {
	_, err := b.doPost("sendDocument", map[string]interface{}{
		"chat_id": chatID, "document": docURL, "caption": caption,
	})
	return err
}

func (b *Bot) SendVideo(chatID int64, videoURL, caption string) error {
	_, err := b.doPost("sendVideo", map[string]interface{}{
		"chat_id": chatID, "video": videoURL, "caption": caption,
	})
	return err
}

func (b *Bot) SendAudio(chatID int64, audioURL, caption string) error {
	_, err := b.doPost("sendAudio", map[string]interface{}{
		"chat_id": chatID, "audio": audioURL, "caption": caption,
	})
	return err
}

func (b *Bot) SendVoice(chatID int64, voiceURL string) error {
	_, err := b.doPost("sendVoice", map[string]interface{}{
		"chat_id": chatID, "voice": voiceURL,
	})
	return err
}

func (b *Bot) SendLocation(chatID int64, lat, lng float64, title string) error {
	_, err := b.doPost("sendLocation", map[string]interface{}{
		"chat_id":   chatID,
		"latitude":  fmt.Sprintf("%.6f", lat),
		"longitude": fmt.Sprintf("%.6f", lng),
	})
	return err
}

func (b *Bot) SendContact(chatID int64, phone, firstName string) error {
	_, err := b.doPost("sendContact", map[string]interface{}{
		"chat_id":    chatID,
		"phone_number": phone,
		"first_name": firstName,
	})
	return err
}

func (b *Bot) SendPoll(chatID int64, question string, options []string) error {
	_, err := b.doPost("sendPoll", map[string]interface{}{
		"chat_id": chatID, "question": question, "options": options,
		"is_anonymous": false,
	})
	return err
}

func (b *Bot) SendInlineKeyboard(chatID int64, text string, buttons [][]InlineKeyboardButton) error {
	keyboard := InlineKeyboardMarkup{InlineKeyboard: buttons}
	kbJSON, _ := json.Marshal(keyboard)
	_, err := b.doPost("sendMessage", map[string]interface{}{
		"chat_id":      chatID,
		"text":         text,
		"reply_markup": json.RawMessage(kbJSON),
	})
	return err
}

func (b *Bot) SendReplyKeyboard(chatID int64, text string, buttons [][]string) error {
	var kb [][]KeyboardButton
	for _, row := range buttons {
		var btnRow []KeyboardButton
		for _, b := range row {
			btnRow = append(btnRow, KeyboardButton{Text: b})
		}
		kb = append(kb, btnRow)
	}
	keyboard := ReplyKeyboardMarkup{Keyboard: kb, ResizeKeyboard: true}
	kbJSON, _ := json.Marshal(keyboard)
	_, err := b.doPost("sendMessage", map[string]interface{}{
		"chat_id":      chatID,
		"text":         text,
		"reply_markup": json.RawMessage(kbJSON),
	})
	return err
}

func (b *Bot) SetWebhook(url string) error {
	_, err := b.doPost("setWebhook", map[string]string{"url": url})
	return err
}

func (b *Bot) DeleteWebhook() error {
	_, err := b.doPost("deleteWebhook", nil)
	return err
}

func (b *Bot) GetWebhookInfo() (map[string]interface{}, error) {
	rb, err := b.doPost("getWebhookInfo", nil)
	if err != nil { return nil, err }
	var result map[string]interface{}
	json.Unmarshal(rb, &result)
	return result, nil
}

func (b *Bot) GetMe() (map[string]interface{}, error) {
	rb, err := b.doPost("getMe", nil)
	if err != nil { return nil, err }
	var result map[string]interface{}
	json.Unmarshal(rb, &result)
	return result, nil
}

func ParseUpdate(body []byte) (*Update, error) {
	var update Update
	if err := json.Unmarshal(body, &update); err != nil { return nil, err }
	return &update, nil
}

func ParseChatID(body []byte) int64 {
	var update Update
	if err := json.Unmarshal(body, &update); err != nil { return 0 }
	return update.Message.Chat.ID
}

type InlineQuery struct {
	ID     string `json:"id"`
	From   struct{ ID int64 `json:"id"` } `json:"from"`
	Query  string `json:"query"`
	Offset string `json:"offset"`
}

func ParseInlineQuery(body []byte) (*InlineQuery, error) {
	var raw struct{ InlineQuery InlineQuery `json:"inline_query"` }
	if err := json.Unmarshal(body, &raw); err != nil { return nil, err }
	return &raw.InlineQuery, nil
}

func (b *Bot) AnswerInlineQuery(queryID string, results []map[string]interface{}) error {
	_, err := b.doPost("answerInlineQuery", map[string]interface{}{
		"inline_query_id": queryID, "results": results,
	})
	return err
}
