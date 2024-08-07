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
	Tenant     string               `json:"tenant" validate:"required"`
	Account    string               `json:"account" validate:"required"`
	Action     string               `json:"action" validate:"required"`
	DateSend   string               `json:"dateSend" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	Message    interface{}          `json:"message" validate:"omitempty"`
	Media      Media                `json:"media" validate:"omitempty"`
	User       User                 `json:"user" validate:"required"`
	Additional webhook.AttributeDTO `json:"additional" validate:"required"`
}

type EndDTO struct {
	AccountID string `json:"account_id" validate:"required"`
	UniqueID  string `json:"unique_id" validate:"required"`
}
