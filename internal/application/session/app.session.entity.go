package appSession

import "time"

type Session struct {
	ID           int       `json:"id"`
	SID          string    `json:"sid"`
	TenantID     string    `json:"tenant_id"`
	UniqueId     string    `json:"unique_id"`
	BotPlatform  string    `json:"bot_platform"`
	State        string    `json:"state"`
	Stacktrace   string    `json:"stacktrace"`
	CustMessage  string    `json:"cust_message"`
	CustName     string    `json:"cust_name"`
	CustMessage1 string    `json:"cust_message1"`
	BotResponse  string    `json:"bot_response"`
	OutgoingRe   string    `json:"outgoing_re"`
	BotURL       string    `json:"bot_url"`
	BotAccount   string    `json:"bot_account"`
	ChannelAcc   string    `json:"channel_acc"`
	ChannelSou   string    `json:"channel_sou"`
	ChannelPla   string    `json:"channel_pla"`
	ChannelID    string    `json:"channel_id"`
	Omnichannel  string    `json:"omnichannel"`
	BotDate      time.Time `json:"bot_date"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
