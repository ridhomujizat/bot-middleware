package webHookTelegram

import "encoding/json"

type PayloadQueue struct {
	Unique_id string `json:"unique_id"`
	Cust_name string `json:"cust_name"`
	AccountId string `json:"account"`
	TenantId  string `json:"tenant_id"`
}

func UnmarshalTelegrampayload(data []byte) (Telegrampayload, error) {
	var r Telegrampayload
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Telegrampayload) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Telegrampayload struct {
	UpdateID      int64                  `json:"update_id"`
	Message       TelegrampayloadMessage `json:"message"`
	CallbackQuery *CallbackQuery         `json:"callback_query"`
}

type CallbackQuery struct {
	ID           string               `json:"id"`
	From         CallbackQueryFrom    `json:"from"`
	Message      CallbackQueryMessage `json:"message"`
	ChatInstance string               `json:"chat_instance"`
	Data         string               `json:"data"`
}

type CallbackQueryFrom struct {
	ID           string  `json:"id"`
	IsBot        bool    `json:"is_bot"`
	FirstName    string  `json:"first_name"`
	LastName     *string `json:"last_name,omitempty"`
	Username     string  `json:"username"`
	LanguageCode string  `json:"language_code"`
}

type CallbackQueryMessage struct {
	MessageID   string      `json:"message_id"`
	From        PurpleFrom  `json:"from"`
	Chat        Chat        `json:"chat"`
	Date        int64       `json:"date"`
	Text        string      `json:"text"`
	ReplyMarkup ReplyMarkup `json:"reply_markup"`
}

type Chat struct {
	ID        string  `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  *string `json:"last_name,omitempty"`
	Username  string  `json:"username"`
	Type      string  `json:"type"`
}

type PurpleFrom struct {
	ID        string `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

type ReplyMarkup struct {
	InlineKeyboard [][]InlineKeyboard `json:"inline_keyboard"`
}

type InlineKeyboard struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

type TelegrampayloadMessage struct {
	MessageID string            `json:"message_id"`
	From      CallbackQueryFrom `json:"from"`
	Chat      Chat              `json:"chat"`
	Date      int64             `json:"date"`
	Text      string            `json:"text"`
	Entities  []Entity          `json:"entities"`
}

type Entity struct {
	Offset int64  `json:"offset"`
	Length int64  `json:"length"`
	Type   string `json:"type"`
}
