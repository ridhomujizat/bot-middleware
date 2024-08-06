package webHookTelegram

import (
	"bot-middleware/internal/pkg/messaging"
)

type TelegramService struct {
	messagingGeneral messaging.MessagingGeneral
}

func NewTelegramService(messagingGeneral messaging.MessagingGeneral) *TelegramService {
	return &TelegramService{
		messagingGeneral: messagingGeneral,
	}
}

func (t *TelegramService) SendQueue(queueName string, message interface{}) error {
	publisher := t.messagingGeneral.GetPublisher()
	return publisher.Publish(queueName+":telegram", message)
}
