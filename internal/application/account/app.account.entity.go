package appAccount

import "time"

type AccountSetting struct {
	ID         int       `json:"id"`
	Alias      string    `json:"alias"`
	Account    string    `json:"account"`
	TenantId   string    `json:"tenant_id"`
	AccountTyp string    `json:"account_typ"`
	AccountPla string    `json:"account_pla"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	Token      string    `json:"token"`
	DeviceId   string    `json:"device_id"`
	AuthURL    string    `json:"auth_url"`
	BaseURL    string    `json:"base_url"`
	CreatedBy  string    `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
