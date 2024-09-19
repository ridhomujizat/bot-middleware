package workerWhatsapp

import (
	"bot-middleware/internal/application/bot/botpress"
	"bot-middleware/internal/entities"
	"bot-middleware/internal/worker"
)

type IncomingDTO struct {
	TenantId  string     `json:"tenant_id" validate:"required"`
	AccountId string     `json:"account_id" validate:"required"`
	TunnelUrl string     `json:"tunnel_url" validate:"required,url"`
	Contacts  []Contacts `json:"contacts" validate:"omitempty,dive"`
	Messages  []Messages `json:"messages" validate:"omitempty,dive"`
}

type Contacts struct {
	WaId    string  `json:"wa_id" validate:"required"`
	Profile Profile `json:"profile" validate:"omitempty"`
}

type Profile struct {
	Name string `json:"name" validate:"required"`
}

type Messages struct {
	Timestamp   string               `json:"timestamp" validate:"required"`
	Type        entities.MessageType `json:"type" validate:"required"`
	From        string               `json:"from" validate:"required"`
	Text        Text                 `json:"text" validate:"required"`
	Interactive InteractiveMessage   `json:"interactive" validate:"required"`
}

type Text struct {
	Body string `json:"body" validate:"required"`
}
type Body struct {
	Text string `json:"text" validate:"required"`
}

type InteractiveMessage struct {
	Type         string               `json:"type" validate:"omitempty,oneof='list button'"`
	List_reply   InteractiveListReply `json:"list_reply" validate:"omitempty"`
	Button_reply InteractiveListReply `json:"button_reply" validate:"omitempty"`
}

type InteractiveListReply struct {
	Title string `json:"title" validate:"required"`
	Id    string `json:"id" validate:"required"`
}

type PayloadDTO struct {
	Incoming         IncomingDTO                      `json:"incoming" validate:"required"`
	MetaData         worker.MetaData                  `json:"metadata" validate:"required"`
	BotResponse      botpress.AnswarPayloadBotpresDTO `json:"bot_response" validate:"omitempty"`
	OutgoingResponse []interface{}                    `json:"outgoing_response" validate:"omitempty"`
}

type OutgoingText struct {
	RecipientType    string `json:"recipient_type"`
	MessagingProduct string `json:"messaging_product"`
	To               string `json:"to"`
	Type             string `json:"type"`
	Text             Text   `json:"text"`
}

type OutgoingButton struct {
	RecipientType    string      `json:"recipient_type"`
	MessagingProduct string      `json:"messaging_product"`
	To               string      `json:"to"`
	Type             string      `json:"type"`
	Interactive      Interactive `json:"interactive"`
}

type Interactive struct {
	Type   string  `json:"type"`
	Body   Body    `json:"body"`
	Action Action  `json:"action"`
	Header *Header `json:"header,omitempty"`
}

type Action struct {
	Buttons  []Button  `json:"buttons,omitempty"`
	Button   string    `json:"button,omitempty"`
	Sections []Section `json:"sections,omitempty"`
}

type Button struct {
	Type  string `json:"type"`
	Reply Reply  `json:"reply"`
}

type Reply struct {
	Title string `json:"title"`
	ID    string `json:"id"`
}

type Section struct {
	Rows  []Rows `json:"rows"`
	Title string `json:"title"`
}
type Rows struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ID          string `json:"id"`
}

type Header struct {
	Type  string `json:"type"`
	Image Image  `json:"image"`
	Text  string `json:"text"`
}

type Image struct {
	ID      string `json:"id"`
	Caption string `json:"caption,omitempty"`
	Link    string `json:"link,omitempty"`
}

type OutgoingList struct {
	RecipientType    string      `json:"recipient_type"`
	MessagingProduct string      `json:"messaging_product"`
	To               string      `json:"to"`
	Type             string      `json:"type"`
	Interactive      Interactive `json:"interactive"`
}

type OutgoingCarousel struct {
	RecipientType    string      `json:"recipient_type"`
	MessagingProduct string      `json:"messaging_product"`
	To               string      `json:"to"`
	Type             string      `json:"type"`
	Interactive      Interactive `json:"interactive"`
}
