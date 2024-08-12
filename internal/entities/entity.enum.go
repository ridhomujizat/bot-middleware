package entities

type BotPlatform string
type Omnichannel string
type ChannelSources string
type ChannelId int
type ChannelPlatform string
type MessageType string
type AccountPlatform string
type AccountType string
type BotpressMessageType string
type BotpressErrorType string

const (
	BotpressMessageTypeText            BotpressMessageType = "text"
	BotpressMessageTypeSingleChoice    BotpressMessageType = "single-choice"
	BotpressMessageTypeCarousel        BotpressMessageType = "carousel"
	BotpressMessageTypePostback        BotpressMessageType = "postback"
	BotpressErrorTypeUnauthorizedError BotpressErrorType   = "UnauthorizedError"
)

const (
	BotpressType     AccountType = "botpress"
	WhatsAppType     AccountType = "whatsapp"
	FBMessengerType  AccountType = "fbmessenger"
	OctopushChatType AccountType = "octopushchat"
	IGDMType         AccountType = "igdm"
	TelegramType     AccountType = "telegram"
)

const (
	WhatsAppSocio    AccountPlatform = "whatsapp_socio"
	FBMSocio         AccountPlatform = "fbm_socio"
	WhatsAppMaytapi  AccountPlatform = "whatsapp_maytapi"
	BotpressPlatform AccountPlatform = "botpress"
	LiveChatOctopush AccountPlatform = "livechat_octopushchat"
	IGDMSocio        AccountPlatform = "igdm_socio"
	Telegram         AccountPlatform = "telegram_official"
)

const (
	WHATSAPP    ChannelSources = "whatsapp"
	FBMESSENGER ChannelSources = "fbmessenger"
	LIVECHAT    ChannelSources = "livechat"
	TELEGRAM    ChannelSources = "telegram"
)

const (
	WHATSAPP_ID    ChannelId = 12
	LIVECHAT_ID    ChannelId = 3
	FBMESSENGER_ID ChannelId = 7
	TELEGRAM_ID    ChannelId = 5
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
