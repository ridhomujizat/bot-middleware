package workerLivechat

import (
	"bot-middleware/internal/application"
	"bot-middleware/internal/application/bot/botpress"
	appSession "bot-middleware/internal/application/session"
	"bot-middleware/internal/pkg/messaging"
	"bot-middleware/internal/pkg/util"
	webhookLivechat "bot-middleware/internal/webhook/livechat"
	"encoding/json"

	"github.com/pterm/pterm"
)

type LivechatService struct {
	messagingGeneral messaging.MessagingGeneral
	application      *application.Services
}

func NewLivechatService(messagingGeneral messaging.MessagingGeneral, application *application.Services) *LivechatService {
	return &LivechatService{
		messagingGeneral: messagingGeneral,
		application:      application,
	}
}

func (l *LivechatService) Subscribe(exchange, routingKey, queueName string, allowNonJsonMessages bool, handleFunc func([]byte) error) {
	go func() {
		if err := l.messagingGeneral.Subscribe(exchange, routingKey, queueName, allowNonJsonMessages, handleFunc); err != nil {
			util.HandleAppError(err, "subscribe", "Subscribe", true)
		}
	}()
}

func (l *LivechatService) Process(body []byte) error {
	var msg webhookLivechat.IncomingDTO
	if err := json.Unmarshal(body, &msg); err != nil {
		util.HandleAppError(err, "Livechat process", "Unmarshal", false)
		return err
	}

	session, err := l.application.SessionService.FindSession(msg.Additional.UniqueId, string(msg.Additional.ChannelPlatform), string(msg.Additional.ChannelSources), msg.Additional.TenantId)
	if err != nil {
		util.HandleAppError(err, "Livechat process", "FindSession", false)
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
		util.HandleAppError(err, "Livechat process", "GetAndUpdateBotServer", false)
		return err
	}

	msg.Additional.BotEndpoint = result.ServerAddress
	msg.Additional.BotAccount = result.ServerAccount
	sid, err := util.GenerateId()
	if err != nil {
		util.HandleAppError(err, "Livechat process", "GenerateId", false)
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

func (l *LivechatService) InitiateBot(body []byte) error {
	var msg webhookLivechat.IncomingDTO
	if err := json.Unmarshal(body, &msg); err != nil {
		util.HandleAppError(err, "initiateBot", "Unmarshal", false)
		return err
	}

	additional := msg.Additional
	switch additional.ChannelPlatform {
	case "OCTOPUSHCHAT":
		if err := l.processBotOctopushchat(&msg); err != nil {
			util.HandleAppError(err, "initiateBot", "processBotOctopushchat", false)
			return err
		}
	default:
		return nil
	}
	return nil
}

func (l *LivechatService) processBotOctopushchat(msg *webhookLivechat.IncomingDTO) error {
	additional := msg.Additional
	switch additional.BotPlatform {
	// case entities.BOTPRESS:
	// 	botPayload, err := l.application.BotService.Botpress.BPTLGOF(*msg)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if botPayload != nil {
	// 		return l.processBotpress(msg, botPayload)
	// 	}
	}
	return nil
}

func (l *LivechatService) processBotpress(msg *webhookLivechat.IncomingDTO, botPayload *botpress.AskPayloadBotpresDTO) error {
	// additional := msg.Additional

	// loginResult, err := l.application.BotService.Botpress.Login(additional.BotAccount, additional.TenantId)
	// if err != nil {
	// 	util.HandleAppError(err, "processBotpress", "Login", false)
	// 	return err
	// }

	// result, err := l.application.BotService.Botpress.AskBotpress(additional.UniqueId, loginResult.Token, loginResult.BaseURL, botPayload)
	// if err != nil {
	// 	util.HandleAppError(err, "processBotpress", "AskBotpress", false)
	// 	return err
	// }

	// msg.Additional.BotResponse = map[string]interface{}{
	// 	"responses":  result.Responses,
	// 	"state":      result.State.Context.CurrentNode,
	// 	"stacktrace": result.State.Stacktrace,
	// 	"bot_date":   time.Now().Format("2006-01-02 15:04:05"),
	// }

	// return l.messagingGeneral.Publish(util.GodotEnv("QUEUE_LIVECHAT_OUTGOING"), *msg)
	return nil
}
