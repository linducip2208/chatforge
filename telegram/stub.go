//go:build !pro

package telegram

type Bot struct{ Token string }
type Update struct{ Message Message }
type Message struct{ Chat Chat; Text string }
type Chat struct{ ID int64; FirstName string }
type InlineKeyboardButton struct{ Text string; CallbackData string }

func New(token string) *Bot { return &Bot{} }
func ParseUpdate(body []byte) (*Update, error) { return nil, nil }
func (b *Bot) SendMessage(chatID int64, text string) error { return nil }
func (b *Bot) SendInlineKeyboard(chatID int64, text string, rows [][]InlineKeyboardButton) error { return nil }
