package appBot

import "time"

type ServerBot struct {
	ID         int       `json:"id"`
	ServerAcco string    `json:"server_acco"`
	ServerName string    `json:"server_name"`
	ServerAddr string    `json:"server_addr"`
	IsActive   string    `json:"is_active"`
	Total      int       `json:"total"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
