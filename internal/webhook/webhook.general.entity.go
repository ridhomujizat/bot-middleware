package webhook

type BotPlatform string
type Omnichannel string
type ChannelSources string
type ChannelID int
type ChannelPlatform string
type MessageType string

const (
	WHATSAPP    ChannelSources = "whatsapp"
	FBMESSENGER ChannelSources = "fbmessenger"
	LIVECHAT    ChannelSources = "livechat"
	TELEGRAM    ChannelSources = "telegram"
)

const (
	WHATSAPP_ID    ChannelID = 12
	LIVECHAT_ID    ChannelID = 3
	FBMESSENGER_ID ChannelID = 7
	TELEGRAM_ID    ChannelID = 5
)

const (
	SOCIOCONNECT ChannelPlatform = "socioconnect"
	MAYTAPI      ChannelPlatform = "maytapi"
	OCTOPUSHCHAT ChannelPlatform = "octopushchat"
	OFFICIAL     ChannelPlatform = "official"
)

const (
	BOTPRESS BotPlatform = "botpress"
)

const (
	ONX Omnichannel = "onx"
	ON5 Omnichannel = "on5"
	ON4 Omnichannel = "on4"
)

const (
	TEXT        MessageType = "text"
	IMAGE       MessageType = "image"
	CONTACTS    MessageType = "contacts"
	DOCUMENT    MessageType = "document"
	INTERACTIVE MessageType = "interactive"
	BUTTON      MessageType = "button"
	LOCATION    MessageType = "location"
	VIDEO       MessageType = "video"
	STICKER     MessageType = "sticker"
	ORDER       MessageType = "order"
	UNKNOWN     MessageType = "unknown"
	VOICE       MessageType = "voice"
	EPHEMERAL   MessageType = "ephemeral"
)

type ParamsDTO struct {
	Omnichannel Omnichannel `json:"omnichannel" validate:"required,oneof=onx on5 on4"`
	TenantId    string      `json:"tenantId" validate:"required"`
	Account     string      `json:"account,omitempty" validate:"omitempty,required"`
	BotPlatform BotPlatform `json:"botplatform" validate:"required,oneof=botpress"`
}

type HandoverDTO struct {
	SID         string      `json:"sid,omitempty" validate:"required"`
	AccountID   string      `json:"account_id" validate:"required"`
	UniqueID    string      `json:"unique_id" validate:"required"`
	Message     string      `json:"message,omitempty" validate:"required"`
	CustMessage interface{} `json:"cust_message,omitempty" validate:"omitempty"`
}

type AttributeDTO struct {
	BotPlatform        BotPlatform     `json:"botplatform" validate:"required,oneof=botpress"`
	Omnichannel        Omnichannel     `json:"omnichannel" validate:"required,oneof=onx on5 on4"`
	TenantId           string          `json:"tenantId" validate:"required"`
	AccountId          string          `json:"accountId" validate:"required"`
	UniqueID           string          `json:"unique_id" validate:"required"`
	ChannelPlatform    ChannelPlatform `json:"channel_platform" validate:"required,oneof=socioconnect maytapi octopushchat official"`
	ChannelSources     ChannelSources  `json:"channel_sources" validate:"required,oneof=whatsapp fbmessenger livechat telegram"`
	ChannelID          ChannelID       `json:"channel_id" validate:"required,oneof=12 3 7 5"`
	DateTimestamp      string          `json:"date_timestamp" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	MiddlewareEndpoint string          `json:"middleware_endpoint" validate:"required,url"`
	CustName           string          `json:"cust_name" validate:"required"`
	StreamID           string          `json:"stream_id,omitempty" validate:"omitempty"`
	CustMessage        interface{}     `json:"cust_message,omitempty" validate:"omitempty"`
}
