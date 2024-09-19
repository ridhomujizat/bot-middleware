package worker

import (
	"bot-middleware/internal/entities"
)

type MetaData struct {
	BotEndpoint     string                   `json:"bot_endpoint" validate:"omitempty,url"`
	BotAccount      string                   `json:"bot_account" validate:"omitempty"`
	AccountId       string                   `json:"accountId" validate:"required"`
	UniqueId        string                   `json:"unique_id" validate:"required"`
	ChannelId       entities.ChannelId       `json:"channel_id" validate:"required,oneof=12 3 7 5"`
	ChannelPlatform entities.ChannelPlatform `json:"channel_platform" validate:"required,oneof=socioconnect maytapi octopushchat official"`
	ChannelSources  entities.ChannelSources  `json:"channel_sources" validate:"required,oneof=whatsapp fbmessenger livechat telegram"`
	DateTimestamp   string                   `json:"date_timestamp" validate:"required"`
	Sid             string                   `json:"sid,omitempty" validate:"omitempty"`
	NewSession      bool                     `json:"new_session,omitempty" validate:"omitempty"`
	CustMessage     string                   `json:"cust_message,omitempty" validate:"omitempty"`
}
