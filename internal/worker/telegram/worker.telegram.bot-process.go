package workerTelegram

import (
	"bot-middleware/internal/application"
	"bot-middleware/internal/entities"
	"bot-middleware/internal/pkg/messaging"
	"bot-middleware/internal/pkg/util"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/streadway/amqp"
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
	handleFunc := func(body []byte, delivery amqp.Delivery) {
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

	responBot, errAsk := t.application.BotService.Botpress.AskBotpress(payload.Additional.UniqueId, loginResutl.Token, loginResutl.BaseURL, botPayload)
	if errAsk != nil {
		util.HandleAppError(errAsk, "ask botpress", "IncomingHandler", false)
	}

	botRespon := map[string]interface{}{
		"responses":  responBot.Responses,
		"state":      responBot.State.Context.CurrentNode,
		"stacktrace": responBot.State.Stacktrace,
		"bot_date":   time.Now().Format("2006-01-02 15:04:05"),
	}

	payload.BotResponse = &botRespon

	// END BOTPRESS ====================================
	queueName := fmt.Sprintf("%s:%s:%s:%s:outgoing", payload.Additional.Omnichannel, payload.Additional.TenantId, util.GodotEnv("TELEGRAM_QUEUE_NAME"), payload.Additional.AccountId)
	t.messagingGeneral.Publish(queueName, payload)

}

func (t *TelegramBotProcess) OutgoingTelegram(tenantId string, accountId string, payload interface{}) ([]byte, error) {
	account, errAcc := t.application.AccountService.GetUserByAccountId(accountId)
	if errAcc != nil {
		util.HandleAppError(errAcc, "get user by account id", "OutgoingTelegramText", false)
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	}

	baseUrl := fmt.Sprintf("%s/sendMessage", account.BaseURL)
	respon, statusCode, errReq := util.HttpPost(baseUrl, []byte(jsonData), map[string]string{})
	if errReq != nil {
		util.HandleAppError(errReq, "http post", "OutgoingTelegramText", false)
	}

	if statusCode == http.StatusOK {
		return []byte(respon), nil
	} else {
		return nil, errReq
	}

}
