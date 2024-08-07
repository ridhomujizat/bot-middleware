package appBot

import (
	"bot-middleware/internal/application/bot/botpress"

	"gorm.io/gorm"
)

type BotService struct {
	db       *gorm.DB
	Botpress *botpress.BotpressService
}

func NewBotService(db *gorm.DB) *BotService {
	return &BotService{db: db, Botpress: botpress.NewBotpressService(db)}
}

func (a *BotService) GetServerBot(name string) (*ServerBot, error) {
	var server ServerBot
	if err := a.db.Where("server_acco = ?", name).
		First(&server).Error; err != nil {
		return nil, err
	}
	return &server, nil
}
