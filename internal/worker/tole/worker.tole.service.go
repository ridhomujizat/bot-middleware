package workerTole

import (
	"bot-middleware/internal/pkg/messaging"
	"bot-middleware/internal/pkg/util"

	"github.com/pterm/pterm"
)

type ToleService struct {
	messagingGeneral messaging.MessagingGeneral
}

func NewToleService(messagingGeneral messaging.MessagingGeneral, queueName string) {
	service := &ToleService{
		messagingGeneral: messagingGeneral,
	}
	service.subscribe(queueName)
}

func (t *ToleService) subscribe(queueName string) {
	subscriber := t.messagingGeneral.GetSubscriber()

	handleFunc := func(body []byte) {
		pterm.Info.Printfln("Received a message: %s", body)
		// Business Logic
	}

	go func() {
		if err := subscriber.Subscribe(queueName, handleFunc); err != nil {
			util.HandleAppError(err, "subscribe", "Subscribe", true)
		}
	}()
}
