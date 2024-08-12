package appSession

import "bot-middleware/internal/entities"

type SessionHistory struct {
	Session

	BotPlatform     entities.BotPlatform     `gorm:"type:session_history_bot_platform_enum;not null;column:bot_platform" json:"bot_platform"`
	ChannelSource   entities.ChannelSources  `gorm:"type:session_history_channel_source_enum;not null;column:channel_source" json:"channel_source"`
	ChannelPlatform entities.ChannelPlatform `gorm:"type:session_history_channel_platform_enum;not null;column:channel_platform" json:"channel_platform"`
	ChannelId       entities.ChannelId       `gorm:"type:session_history_channel_id_enum;not null;column:channel_id" json:"channel_id"`
	Omnichannel     entities.Omnichannel     `gorm:"type:session_history_omnichannel_enum;not null;column:omnichannel" json:"omnichannel"`
}
