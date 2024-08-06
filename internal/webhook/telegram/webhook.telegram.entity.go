package webhookTelegram

type EndDTO struct {
	AccountID string `json:"account_id" validate:"required"`
	UniqueID  string `json:"unique_id" validate:"required"`
	Message   string `json:"message,omitempty" validate:"required"`
	SID       string `json:"sid,omitempty" validate:"required"`
}
