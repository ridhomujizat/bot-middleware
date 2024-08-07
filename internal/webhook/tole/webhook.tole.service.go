package webhookTole

import (
	"bot-middleware/internal/pkg/messaging"
)

type ToleService struct {
	messagingGeneral messaging.MessagingGeneral
}

func NewToleService(messagingGeneral messaging.MessagingGeneral) *ToleService {
	return &ToleService{
		messagingGeneral: messagingGeneral,
	}
}

func (t *ToleService) Send(queueName string, message interface{}) error {
	return t.messagingGeneral.Publish(queueName+":tole", message)
}
