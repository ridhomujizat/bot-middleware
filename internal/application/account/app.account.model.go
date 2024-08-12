package appAccount

import (
	"bot-middleware/internal/entities"
	"time"
)

type AccountSetting struct {
	Id              uint                     `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Alias           string                   `gorm:"type:varchar;column:alias" json:"alias"`
	Account         string                   `gorm:"type:varchar;column:account" json:"account"`
	TenantId        string                   `gorm:"type:varchar;column:tenant_id;not null" json:"tenant_id"`
	AccountType     entities.AccountType     `gorm:"type:account_setting_account_type_enum;column:account_type;not null" json:"account_type"`
	AccountPlatform entities.AccountPlatform `gorm:"type:account_setting_account_platform_enum;column:account_platform;not null" json:"account_platform"`
	Username        string                   `gorm:"type:varchar;column:username" json:"username"`
	Password        string                   `gorm:"type:varchar;column:password" json:"password"`
	Token           string                   `gorm:"type:text;column:token" json:"token"`
	DeviceId        string                   `gorm:"type:varchar;column:device_id" json:"device_id"`
	AuthURL         string                   `gorm:"type:text;column:auth_url" json:"auth_url"`
	BaseURL         string                   `gorm:"type:text;column:base_url" json:"base_url"`
	CreatedAt       time.Time                `gorm:"column:create_at;autoCreateTime;not null" json:"created_at"`
	UpdatedAt       time.Time                `gorm:"column:update_at;autoUpdateTime;not null" json:"updated_at"`
}
