package appBot

import (
	"gorm.io/gorm"
)

type BotService struct {
	db *gorm.DB
}

func NewBotService(db *gorm.DB) *BotService {
	return &BotService{db: db}
}

func (a *BotService) GetServerBot(name string) (*ServerBot, error) {
	var server ServerBot
	if err := a.db.Where("server_acco = ?", name).
		First(&server).Error; err != nil {
		return nil, err
	}
	return &server, nil
}
