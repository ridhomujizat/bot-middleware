package workerTelegram

import (
	"bot-middleware/internal/application"
	"bot-middleware/internal/application/bot/botpress"
	appSession "bot-middleware/internal/application/session"
	"bot-middleware/internal/entities"
	"bot-middleware/internal/pkg/messaging"
	"bot-middleware/internal/pkg/util"
	webhookTelegram "bot-middleware/internal/webhook/telegram"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pterm/pterm"
)

type TelegramService struct {
	messagingGeneral messaging.MessagingGeneral
	application      *application.Services
}

func NewTelegramService(messagingGeneral messaging.MessagingGeneral, application *application.Services) *TelegramService {
	return &TelegramService{
		messagingGeneral: messagingGeneral,
		application:      application,
	}
}

func (t *TelegramService) Subscribe(exchange, routingKey, queueName string, allowNonJsonMessages bool, handleFunc func([]byte) error) {
	go func() {
		if err := t.messagingGeneral.Subscribe(exchange, routingKey, queueName, allowNonJsonMessages, handleFunc); err != nil {
			util.HandleAppError(err, "subscribe", "Subscribe", true)
		}
	}()
}

func (t *TelegramService) Process(body []byte) error {
	var msg webhookTelegram.IncomingDTO
	if err := json.Unmarshal(body, &msg); err != nil {
		util.HandleAppError(err, "Telegram process", "Unmarshal", false)
		return err
	}
	session, err := t.application.SessionService.FindSession(msg.Additional.UniqueId, string(msg.Additional.ChannelPlatform), string(msg.Additional.ChannelSources), msg.Additional.TenantId)

	if err != nil {
		util.HandleAppError(err, "Telegram process", "FindSession telegram", false)
		// return err
	}
	if session == nil {
		return t.handleNewSession(&msg)
	} else {
		return t.handleExistingSession(&msg, session)
	}
}

func (t *TelegramService) handleNewSession(msg *webhookTelegram.IncomingDTO) error {
	result, err := t.application.BotService.GetAndUpdateBotServer()
	if err != nil {
		util.HandleAppError(err, "Telegram process", "GetAndUpdateBotServer", false)
		return err
	}

	msg.Additional.BotEndpoint = result.ServerAddress
	msg.Additional.BotAccount = result.ServerAccount
	sid, err := util.GenerateId()
	if err != nil {
		pterm.Error.Printfln("Error: %s", err)
	}
	msg.Additional.SID = sid
	msg.Additional.NewSession = true

	additional := msg.Additional
	queueName := fmt.Sprintf("%s:%s:%s:%s:bot", additional.Omnichannel, additional.TenantId, util.GodotEnv("TELEGRAM_QUEUE_NAME"), additional.AccountId)
	pterm.Info.Println("step", queueName)
	return t.messagingGeneral.Publish(queueName, msg)
}

func (t *TelegramService) handleExistingSession(msg *webhookTelegram.IncomingDTO, session *appSession.Session) error {
	msg.Additional.BotEndpoint = session.BotURL
	msg.Additional.BotAccount = session.BotAccount
	msg.Additional.SID = session.SID
	msg.Additional.NewSession = false

	additional := msg.Additional
	queueNameInit := fmt.Sprintf("%s:%s:%s:%s", additional.Omnichannel, additional.TenantId, util.GodotEnv("TELEGRAM_QUEUE_NAME"), additional.AccountId)
	pterm.Info.Printfln("msg => %+v", msg)

	if session.State == "handover" {
		queueName := fmt.Sprintf("%s:handover", queueNameInit)
		return t.messagingGeneral.Publish(queueName, msg)
	}

	queueName := fmt.Sprintf("%s:bot", queueNameInit)
	pterm.Info.Println("step", queueName)
	return t.messagingGeneral.Publish(queueName, msg)
}

func (t *TelegramService) InitiateBot(body []byte) error {
	var msg webhookTelegram.IncomingDTO
	if err := json.Unmarshal(body, &msg); err != nil {
		util.HandleAppError(err, "initiateBot", "Unmarshal", false)
		return err
	}

	switch msg.Additional.BotPlatform {
	case entities.BOTPRESS:
		t.processBotOfficial(body)

	}
	return nil
}

func (t *TelegramService) processBotOfficial(body []byte) error {
	var msg webhookTelegram.IncomingDTO
	if err := json.Unmarshal(body, &msg); err != nil {
		util.HandleAppError(err, "initiateBot", "Unmarshal", false)
		return err
	}
	additional := msg.Additional
	switch additional.BotPlatform {
	case entities.BOTPRESS:
		botPayload, err := t.application.BotService.Botpress.BPTLGOF(&msg)
		if err != nil {
			return err
		}
		if botPayload != nil {
			if err := t.botpressProcess(msg, botPayload); err != nil {
				util.HandleAppError(err, "initiateBot", "processBotOfficial Telegram", false)
				return err
			}
		}
	}
	return nil
}

func (t *TelegramService) botpressProcess(payload webhookTelegram.IncomingDTO, botPayload *botpress.AskPayloadBotpresDTO) error {
	loginResutl, err := t.application.BotService.Botpress.Login("libra_onx", "onx_dev")
	if err != nil {
		util.HandleAppError(err, "login botpress", "Incoming", false)
		return err
	}

	responBot, err := t.application.BotService.Botpress.AskBotpress(payload.Additional.UniqueId, loginResutl.Token, loginResutl.BaseURL, botPayload)
	if err != nil {
		util.HandleAppError(err, "ask botpress", "IncomingHandler", false)
		return err
	}

	botRespon := map[string]interface{}{
		"responses":  responBot.Responses,
		"state":      responBot.State.Context.CurrentNode,
		"stacktrace": responBot.State.Stacktrace,
		"bot_date":   time.Now().Format("2006-01-02 15:04:05"),
	}

	payload.BotResponse = &botRespon

	queueName := fmt.Sprintf("%s:%s:%s:%s:outgoing", payload.Additional.Omnichannel, payload.Additional.TenantId, util.GodotEnv("TELEGRAM_QUEUE_NAME"), payload.Additional.AccountId)
	pterm.Info.Println("step", queueName)
	t.messagingGeneral.Publish(queueName, payload)
	return nil
}

func (t *TelegramService) Outgoing(body []byte) error {
	var msg webhookTelegram.IncomingDTO
	if err := json.Unmarshal(body, &msg); err != nil {
		util.HandleAppError(err, "Telegram process", "Unmarshal", false)
		return err
	}

	var newBotResponse botpress.BotPressResponseDTO
	mapData, err := json.Marshal(msg.BotResponse)
	if err != nil {
		fmt.Println("Error marshaling map to JSON:", err)
		return err
	}
	err = json.Unmarshal(mapData, &newBotResponse)
	if err != nil {
		fmt.Println("Error unmarshaling JSON to BotResponse:", err)
		return err
	}

	var msgResult []map[string]interface{}
	// Assuming newBotResponse.Responses is a slice of responses
	for _, response := range newBotResponse.Responses {
		// Call OutgoingTelegram for each response
		payloadOutgoing := map[string]interface{}{
			"chat_id": msg.Additional.UniqueId,
		}
		switch response.Type {
		case string(botpress.TEXT):
			payloadOutgoing["text"] = response.Text
		case string(botpress.SINGLE_CHOICE):
			InlineKeyboard := []map[string]interface{}{}

			for _, choise := range response.Choices {
				InlineKeyboard = append(InlineKeyboard, map[string]interface{}{
					"text":          choise.Title,
					"callback_data": choise.Value,
				})
			}

			result := [][]map[string]interface{}{}
			result = append(result, InlineKeyboard)

			payloadOutgoing["text"] = response.Text
			payloadOutgoing["reply_markup"] = map[string]interface{}{
				"inline_keyboard": result,
			}
		case string(botpress.CAROUSEL):
			payloadOutgoing["text"] = "Silahkan pilih menu berikut"

			InlineKeyboar := [][]map[string]interface{}{}

			for _, choise := range response.Items {
				action := []map[string]interface{}{}

				for _, c := range choise.Actions {
					action = append(action, map[string]interface{}{
						"subtitile":     choise.SubTitle,
						"image":         choise.Image,
						"callback_data": c.Payload,
						"text":          choise.Title,
					})
				}

				InlineKeyboar = append(InlineKeyboar, action)
			}

			payloadOutgoing["reply_markup"] = map[string]interface{}{
				"inline_keyboard": InlineKeyboar,
			}

		}

		res, err := t.outgoingTelegram(msg.Additional.TenantId, msg.Additional.AccountId, payloadOutgoing)

		if err != nil {
			fmt.Println("Error calling OutgoingTelegram:", err)
			continue
		}

		// Append the result to msgResult
		msgResult = append(msgResult, map[string]interface{}{
			"response":  res,
			"sent_date": time.Now().Format("2006-01-02 15:04:05"),
		})

		for k, v := range payloadOutgoing {
			msgResult[len(msgResult)-1][k] = v
		}
	}

	msg.OutgoingResponse = &msgResult

	queueName := fmt.Sprintf("%s:%s:%s:%s:finish", msg.Additional.Omnichannel, msg.Additional.TenantId, util.GodotEnv("TELEGRAM_QUEUE_NAME"), msg.Additional.AccountId)
	pterm.Info.Println("step", queueName)
	t.messagingGeneral.Publish(queueName, msg)

	return nil
}

func (t *TelegramService) outgoingTelegram(tenantId string, accountId string, payload interface{}) ([]byte, error) {
	account, errAcc := t.application.AccountService.GetAccount(accountId, tenantId)
	if errAcc != nil {
		util.HandleAppError(errAcc, "get user by account id", "OutgoingTelegramText", false)
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		util.HandleAppError(errAcc, "Error marshaling JSON", "OutgoingTelegramText", false)
	}

	baseUrl := fmt.Sprintf("%s/sendMessage", account.BaseURL)

	respon, statusCode, errReq := util.HttpPost(baseUrl, []byte(jsonData), map[string]string{
		"Content-Type": "application/json",
	})
	if errReq != nil {
		util.HandleAppError(errReq, "http post", "OutgoingTelegramText", false)
	}

	fmt.Println("Account:", accountId, "tenant:", tenantId, "URL:", baseUrl, "Payload:", string(jsonData))
	fmt.Println("Response:", string(respon), "Status Code:", statusCode)

	if statusCode == http.StatusOK {
		return []byte(respon), nil
	} else {
		return nil, errReq
	}

}

func (t *TelegramService) Finish(body []byte) error {
	var msg webhookTelegram.IncomingDTO
	if err := json.Unmarshal(body, &msg); err != nil {
		util.HandleAppError(err, "Telegram process", "Unmarshal", false)
	}

	session, err := t.application.SessionService.FindSession(msg.Additional.UniqueId, string(msg.Additional.ChannelPlatform), string(msg.Additional.ChannelSources), msg.Additional.TenantId)
	if err != nil {
		util.HandleAppError(err, "Telegram process", "FindSession telegram", false)
	}

	payload := appSession.Session{}
	payload.TenantID = msg.Additional.TenantId
	payload.UniqueID = msg.Additional.UniqueId
	payload.BotPlatform = msg.Additional.BotPlatform
	payload.BotAccount = msg.Additional.BotAccount
	payload.BotAccount = msg.Additional.BotAccount

	payload.State = (*msg.BotResponse)["state"].(string)
	stacktrace, _ := json.Marshal((*msg.BotResponse)["stacktrace"])
	payload.Stacktrace = string(stacktrace)
	respons, _ := json.Marshal((*msg.BotResponse)["responses"])
	payload.BotResponse = string(respons)
	payload.BotURL = msg.Additional.BotEndpoint
	payload.ChannelSource = msg.Additional.ChannelSources
	payload.ChannelPlatform = msg.Additional.ChannelPlatform
	payload.ChannelID = msg.Additional.ChannelID
	payload.Omnichannel = msg.Additional.Omnichannel
	botDate, _ := time.Parse(time.RFC3339, (*msg.BotResponse)["bot_date"].(string))
	payload.BotDate = botDate

	resOut, _ := json.Marshal(msg.OutgoingResponse)
	payload.OutgoingResponse = string(resOut)
	payload.BotAccount = msg.Additional.BotAccount
	payload.ChannelAccount = msg.Additional.AccountId

	if msg.CallbackQuery != nil {
		payload.CustName = msg.Additional.CustName
		payload.CustMessage = msg.CallbackQuery.Data
		payload.CustMessageType = "postback"
	} else {
		payload.CustName = msg.Additional.CustName
		payload.CustMessage = msg.Message.Text
		payload.CustMessageType = "text"
	}

	if session == nil {
		if err := t.application.SessionService.CreateSession(&payload); err != nil {
			util.HandleAppError(err, "Telegram process", "CreateSession", false)
			return err
		}
	} else {
		fmt.Println("Update session")
		// if err := t.application.SessionService.UpdateSession(&payload); err != nil {
		// 	util.HandleAppError(err, "Telegram process", "UpdateSession", false)
		// 	return err
		// }
	}

	return nil
}
