package workerTelegram

import (
	"bot-middleware/internal/application"
	"bot-middleware/internal/entities"
	"bot-middleware/internal/pkg/messaging"
	"bot-middleware/internal/pkg/util"
	"fmt"
)

type TelegramBotProcess struct {
	messagingGeneral messaging.MessagingGeneral
	application      *application.Services
}

func NewTelegramBotProcess(messagingGeneral messaging.MessagingGeneral, application *application.Services, exchange, routingKey, queueName string, allowNonJsonMessages bool) {
	service := &TelegramBotProcess{
		messagingGeneral: messagingGeneral,
		application:      application,
	}
	service.subscribe(exchange, routingKey, queueName, allowNonJsonMessages)
}

func (t *TelegramBotProcess) subscribe(exchange, routingKey, queueName string, allowNonJsonMessages bool) {
	handleFunc := func(body []byte) {
		t.botProcess(body)
	}

	go func() {
		if err := t.messagingGeneral.Subscribe(exchange, routingKey, queueName, allowNonJsonMessages, handleFunc); err != nil {
			util.HandleAppError(err, "subscribe", "Subscribe", false)
		}
	}()
}

func (t *TelegramBotProcess) botProcess(body []byte) {
	payload, errBody := entities.UnmarshalTelegramDTO(body)
	if errBody != nil {
		util.HandleAppError(errBody, "unmarshal telegram dto", "IncomingHandler", false)
	}

	// // BOTPRESS ========================================
	loginResutl, loginErr := t.application.BotService.Botpress.Login()
	if loginErr != nil {
		util.HandleAppError(loginErr, "login botpress", "Incoming", false)
	}

	botPayload, botPayloadErr := t.application.BotService.Botpress.ParsingPayloadTelegram(payload)
	if botPayloadErr != nil {
		util.HandleAppError(botPayloadErr, "parsing payload telegram", "IncomingHandler", false)
	}

	botRespon, errAsk := t.application.BotService.Botpress.AskBotpress(payload.Additional.UniqueID, loginResutl.Token, loginResutl.BaseURL, botPayload)
	if errAsk != nil {
		util.HandleAppError(errAsk, "ask botpress", "IncomingHandler", false)
	}
	queueName := fmt.Sprintf("%s:%s:%s:%s:outgoing", payload.Additional.Omnichannel, payload.Additional.TenantId, util.GodotEnv("TELEGRAM_QUEUE_NAME"), payload.Additional.AccountId)

	// END BOTPRESS ====================================
	t.messagingGeneral.Publish(queueName, botRespon)

	// login, err := t.application.BotService.Botpress.Login()
	// if err != nil {
	// 	pterm.Error.Printfln("Error: %s", err)
	// }
	// fmt.Println("Login: ", login)

}
