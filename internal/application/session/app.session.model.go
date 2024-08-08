package appSession

import (
	"bot-middleware/internal/entities"
	"time"
)

type Session struct {
	ID               uint                     `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	SID              string                   `gorm:"type:varchar;not null;unique;column:sid" json:"sid"`
	TenantID         string                   `gorm:"type:varchar;not null;column:tenant_id" json:"tenant_id"`
	UniqueID         string                   `gorm:"type:varchar;not null;unique;column:unique_id" json:"unique_id"`
	BotPlatform      entities.BotPlatform     `gorm:"type:session_bot_platform_enum;not null;column:bot_platform" json:"bot_platform"`
	State            string                   `gorm:"type:varchar;column:state" json:"state"`
	Stacktrace       string                   `gorm:"type:text;column:stacktrace" json:"stacktrace"`
	CustMessage      string                   `gorm:"type:text;column:cust_message" json:"cust_message"`
	CustName         string                   `gorm:"type:varchar;column:cust_name" json:"cust_name"`
	CustMessageType  string                   `gorm:"type:varchar;column:cust_message_type" json:"cust_message_type"`
	BotResponse      string                   `gorm:"type:text;column:bot_response" json:"bot_response"`
	OutgoingResponse string                   `gorm:"type:text;column:outgoing_response" json:"outgoing_response"`
	BotURL           string                   `gorm:"type:text;column:bot_url" json:"bot_url"`
	BotAccount       string                   `gorm:"type:varchar;column:bot_account" json:"bot_account"`
	ChannelAccount   string                   `gorm:"type:varchar;column:channel_account" json:"channel_account"`
	ChannelSource    entities.ChannelSources  `gorm:"type:session_channel_source_enum;not null;column:channel_source" json:"channel_source"`
	ChannelPlatform  entities.ChannelPlatform `gorm:"type:session_channel_platform_enum;not null;column:channel_platform" json:"channel_platform"`
	ChannelID        entities.ChannelID       `gorm:"type:session_channel_id_enum;not null;column:channel_id" json:"channel_id"`
	Omnichannel      entities.Omnichannel     `gorm:"type:session_omnichannel_enum;not null;column:omnichannel" json:"omnichannel"`
	BotDate          time.Time                `gorm:"type:timestamp;column:bot_date" json:"bot_date"`
	CreatedAt        time.Time                `gorm:"column:create_at;autoCreateTime;not null" json:"created_at"`
	UpdatedAt        time.Time                `gorm:"column:update_at;autoUpdateTime;not null" json:"updated_at"`
}
