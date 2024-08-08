package webhookLivechat

import (
	"bot-middleware/internal/pkg/messaging"
	"bot-middleware/internal/pkg/util"
	"bot-middleware/internal/webhook"
	"fmt"

	"github.com/pterm/pterm"
)

type LivechatService struct {
	messagingGeneral messaging.MessagingGeneral
}

func NewLivechatService(messagingGeneral messaging.MessagingGeneral) *LivechatService {
	return &LivechatService{
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

func (l *LivechatService) Incoming(params webhook.ParamsDTO, payload IncomingDTO) (interface{}, error) {
	queueName := fmt.Sprintf("%s:%s:%s:%s", params.Omnichannel, params.TenantId, util.GodotEnv("LIVECHAT_QUEUE_NAME"), params.Account)
	pterm.Info.Println("queueName", queueName)

	data := webhook.AttributeDTO{
		UniqueId:           payload.User.Token,
		CustName:           "",
		BotPlatform:        params.BotPlatform,
		Omnichannel:        params.Omnichannel,
		TenantId:           params.TenantId,
		AccountId:          payload.Account,
		ChannelPlatform:    webhook.SOCIOCONNECT,
		ChannelSources:     webhook.LIVECHAT,
		ChannelID:          webhook.LIVECHAT_ID,
		MiddlewareEndpoint: fmt.Sprintf("%s/octopushchat/livechat/%s/%s", util.GodotEnv("BASE_URL"), params.Omnichannel, params.TenantId),
		DateTimestamp:      payload.DateSend,
		CustMessage:        payload,
	}

	payload.Additional = data

	pterm.Info.Println("Payload:", payload)

	err := l.messagingGeneral.Publish(queueName, payload)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (l *LivechatService) Handover(params webhook.ParamsDTO, payload webhook.HandoverDTO) (interface{}, error) {
	queueName := fmt.Sprintf("%s:%s:%s:%s:handover", params.Omnichannel, params.TenantId, util.GodotEnv("LIVECHAT_QUEUE_NAME"), payload.AccountID)
	pterm.Info.Println("queueName", queueName)

	additional := map[string]webhook.ChannelPlatform{
		"channel_platform": webhook.SOCIOCONNECT,
	}

	data := map[string]interface{}{
		"payload":    payload,
		"params":     params,
		"additional": additional,
	}

	pterm.Info.Println("data handover livechat", data)

	err := l.messagingGeneral.Publish(queueName, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (l *LivechatService) End(params webhook.ParamsDTO, payload EndDTO) (interface{}, error) {
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

	err := l.messagingGeneral.Publish(queueName, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
