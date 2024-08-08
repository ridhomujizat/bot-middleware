package webhook

import 	"bot-middleware/internal/entities"


type HopContext struct {
	AppID    int    `json:"app_id" validate:"required,number"`
	Metadata string `json:"metadata,omitempty" validate:"omitempty"`
}

type Sender struct {
	ID         string `json:"id" validate:"required"`
	FirstName  string `json:"first_name" validate:"required"`
	LastName   string `json:"last_name" validate:"required"`
	ProfilePic string `json:"profile_pic" validate:"required,url"`
}

type Recipient struct {
	ID string `json:"id" validate:"required"`
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
	ID         string       `json:"id" validate:"required"`
	Time       int          `json:"time" validate:"required,number"`
	Messaging  []Messaging  `json:"messaging" validate:"required,dive"`
	HopContext []HopContext `json:"hop_context,omitempty" validate:"omitempty,dive"`
}

type Data struct {
	Object string  `json:"object" validate:"required"`
	Entry  []Entry `json:"entry" validate:"required,dive"`
}

type AttributeDTO struct {
	BotPlatform        entities.BotPlatform     `json:"botplatform" validate:"required,oneof=botpress"`
	Omnichannel        entities.Omnichannel     `json:"omnichannel" validate:"required,oneof=onx on5 on4"`
	TenantId           string          `json:"tenantId" validate:"required"`
	AccountId          string          `json:"accountId" validate:"required"`
	UniqueId           string          `json:"unique_id" validate:"required"`
	ChannelPlatform    entities.ChannelPlatform `json:"channel_platform" validate:"required,oneof=socioconnect maytapi octopushchat official"`
	ChannelSources     entities.ChannelSources  `json:"channel_sources" validate:"required,oneof=whatsapp fbmessenger livechat telegram"`
	ChannelID          entities.ChannelID       `json:"channel_id" validate:"required,oneof=12 3 7 5"`
	DateTimestamp      string          `json:"date_timestamp" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	MiddlewareEndpoint string          `json:"middleware_endpoint" validate:"required,url"`
	CustName           string          `json:"cust_name" validate:"required"`
	StreamID           string          `json:"stream_id,omitempty" validate:"omitempty"`
	CustMessage        interface{}     `json:"cust_message,omitempty" validate:"omitempty"`
	BotEndpoint        string          `json:"bot_endpoint" validate:"omitempty,url"`
	BotAccount         string          `json:"bot_account" validate:"omitempty"`
	SID                string          `json:"sid,omitempty" validate:"omitempty"`
	NewSession         bool            `json:"new_session,omitempty" validate:"omitempty"`
	BotResponse        interface{}     `json:"botResponse,omitempty" validate:"omitempty"`
}

type IncomingDTO struct {
	Tenant      string       `json:"tenant" validate:"required"`
	Channel     string       `json:"channel" validate:"required"`
	Account     string       `json:"account" validate:"required"`
	AccountName string       `json:"account_name" validate:"required"`
	Data        Data         `json:"data" validate:"required"`
	Additional  AttributeDTO `json:"additional" validate:"required"`
	Test        string       `json:"test" validate:"required,oneof=12 3 7 5"`
}

type ParamsDTO struct {
	Omnichannel entities.Omnichannel `json:"omnichannel" validate:"required,oneof=onx on5 on4"`
	TenantId    string      `json:"tenantId" validate:"required"`
	Account     string      `json:"account,omitempty" validate:"omitempty,required"`
	BotPlatform entities.BotPlatform `json:"botplatform" validate:"required,oneof=botpress"`
}

type HandoverDTO struct {
	SID         string      `json:"sid,omitempty" validate:"required"`
	AccountID   string      `json:"account_id" validate:"required"`
	UniqueId    string      `json:"unique_id" validate:"required"`
	Message     string      `json:"message,omitempty" validate:"required"`
	CustMessage interface{} `json:"cust_message,omitempty" validate:"omitempty"`
}

type EndDTO struct {
	AccountID string `json:"account_id" validate:"required"`
	UniqueId  string `json:"unique_id" validate:"required"`
	Message   string `json:"message,omitempty" validate:"required"`
	SID       string `json:"sid,omitempty" validate:"required"`
}
