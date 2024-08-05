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

func (t *TelegramService) Incoming(params webhook.ParamsDTO, payload webhook.IncomingDTO) (interface{}, error) {

	queueName := fmt.Sprintf("%s:%s:%s:%s", params.Omnichannel, params.TenantId, util.GodotEnv("TELEGRAM_QUEUE_NAME"), payload.Account)
	pterm.Info.Println("queueName", queueName)

	data := webhook.AttributeDTO{
		UniqueID:           payload.Data.Entry[0].Messaging[0].Sender.ID,
		CustName:           fmt.Sprintf("%s %s", payload.Data.Entry[0].Messaging[0].Sender.FirstName, payload.Data.Entry[0].Messaging[0].Sender.LastName),
		BotPlatform:        params.BotPlatform,
		Omnichannel:        params.Omnichannel,
		TenantId:           params.TenantId,
		AccountId:          payload.Account,
		ChannelPlatform:    webhook.OFFICIAL,
		ChannelSources:     webhook.TELEGRAM,
		ChannelID:          webhook.TELEGRAM_ID,
		MiddlewareEndpoint: fmt.Sprintf("%s/official/telegram/%s/%s/%s", util.GodotEnv("BASE_URL"), params.Omnichannel, params.TenantId, params.Account),
		DateTimestamp:      time.Unix(int64(payload.Data.Entry[0].Messaging[0].Timestamp/1000), 0).Format("2006-01-02 15:04:05"),
		CustMessage:        payload.Data.Entry[0].Messaging[0].Message.Text,
	}

	payload.Additional = data

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
