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

func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

func getFloat64(m map[string]interface{}, key string) float64 {
	if val, ok := m[key].(float64); ok {
		return val
	}
	return 0
}

func (t *TelegramService) Incoming(params webhook.ParamsDTO, payload map[string]interface{}) (interface{}, error) {
	queueName := fmt.Sprintf("%s:%s:%s:%s", params.Omnichannel, params.TenantId, util.GodotEnv("TELEGRAM_QUEUE_NAME"), params.Account)
	pterm.Info.Println("queueName", queueName)

	var data webhook.AttributeDTO
	lastName := ""

	if callbackQuery, ok := payload["callback_query"].(map[string]interface{}); ok {
		from := callbackQuery["from"].(map[string]interface{})
		lastName = getString(from, "last_name")
		data = webhook.AttributeDTO{
			UniqueID:           getString(from, "id"),
			CustName:           fmt.Sprintf("%s %s", getString(from, "first_name"), lastName),
			BotPlatform:        params.BotPlatform,
			Omnichannel:        params.Omnichannel,
			TenantId:           params.TenantId,
			AccountId:          params.Account,
			ChannelPlatform:    webhook.OFFICIAL,
			ChannelSources:     webhook.TELEGRAM,
			ChannelID:          webhook.TELEGRAM_ID,
			MiddlewareEndpoint: fmt.Sprintf("%s/official/telegram/%s/%s/%s", util.GodotEnv("BASE_URL"), params.Omnichannel, params.TenantId, params.Account),
			DateTimestamp:      time.Unix(int64(getFloat64(callbackQuery["message"].(map[string]interface{}), "date")), 0).Format("2006-01-02 15:04:05"),
			CustMessage:        payload,
		}
	} else if message, ok := payload["message"].(map[string]interface{}); ok {
		from := message["from"].(map[string]interface{})
		lastName = getString(from, "last_name")
		data = webhook.AttributeDTO{
			UniqueID:           getString(from, "id"),
			CustName:           fmt.Sprintf("%s %s", getString(from, "first_name"), lastName),
			BotPlatform:        params.BotPlatform,
			Omnichannel:        params.Omnichannel,
			TenantId:           params.TenantId,
			AccountId:          params.Account,
			ChannelPlatform:    webhook.OFFICIAL,
			ChannelSources:     webhook.TELEGRAM,
			ChannelID:          webhook.TELEGRAM_ID,
			MiddlewareEndpoint: fmt.Sprintf("%s/official/telegram/%s/%s/%s", util.GodotEnv("BASE_URL"), params.Omnichannel, params.TenantId, params.Account),
			DateTimestamp:      time.Unix(int64(getFloat64(message, "date")), 0).Format("2006-01-02 15:04:05"),
			CustMessage:        payload,
		}
	}

	payload["additional"] = data

	pterm.Info.Println("Payload:", payload)

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

func (t *TelegramService) End(params webhook.ParamsDTO, payload EndDTO) (interface{}, error) {
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
