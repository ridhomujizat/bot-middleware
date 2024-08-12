package webhookWhatsapp

import (
	"bot-middleware/internal/entities"
	"bot-middleware/internal/pkg/messaging"
	"bot-middleware/internal/pkg/util"
	"bot-middleware/internal/webhook"
	"fmt"
	"time"

	"github.com/pterm/pterm"
)

type WhatsappService struct {
	messagingGeneral messaging.MessagingGeneral
}

func NewWhatsappService(messagingGeneral messaging.MessagingGeneral) *WhatsappService {
	return &WhatsappService{
		messagingGeneral: messagingGeneral,
	}
}

func (w *WhatsappService) Incoming(params webhook.ParamsDTO, payload IncomingDTO) (interface{}, error) {
	queueName := fmt.Sprintf("%s:%s:%s:%s", params.Omnichannel, params.TenantId, util.GodotEnv("WA_QUEUE_NAME"), payload.AccountId)
	pterm.Info.Println("queueName", queueName)
	pterm.Info.Println("payloaddddddddddddddddddddddddddddddddddddd", payload)

	data := webhook.AttributeDTO{
		UniqueId:           payload.Contacts[0].WaId,
		CustName:           payload.Contacts[0].Profile.Name,
		BotPlatform:        params.BotPlatform,
		Omnichannel:        params.Omnichannel,
		TenantId:           params.TenantId,
		AccountId:          payload.AccountId,
		ChannelPlatform:    entities.SOCIOCONNECT,
		ChannelSources:     entities.WHATSAPP,
		ChannelId:          entities.WHATSAPP_ID,
		MiddlewareEndpoint: fmt.Sprintf("%s/socioconnect/whatsapp/%s/%s", util.GodotEnv("BASE_URL"), params.Omnichannel, params.TenantId),
		DateTimestamp:      time.Unix(0, int64(payload.Messages[0].Timestamp)).Format("2006-01-02 15:04:05"),
	}

	payload.Additional = data

	err := w.messagingGeneral.Publish(queueName, payload)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (w *WhatsappService) End(params webhook.ParamsDTO, payload EndDTO) (interface{}, error) {
	queueName := fmt.Sprintf("%s:%s:%s:%s:end", params.Omnichannel, params.TenantId, util.GodotEnv("WA_QUEUE_NAME"), payload.AccountId)
	pterm.Info.Println("queueName", queueName)

	data := map[string]interface{}{
		"payload":    payload,
		"params":     params,
		"additional": map[string]interface{}{"channel_platform": entities.SOCIOCONNECT},
	}

	err := w.messagingGeneral.Publish(queueName, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (w *WhatsappService) Handover(params webhook.ParamsDTO, payload webhook.HandoverDTO) (interface{}, error) {
	queueName := fmt.Sprintf("%s:%s:%s:%s:handover", params.Omnichannel, params.TenantId, "wa_queue_name", payload.AccountId)
	pterm.Info.Println("handover params:", params)
	pterm.Info.Println("handover payload:", payload)
	pterm.Info.Println("handover queueName:", queueName)

	data := map[string]interface{}{
		"payload":    payload,
		"params":     params,
		"additional": map[string]interface{}{"channel_platform": entities.SOCIOCONNECT},
	}

	pterm.Info.Println("handover data", data)

	err := w.messagingGeneral.Publish(queueName, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (w *WhatsappService) Commerce(params webhook.ParamsDTO, payload IncomingDTO) (interface{}, error) {
	uniqueId := "62111"
	custName := "User"
	if len(payload.Contacts) > 0 {
		if payload.Contacts[0].WaId != "" {
			uniqueId = payload.Contacts[0].WaId
		}
		if payload.Contacts[0].Profile.Name != "" {
			custName = payload.Contacts[0].Profile.Name
		}
	}

	data := webhook.AttributeDTO{
		UniqueId:           uniqueId,
		CustName:           custName,
		BotPlatform:        params.BotPlatform,
		Omnichannel:        params.Omnichannel,
		TenantId:           "wa_commerce",
		AccountId:          payload.AccountId,
		ChannelPlatform:    entities.SOCIOCONNECT,
		ChannelSources:     entities.WHATSAPP,
		ChannelId:          entities.WHATSAPP_ID,
		MiddlewareEndpoint: fmt.Sprintf("%s/socioconnect/whatsapp/commerce/%s/%s", util.GodotEnv("BASE_URL"), params.Omnichannel, params.TenantId),
		DateTimestamp:      time.Unix(int64(payload.Messages[0].Timestamp), 0).Format("2006-01-02 15:04:05"),
	}

	if uniqueId == "62111" {
		return "OK", nil
	}

	payload.Additional = data

	pterm.Info.Println("payload", payload)

	queueName := fmt.Sprintf("%s:wa_commerce:%s:%s", params.Omnichannel, util.GodotEnv("WA_QUEUE_NAME"), payload.AccountId)
	if payload.Messages[0].Type == entities.ORDER {
		queueName = fmt.Sprintf("%s:wa_commerce:%s:%s:order_received", params.Omnichannel, util.GodotEnv("WA_QUEUE_NAME"), payload.AccountId)
	}

	pterm.Info.Println("queueName", queueName)

	err := w.messagingGeneral.Publish(queueName, payload)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (w *WhatsappService) Midtrans(params webhook.ParamsDTO, payload interface{}) (interface{}, error) {
	queueName := fmt.Sprintf("%s:wa_commerce:%s:midtrans_notification", params.Omnichannel, util.GodotEnv("WA_QUEUE_NAME"))

	pterm.Info.Println("queueName", queueName)

	err := w.messagingGeneral.Publish(queueName, payload)
	if err != nil {
		return nil, err
	}

	return payload, nil
}
