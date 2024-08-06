package workerTelegram

import (
	"bot-middleware/internal/application"
	"bot-middleware/internal/pkg/messaging"
	"bot-middleware/internal/pkg/util"
	webhookTelegram "bot-middleware/internal/webhook/telegram"
	"fmt"
	"log"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/pterm/pterm"
)

type TelegramService struct {
	messagingGeneral messaging.MessagingGeneral
	application      *application.Services
}

func NewTelegramService(messagingGeneral messaging.MessagingGeneral, application *application.Services, exchange, routingKey, queueName string, allowNonJsonMessages bool) {
	service := &TelegramService{
		messagingGeneral: messagingGeneral,
		application:      application,
	}
	service.subscribe(exchange, routingKey, queueName, allowNonJsonMessages)
}

func (t *TelegramService) subscribe(exchange, routingKey, queueName string, allowNonJsonMessages bool) {
	handleFunc := func(body []byte) {
		switch routingKey {
		case "incoming":
			fmt.Println("Initializing Telegram service...")
			t.incomingHandler(body)
		case "bot-process":
			fmt.Println("Starting bot service...")
			t.botHandler(body)
		default:
			fmt.Println("Unknown routing")
		}
	}

	go func() {
		if err := t.messagingGeneral.Subscribe(exchange, routingKey, queueName, allowNonJsonMessages, handleFunc); err != nil {
			util.HandleAppError(err, "subscribe", "Subscribe", true)
		}
	}()
}

func (t *TelegramService) incomingHandler(body []byte) {

	payload, errBody := webhookTelegram.UnmarshalTelegramDTO(body)
	if errBody != nil {
		util.HandleAppError(errBody, "unmarshal telegram dto", "IncomingHandler", true)
	}

	var additional = payload.Additional

	session, errSession := t.application.SessinonService.FindSession(additional.UniqueID, string(additional.ChannelPlatform), string(additional.ChannelSources), additional.TenantId)
	if errSession != nil {
		botserver, errBot := t.application.BotService.GetServerBot("libra_onx")
		if errBot != nil {
			pterm.Error.Printfln("Error: %s", errBot)
		}

		sid, err := gonanoid.New()
		if err != nil {
			log.Fatalf("Failed to generate NanoID: %v", err)
		}

		payload.Additional.BotEndpoint = botserver.ServerAddr
		payload.Additional.BotAccount = botserver.ServerAcco
		payload.Additional.SID = sid
		payload.Additional.NewSession = true
	} else {
		payload.Additional.BotEndpoint = session.BotURL
		payload.Additional.BotAccount = session.BotAccount
		payload.Additional.SID = session.SID
		payload.Additional.NewSession = false
		pterm.Info.Printfln("payload: %v", payload.Additional)
	}

	queueName := fmt.Sprintf("%s:%s:%s:%s:bot", additional.Omnichannel, additional.TenantId, util.GodotEnv("TELEGRAM_QUEUE_NAME"), additional.AccountId)
	t.messagingGeneral.Publish(queueName, payload)

}

func (t *TelegramService) botHandler(body []byte) {

}
