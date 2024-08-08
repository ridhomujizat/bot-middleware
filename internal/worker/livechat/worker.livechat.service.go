package workerLivechat

import (
	"bot-middleware/internal/application"
	appSession "bot-middleware/internal/application/session"
	"bot-middleware/internal/pkg/messaging"
	"bot-middleware/internal/pkg/util"
	webhookLivechat "bot-middleware/internal/webhook/livechat"
	"encoding/json"

	"github.com/pterm/pterm"
	"github.com/streadway/amqp"
)

type LivechatService struct {
	messagingGeneral messaging.MessagingGeneral
	application      *application.Services
}

func NewLivechatService(messagingGeneral messaging.MessagingGeneral, application *application.Services, exchange, routingKey, queueName string, allowNonJsonMessages bool) {
	service := &LivechatService{
		messagingGeneral: messagingGeneral,
		application:      application,
	}
	service.subscribe(exchange, routingKey, queueName, allowNonJsonMessages)
}

func (l *LivechatService) subscribe(exchange, routingKey, queueName string, allowNonJsonMessages bool) {
	handleFunc := func(body []byte, delivery amqp.Delivery) {
		err := l.process(body)
		if err != nil {
			pterm.Error.Printf("Failed to process message: %v. Requeuing...\n", err)
			delivery.Nack(false, true)
		} else {
			delivery.Ack(false)
		}
	}

	go func() {
		if err := l.messagingGeneral.Subscribe(exchange, routingKey, queueName, allowNonJsonMessages, handleFunc); err != nil {
			util.HandleAppError(err, "subscribe", "Subscribe", true)
		}
	}()
}

func (l *LivechatService) process(body []byte) error {
	var msg webhookLivechat.IncomingDTO
	if err := json.Unmarshal(body, &msg); err != nil {
		util.HandleAppError(err, "Livechat process", "Unmarshal", true)
		return err
	}

	session, err := l.application.SessionService.FindSession(msg.Additional.UniqueId, string(msg.Additional.ChannelPlatform), string(msg.Additional.ChannelSources), msg.Additional.TenantId)
	if err != nil {
		util.HandleAppError(err, "Livechat process", "FindSession", true)
		return err
	}

	pterm.Info.Println("session", session)

	if session == nil {
		return l.handleNewSession(&msg)
	} else {
		return l.handleExistingSession(&msg, session)
	}
}

func (l *LivechatService) handleNewSession(msg *webhookLivechat.IncomingDTO) error {
	result, err := l.application.BotService.GetAndUpdateBotServer()
	if err != nil {
		util.HandleAppError(err, "Livechat process", "GetAndUpdateBotServer", true)
		return err
	}

	msg.Additional.BotEndpoint = result.ServerAddress
	msg.Additional.BotAccount = result.ServerAccount
	sid, err := util.GenerateId()
	if err != nil {
		util.HandleAppError(err, "Livechat process", "GenerateId", true)
		return err
	}

	msg.Additional.SID = sid
	msg.Additional.NewSession = true
	pterm.Info.Printfln("msg => %+v", msg)
	return l.messagingGeneral.Publish(util.GodotEnv("QUEUE_LIVECHAT_BOT"), *msg)
}

func (l *LivechatService) handleExistingSession(msg *webhookLivechat.IncomingDTO, session *appSession.Session) error {
	msg.Additional.BotEndpoint = session.BotURL
	msg.Additional.BotAccount = session.BotAccount
	msg.Additional.SID = session.SID
	msg.Additional.NewSession = false
	if session.State == "handover" {
		return l.messagingGeneral.Publish(util.GodotEnv("QUEUE_LIVECHAT_FORWARD"), *msg)
	}
	return l.messagingGeneral.Publish(util.GodotEnv("QUEUE_LIVECHAT_BOT"), *msg)
}
