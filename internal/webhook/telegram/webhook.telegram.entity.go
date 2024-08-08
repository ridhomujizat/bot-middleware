package webhookTelegram

import (
	"bot-middleware/internal/webhook"
	"encoding/json"
)

type EndDTO struct {
	AccountID string `json:"account_id" validate:"required"`
	UniqueId  string `json:"unique_id" validate:"required"`
	Message   string `json:"message,omitempty" validate:"required"`
	SID       string `json:"sid,omitempty" validate:"required"`
}

func UnmarshalTelegramDTO(data []byte) (IncomingTelegramDTO, error) {
	var r IncomingTelegramDTO
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *IncomingTelegramDTO) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type IncomingTelegramDTO struct {
	UpdateID      int64                  `json:"update_id"`
	Message       TelegrampayloadMessage `json:"message"`
	CallbackQuery *CallbackQuery         `json:"callback_query"`
	Additional    *webhook.AttributeDTO  `json:"additional"`
}

type CallbackQuery struct {
	ID           int64                `json:"id"`
	From         CallbackQueryFrom    `json:"from"`
	Message      CallbackQueryMessage `json:"message"`
	ChatInstance string               `json:"chat_instance"`
	Data         string               `json:"data"`
}

type CallbackQueryFrom struct {
	ID           int64   `json:"id"`
	IsBot        bool    `json:"is_bot"`
	FirstName    string  `json:"first_name"`
	LastName     *string `json:"last_name,omitempty"`
	Username     string  `json:"username"`
	LanguageCode string  `json:"language_code"`
}

type CallbackQueryMessage struct {
	MessageID   int64       `json:"message_id"`
	From        PurpleFrom  `json:"from"`
	Chat        Chat        `json:"chat"`
	Date        int64       `json:"date"`
	Text        string      `json:"text"`
	ReplyMarkup ReplyMarkup `json:"reply_markup"`
}

type Chat struct {
	ID        int64   `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  *string `json:"last_name,omitempty"`
	Username  string  `json:"username"`
	Type      string  `json:"type"`
}

type PurpleFrom struct {
	ID        int64  `json:"id"`
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
	MessageID int64             `json:"message_id"`
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
