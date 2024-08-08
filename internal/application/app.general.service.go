package application

import (
	appAccount "bot-middleware/internal/application/account"
	appBot "bot-middleware/internal/application/bot"
	appSession "bot-middleware/internal/application/session"
)

type Services struct {
	AccountService *appAccount.AccountService
	SessionService *appSession.SessionService
	BotService     *appBot.BotService
	// SessionService *session.SessionService
}
