package webhookLivechat

import "bot-middleware/internal/webhook"

type User struct {
	RecaptchaResponse string `json:"g-recaptcha-response" validate:"required"`
	Token             string `json:"token" validate:"required"`
}

type Media struct {
	Url      string `json:"url" validate:"required,url"`
	FileName string `json:"fileName" validate:"required"`
	MimeType string `json:"mimeType" validate:"required"`
	FileSize int    `json:"fileSize" validate:"required"`
}

type IncomingDTO struct {
	Tenant           string               `json:"tenant" validate:"required"`
	Account          string               `json:"account" validate:"required"`
	Action           string               `json:"action" validate:"required"`
	DateSend         string               `json:"dateSend" validate:"required"`
	Message          string               `json:"message" validate:"omitempty"`
	MessageOrigin    interface{}          `json:"message_origin" validate:"omitempty"`
	Media            Media                `json:"media" validate:"omitempty"`
	User             User                 `json:"user" validate:"required"`
	Additional       webhook.AttributeDTO `json:"additional" validate:"omitempty"`
	BotResponse      interface{}          `json:"botResponse,omitempty" validate:"omitempty"`
	OutgoingResponse []interface{}        `json:"outgoingResponse,omitempty" validate:"omitempty"`
}

type EndDTO struct {
	AccountId string `json:"account_id" validate:"required"`
	UniqueId  string `json:"unique_id" validate:"required"`
}
