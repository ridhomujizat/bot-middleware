package appBot

import (
	"time"
)

type ServerBot struct {
	ID            uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	ServerAccount string    `gorm:"type:varchar;not null;column:server_account" json:"server_account"`
	ServerName    string    `gorm:"type:varchar;not null;column:server_name" json:"server_name"`
	ServerAddress string    `gorm:"type:varchar;not null;column:server_address" json:"server_address"`
	IsActive      bool      `gorm:"type:bool;not null;default:true;column:is_active" json:"is_active"`
	Total         int       `gorm:"type:int;not null;default:0;column:total" json:"total"`
	CreatedAt     time.Time `gorm:"column:create_at;autoCreateTime;not null" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:update_at;autoUpdateTime;not null" json:"updated_at"`
}
