package appBot

import (
	"time"

	"gorm.io/gorm"
)

type ServerBot struct {
	ID            int       `json:"id"`
	ServerAccount string    `json:"server_account"`
	ServerName    string    `json:"server_name"`
	ServerAddress string    `json:"server_address"`
	IsActive      string    `json:"is_active"`
	Total         int       `json:"total"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (serverbot *ServerBot) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	serverbot.CreatedAt = now
	serverbot.UpdatedAt = now
	return
}

func (serverbot *ServerBot) BeforeUpdate(tx *gorm.DB) (err error) {
	serverbot.UpdatedAt = time.Now()
	return
}
