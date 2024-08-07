package webhookTelegram

import (
	"bot-middleware/internal/pkg/messaging"
	"bot-middleware/internal/pkg/util"
	"bot-middleware/internal/webhook"
	"fmt"
	"time"

	"github.com/pterm/pterm"
)

type TelegramService struct {
	messagingGeneral messaging.MessagingGeneral
}

func NewTelegramService(messagingGeneral messaging.MessagingGeneral) *TelegramService {
	return &TelegramService{
		messagingGeneral: messagingGeneral,
	}
}

func (t *TelegramService) Incoming(params webhook.ParamsDTO, payload IncomingTelegramDTO) (interface{}, error) {

	queueName := fmt.Sprintf("%s:%s:%s:%s", params.Omnichannel, params.TenantId, util.GodotEnv("TELEGRAM_QUEUE_NAME"), params.Account)
	pterm.Info.Println("queueName", queueName)

	data := webhook.AttributeDTO{
		// UniqueID:           payload.Data.Entry[0].Messaging[0].Sender.ID,
		BotPlatform:        params.BotPlatform,
		Omnichannel:        params.Omnichannel,
		TenantId:           params.TenantId,
		AccountId:          params.Account,
		ChannelPlatform:    webhook.OFFICIAL,
		ChannelSources:     webhook.TELEGRAM,
		ChannelID:          webhook.TELEGRAM_ID,
		MiddlewareEndpoint: fmt.Sprintf("%s/official/telegram/%s/%s/%s", util.GodotEnv("BASE_URL"), params.Omnichannel, params.TenantId, params.Account),
		// DateTimestamp:      time.Unix(int64(payload.Data.Entry[0].Messaging[0].Timestamp/1000), 0).Format("2006-01-02 15:04:05"),
		// CustMessage:        payload.Data.Entry[0].Messaging[0].Message.Text,
	}

	if payload.CallbackQuery != nil {
		var lastName string
		if payload.Message.From.LastName != nil {
			lastName = *payload.Message.From.LastName
		} else {
			lastName = "" // Provide a default value or handle it appropriately
		}
		custName := fmt.Sprintf("%s %s", payload.CallbackQuery.From.FirstName, lastName)
		UniqueID := fmt.Sprintf("%d", payload.CallbackQuery.From.ID)
		data.CustName = custName
		data.UniqueID = UniqueID
		data.CustMessage = payload.CallbackQuery.Message.Text
		data.DateTimestamp = time.Unix(int64(payload.CallbackQuery.Message.Date), 0).Format("2006-01-02 15:04:05")
	} else {
		var lastName string
		if payload.Message.From.LastName != nil {
			lastName = *payload.Message.From.LastName
		} else {
			lastName = "" // Provide a default value or handle it appropriately
		}
		custName := fmt.Sprintf("%s %s", payload.Message.From.FirstName, lastName)
		UniqueID := fmt.Sprintf("%d", payload.Message.From.ID)
		data.UniqueID = UniqueID
		data.CustName = custName
		data.CustMessage = payload.Message.Text
		data.DateTimestamp = time.Unix(int64(payload.Message.Date), 0).Format("2006-01-02 15:04:05")
	}

	payload.Additional = &data

	pterm.Info.Println("payload", payload)
	pterm.Info.Println("queueName", queueName)

	err := t.messagingGeneral.Publish(queueName, payload)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (t *TelegramService) Handover(params webhook.ParamsDTO, payload webhook.HandoverDTO) (interface{}, error) {
	queueName := fmt.Sprintf("%s:%s:%s:%s:handover", params.Omnichannel, params.TenantId, util.GodotEnv("TELEGRAM_QUEUE_NAME"), payload.AccountID)
	pterm.Info.Println("queueName", queueName)

	additional := map[string]webhook.ChannelPlatform{
		"channel_platform": webhook.OFFICIAL,
	}

	data := map[string]interface{}{
		"payload":    payload,
		"params":     params,
		"additional": additional,
	}

	pterm.Info.Println("data handover telegram", data)

	err := t.messagingGeneral.Publish(queueName, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (t *TelegramService) End(params webhook.ParamsDTO, payload webhook.EndDTO) (interface{}, error) {
	queueName := fmt.Sprintf("%s:%s:%s:%s:end", params.Omnichannel, params.TenantId, util.GodotEnv("TELEGRAM_QUEUE_NAME"), payload.AccountID)
	pterm.Info.Println("queueName", queueName)

	additional := map[string]webhook.ChannelPlatform{
		"channel_platform": webhook.OFFICIAL,
	}

	data := map[string]interface{}{
		"payload":    payload,
		"params":     params,
		"additional": additional,
	}

	pterm.Info.Println("data====>", data)

	err := t.messagingGeneral.Publish(queueName, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
