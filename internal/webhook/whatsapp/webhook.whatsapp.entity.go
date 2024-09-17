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
	Profile Profile `json:"profile" validate:"omitempty"`
}

type Messages struct {
	Timestamp   string               `json:"timestamp" validate:"required"`
	Type        entities.MessageType `json:"type" validate:"required"`
	From        string               `json:"from" validate:"required"`
	Text        MessageText          `json:"text" validate:"required"`
	Interactive MessageInteractive   `json:"interactive" validate:"required"`
}

type MessageText struct {
	Body string `json:"body" validate:"required"`
}

type MessageInteractive struct {
	Type         string                      `json:"type" validate:"omitempty,oneof='list button'"`
	List_reply   MessageInteractiveListReply `json:"list_reply" validate:"omitempty"`
	Button_reply MessageInteractiveListReply `json:"button_reply" validate:"omitempty"`
}

type MessageInteractiveListReply struct {
	Title string `json:"title" validate:"required"`
	Id    string `json:"id" validate:"required"`
}

type IncomingDTO struct {
	TenantId         string               `json:"tenant_id" validate:"required"`
	AccountId        string               `json:"account_id" validate:"required"`
	TunnelUrl        string               `json:"tunnel_url" validate:"required,url"`
	Contacts         []Contacts           `json:"contacts" validate:"omitempty,dive"`
	Messages         []Messages           `json:"messages" validate:"omitempty,dive"`
	Additional       webhook.AttributeDTO `json:"additional" validate:"omitempty"`
	BotResponse      interface{}          `json:"botResponse,omitempty" validate:"omitempty"`
	OutgoingResponse []interface{}        `json:"outgoingResponse,omitempty" validate:"omitempty"`
}

type EndDTO struct {
	AccountId string `json:"account_id" validate:"required"`
	UniqueId  string `json:"unique_id" validate:"required"`
	Message   string `json:"message,omitempty" validate:"required"`
	SId       string `json:"sid,omitempty" validate:"required"`
}
