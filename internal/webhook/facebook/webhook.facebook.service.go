package webhookFacebook

import (
	"bot-middleware/internal/entities"
	"bot-middleware/internal/pkg/messaging"
	"bot-middleware/internal/pkg/util"
	"bot-middleware/internal/webhook"
	"fmt"
	"time"

	"github.com/pterm/pterm"
)

type FacebookService struct {
	messagingGeneral messaging.MessagingGeneral
}

func NewFacebookService(messagingGeneral messaging.MessagingGeneral) *FacebookService {
	return &FacebookService{
		messagingGeneral: messagingGeneral,
	}
}

func (f *FacebookService) Incoming(params webhook.ParamsDTO, payload IncomingDTO) (interface{}, error) {
	queueName := fmt.Sprintf("%s:%s:%s:%s", params.Omnichannel, params.TenantId, util.GodotEnv("FBM_QUEUE_NAME"), payload.Account)
	pterm.Info.Println("queueName", queueName)

	data := webhook.AttributeDTO{
		UniqueId:           payload.Data.Entry[0].Messaging[0].Sender.Id,
		CustName:           fmt.Sprintf("%s %s", payload.Data.Entry[0].Messaging[0].Sender.FirstName, payload.Data.Entry[0].Messaging[0].Sender.LastName),
		BotPlatform:        params.BotPlatform,
		Omnichannel:        params.Omnichannel,
		TenantId:           params.TenantId,
		AccountId:          payload.Account,
		ChannelPlatform:    entities.SOCIOCONNECT,
		ChannelSources:     entities.FBMESSENGER,
		ChannelId:          entities.FBMESSENGER_ID,
		MiddlewareEndpoint: fmt.Sprintf("%s/socioconnect/fbmessenger/%s/%s", util.GodotEnv("BASE_URL"), params.Omnichannel, params.TenantId),
		DateTimestamp:      time.Unix(int64(payload.Data.Entry[0].Messaging[0].Timestamp/1000), 0).Format("2006-01-02 15:04:05"),
		CustMessage:        payload.Data.Entry[0].Messaging[0].Message.Text,
	}

	payload.Additional = data

	err := f.messagingGeneral.Publish(queueName, payload)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (f *FacebookService) Handover(params webhook.ParamsDTO, payload webhook.HandoverDTO) (interface{}, error) {
	queueName := fmt.Sprintf("%s:%s:%s:%s:handover", params.Omnichannel, params.TenantId, util.GodotEnv("FBM_QUEUE_NAME"), payload.AccountId)
	pterm.Info.Println("queueName", queueName)

	additional := map[string]entities.ChannelPlatform{
		"channel_platform": entities.SOCIOCONNECT,
	}

	data := map[string]interface{}{
		"payload":    payload,
		"params":     params,
		"additional": additional,
	}

	pterm.Info.Println("data handover fbm", data)

	err := f.messagingGeneral.Publish(queueName, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (f *FacebookService) End(params webhook.ParamsDTO, payload EndDTO) (interface{}, error) {
	queueName := fmt.Sprintf("%s:%s:%s:%s:end", params.Omnichannel, params.TenantId, util.GodotEnv("FBM_QUEUE_NAME"), payload.AccountId)
	pterm.Info.Println("queueName", queueName)

	additional := map[string]entities.ChannelPlatform{
		"channel_platform": entities.SOCIOCONNECT,
	}

	data := map[string]interface{}{
		"payload":    payload,
		"params":     params,
		"additional": additional,
	}

	pterm.Info.Println("data end fbm", data)

	err := f.messagingGeneral.Publish(queueName, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
