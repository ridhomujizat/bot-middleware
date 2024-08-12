package webhookFacebook

import "bot-middleware/internal/webhook"

type HopContext struct {
	AppId    int    `json:"app_id" validate:"required,number"`
	Metadata string `json:"metadata,omitempty" validate:"omitempty"`
}

type Sender struct {
	Id         string `json:"id" validate:"required"`
	FirstName  string `json:"first_name" validate:"required"`
	LastName   string `json:"last_name" validate:"required"`
	ProfilePic string `json:"profile_pic" validate:"required,url"`
}

type Recipient struct {
	Id string `json:"id" validate:"required"`
}

type QuickReply struct {
	Payload string `json:"payload" validate:"required"`
}

type AttachmentPayload struct {
	URL string `json:"url" validate:"required,url"`
}

type Attachments struct {
	Type    string            `json:"type" validate:"required"`
	Payload AttachmentPayload `json:"payload" validate:"required"`
}

type Message struct {
	Mid         string        `json:"mid" validate:"required"`
	Text        string        `json:"text,omitempty" validate:"omitempty"`
	QuickReply  QuickReply    `json:"quick_reply,omitempty" validate:"omitempty"`
	Attachments []Attachments `json:"attachments,omitempty" validate:"omitempty,dive"`
}

type MessagePostback struct {
	Title   string `json:"title" validate:"required"`
	Payload string `json:"payload" validate:"required"`
}

type Messaging struct {
	Sender    Sender           `json:"sender" validate:"required"`
	Recipient Recipient        `json:"recipient" validate:"required"`
	Timestamp int              `json:"timestamp" validate:"required,number"`
	Message   *Message         `json:"message,omitempty" validate:"omitempty"`
	Postback  *MessagePostback `json:"postback,omitempty" validate:"omitempty"`
}

type Entry struct {
	Id         string       `json:"id" validate:"required"`
	Time       int          `json:"time" validate:"required,number"`
	Messaging  []Messaging  `json:"messaging" validate:"required,dive"`
	HopContext []HopContext `json:"hop_context,omitempty" validate:"omitempty,dive"`
}

type Data struct {
	Object string  `json:"object" validate:"required"`
	Entry  []Entry `json:"entry" validate:"required,dive"`
}

type IncomingDTO struct {
	Tenant      string               `json:"tenant" validate:"required"`
	Channel     string               `json:"channel" validate:"required"`
	Account     string               `json:"account" validate:"required"`
	AccountName string               `json:"account_name" validate:"required"`
	Data        Data                 `json:"data" validate:"required"`
	Additional  webhook.AttributeDTO `json:"additional" validate:"required"`
}

type EndDTO struct {
	AccountId string `json:"account_id" validate:"required"`
	UniqueId  string `json:"unique_id" validate:"required"`
	Message   string `json:"message,omitempty" validate:"required"`
	SId       string `json:"sid,omitempty" validate:"required"`
}
