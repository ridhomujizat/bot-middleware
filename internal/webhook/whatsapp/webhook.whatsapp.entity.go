package webhookWhatsapp

import (
	"bot-middleware/internal/webhook"
	"bot-middleware/internal/entities"
)

type Profile struct {
	Name string `json:"name" validate:"required"`
}

type Contacts struct {
	WaId    string  `json:"wa_id" validate:"required"`
	Profile Profile `json:"profile" validate:"omitmepty"`
}

type Messages struct {
	Timestamp int                 `json:"timestamp" validate:"required"`
	Type      entities.MessageType `json:"type" validate:"required,oneof='text image contacts document interactive button location video sticker order unknown voice ephemeral'"`
}

type IncomingDTO struct {
	TenantId   string               `json:"tenant_id" validate:"required"`
	AccountId  string               `json:"account_id" validate:"required"`
	TunnelUrl  string               `json:"tunnel_url" validate:"required,url"`
	Contacts   []Contacts           `json:"contacts" validate:"omitmepty,dive"`
	Messages   []Messages           `json:"messages" validate:"omitmepty,dive"`
	Additional webhook.AttributeDTO `json:"additional" validate:"omitmepty"`
}

type EndDTO struct {
	AccountID string `json:"account_id" validate:"required"`
	UniqueId  string `json:"unique_id" validate:"required"`
	Message   string `json:"message,omitempty" validate:"required"`
	SID       string `json:"sid,omitempty" validate:"required"`
}
