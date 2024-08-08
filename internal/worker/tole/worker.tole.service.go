package workerTole

import (
	"bot-middleware/internal/pkg/messaging"
	"bot-middleware/internal/pkg/util"

	"github.com/pterm/pterm"
)

type ToleService struct {
	messagingGeneral messaging.MessagingGeneral
}

func NewToleService(messagingGeneral messaging.MessagingGeneral, exchange, routingKey, queueName string, allowNonJsonMessages bool) {
	service := &ToleService{
		messagingGeneral: messagingGeneral,
	}
	service.subscribe(exchange, routingKey, queueName, allowNonJsonMessages)
}

func (t *ToleService) subscribe(exchange, routingKey, queueName string, allowNonJsonMessages bool) {
	handleFunc := func(body []byte) error {
		pterm.Info.Printfln("Received a message: %s", body)
		// Business Logic
		return nil
	}

	go func() {
		if err := t.messagingGeneral.Subscribe(exchange, routingKey, queueName, allowNonJsonMessages, handleFunc); err != nil {
			util.HandleAppError(err, "subscribe", "Subscribe", true)
		}
	}()
}
