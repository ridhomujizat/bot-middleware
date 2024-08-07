package workerTelegram

import (
	"bot-middleware/internal/application"
	"bot-middleware/internal/pkg/messaging"
	"bot-middleware/internal/pkg/util"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pterm/pterm"
)

type TelegramOutgoingHandler struct {
	messagingGeneral messaging.MessagingGeneral
	application      *application.Services
}

func NewTelegramOutgoingHandler(messagingGeneral messaging.MessagingGeneral, application *application.Services, exchange, routingKey, queueName string, allowNonJsonMessages bool) {
	service := &TelegramOutgoingHandler{
		messagingGeneral: messagingGeneral,
		application:      application,
	}
	service.subscribe(exchange, routingKey, queueName, allowNonJsonMessages)
}

func (t *TelegramOutgoingHandler) subscribe(exchange, routingKey, queueName string, allowNonJsonMessages bool) {
	handleFunc := func(body []byte) {
		t.OutgoingHandler(body)
	}

	go func() {
		if err := t.messagingGeneral.Subscribe(exchange, routingKey, queueName, allowNonJsonMessages, handleFunc); err != nil {
			util.HandleAppError(err, "subscribe", "Subscribe", false)
		}
	}()
}

func (t *TelegramOutgoingHandler) OutgoingHandler(body []byte) {
	pterm.Info.Printfln("Received a message%s", body)

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
