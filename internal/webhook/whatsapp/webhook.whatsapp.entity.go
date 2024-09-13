package webhookWhatsapp

import (
	"bot-middleware/internal/entities"
	"bot-middleware/internal/webhook"
)

type Profile struct {
	Name string `json:"name" validate:"required"`
}

type Contacts struct {
	WaId    string  `json:"wa_id" validate:"required"`
	Profile Profile `json:"profile" validate:"omitmepty"`
}

type Messages struct {
	Timestamp   int                  `json:"timestamp" validate:"required"`
	Type        entities.MessageType `json:"type" validate:"required,oneof='text image contacts document interactive button location video sticker order unknown voice ephemeral'"`
	From        string               `json:"from" validate:"required"`
	Text        MessageText          `json:"text" validate:"required"`
	Interactive MessageInteractive   `json:"interactive" validate:"required"`
}

type MessageText struct {
	Body string `json:"body" validate:"required"`
}

type MessageInteractive struct {
	Type         string                      `json:"type" validate:"required,oneof='list button'"`
	List_reply   MessageInteractiveListReply `json:"list_reply" validate:"required"`
	Button_reply MessageInteractiveListReply `json:"button_reply" validate:"required"`
}

type MessageInteractiveListReply struct {
	Title string `json:"title" validate:"required"`
	Id    string `json:"id" validate:"required"`
}

type IncomingDTO struct {
	TenantId         string               `json:"tenant_id" validate:"required"`
	AccountId        string               `json:"account_id" validate:"required"`
	TunnelUrl        string               `json:"tunnel_url" validate:"required,url"`
	Contacts         []Contacts           `json:"contacts" validate:"omitmepty,dive"`
	Messages         []Messages           `json:"messages" validate:"omitmepty,dive"`
	Additional       webhook.AttributeDTO `json:"additional" validate:"omitmepty"`
	BotResponse      interface{}          `json:"botResponse,omitempty" validate:"omitempty"`
	OutgoingResponse []interface{}        `json:"outgoingResponse,omitempty" validate:"omitempty"`
}

type EndDTO struct {
	AccountId string `json:"account_id" validate:"required"`
	UniqueId  string `json:"unique_id" validate:"required"`
	Message   string `json:"message,omitempty" validate:"required"`
	SId       string `json:"sid,omitempty" validate:"required"`
}
