package application

import (
	appAccount "bot-middleware/internal/application/account"
)

type Services struct {
	AccountService *appAccount.AccountService
	// SessionService *session.SessionService
}
