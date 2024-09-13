package workerWhatsapp

import (
	"bot-middleware/internal/application"
	"bot-middleware/internal/application/bot/botpress"
	appSession "bot-middleware/internal/application/session"
	"bot-middleware/internal/entities"
	"bot-middleware/internal/pkg/libs"
	"bot-middleware/internal/pkg/messaging"
	"bot-middleware/internal/pkg/util"
	"bot-middleware/internal/webhook"
	webhookWhatsapp "bot-middleware/internal/webhook/whatsapp"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pterm/pterm"
)

type WhatsappService struct {
	messagingGeneral messaging.MessagingGeneral
	application      *application.Services
	libsService      *libs.LibsService
}

func NewWhatsappService(messagingGeneral messaging.MessagingGeneral, application *application.Services, libsService *libs.LibsService) *WhatsappService {
	return &WhatsappService{
		messagingGeneral: messagingGeneral,
		application:      application,
		libsService:      libsService,
	}
}

func (l *WhatsappService) Subscribe(exchange, routingKey, queueName string, allowNonJsonMessages bool, handleFunc func([]byte) error) {
	go func() {
		if err := l.messagingGeneral.Subscribe(exchange, routingKey, queueName, allowNonJsonMessages, handleFunc); err != nil {
			util.HandleAppError(err, "subscribe", "Subscribe", true)
		}
	}()
}

func (l *WhatsappService) Process(body []byte) error {
	var msg webhookWhatsapp.IncomingDTO
	if err := json.Unmarshal(body, &msg); err != nil {
		return util.HandleAppError(err, "Livechat process", "Unmarshal", true)
	}

	session, err := l.application.SessionService.FindSession(msg.Additional.UniqueId, string(msg.Additional.ChannelPlatform), string(msg.Additional.ChannelSources), msg.Additional.TenantId)
	if err != nil {
		util.HandleAppError(err, "Whatsapp process", "FindSession", false)
		return err
	}

	pterm.Info.Println("session", session)

	if session == nil {
		return l.handleNewSession(&msg)
	} else {
		return l.handleExistingSession(&msg, session)
	}
}

func (l *WhatsappService) handleNewSession(msg *webhookWhatsapp.IncomingDTO) error {
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

	msg.Additional.Sid = sid
	msg.Additional.NewSession = true
	pterm.Info.Printfln("msg => %+v", msg)
	queueName := fmt.Sprintf("%s:%s:%s:%s:bot", msg.Additional.Omnichannel, msg.Additional.TenantId, util.GodotEnv("WHATSAPP_QUEUE_NAME"), msg.Additional.AccountId)
	return l.messagingGeneral.Publish(queueName, *msg)
}

func (l *WhatsappService) handleExistingSession(msg *webhookWhatsapp.IncomingDTO, session *appSession.Session) error {
	msg.Additional.BotEndpoint = session.BotURL
	msg.Additional.BotAccount = session.BotAccount
	msg.Additional.Sid = session.Sid
	msg.Additional.NewSession = false
	if session.State == "handover" {
		queueName := fmt.Sprintf("%s:%s:%s:%s:handover", msg.Additional.Omnichannel, msg.Additional.TenantId, util.GodotEnv("WHATSAPP_QUEUE_NAME"), msg.Additional.AccountId)
		return l.messagingGeneral.Publish(queueName, *msg)
	}
	queueName := fmt.Sprintf("%s:%s:%s:%s:bot", msg.Additional.Omnichannel, msg.Additional.TenantId, util.GodotEnv("WHATSAPP_QUEUE_NAME"), msg.Additional.AccountId)
	return l.messagingGeneral.Publish(queueName, *msg)
}

func (l *WhatsappService) InitiateBot(body []byte) error {
	var msg webhookWhatsapp.IncomingDTO
	if err := json.Unmarshal(body, &msg); err != nil {
		util.HandleAppError(err, "initiateBot", "Unmarshal", false)
		return err
	}

	if err := l.processBotOctopushchat(&msg); err != nil {
		util.HandleAppError(err, "initiateBot", "processBotOctopushchat", false)
		return err
	}

	return nil
}

func (l *WhatsappService) processBotOctopushchat(msg *webhookWhatsapp.IncomingDTO) error {
	additional := msg.Additional
	switch additional.BotPlatform {
	case entities.BOTPRESS:
		botPayload, err := l.application.BotService.Botpress.BPWASC(*msg)
		if err != nil {
			return err
		}
		if botPayload != nil {
			return l.processBotpress(msg, botPayload)
		}
	}
	return nil
}

func (l *WhatsappService) processBotpress(msg *webhookWhatsapp.IncomingDTO, botPayload *botpress.AskPayloadBotpresDTO) error {
	additional := msg.Additional

	loginResult, err := l.application.BotService.Botpress.Login(additional.BotAccount, additional.TenantId)
	if err != nil {
		util.HandleAppError(err, "processBotpress", "Login", false)
		return err
	}

	result, err := l.application.BotService.Botpress.AskBotpress(additional.UniqueId, loginResult.Token, loginResult.BaseURL, botPayload)
	if err != nil {
		util.HandleAppError(err, "processBotpress", "AskBotpress", false)
		return err
	}

	msg.BotResponse = map[string]interface{}{
		"responses":  result.Responses,
		"state":      result.State.Context.CurrentNode,
		"stacktrace": result.State.Stacktrace,
		"bot_date":   time.Now().Format("2006-01-02 15:04:05"),
	}

	queueName := fmt.Sprintf("%s:%s:%s:%s:outgoing", msg.Additional.Omnichannel, msg.Additional.TenantId, util.GodotEnv("WHATSAPP_QUEUE_NAME"), msg.Additional.AccountId)
	return l.messagingGeneral.Publish(queueName, *msg)
}

func (l *WhatsappService) Outgoing(body []byte) error {
	var msg webhookWhatsapp.IncomingDTO
	if err := json.Unmarshal(body, &msg); err != nil {
		return util.HandleAppError(err, "Livechat Outgoing", "JSON Unmarshal", true)
	}

	additional := msg.Additional

	switch additional.BotPlatform {
	case entities.BOTPRESS:
		switch additional.ChannelPlatform {
		case entities.SOCIOCONNECT:
			return l.processOutgoingLivechatBotpress(&msg)
		}
	}

	return nil
}

func (l *WhatsappService) processOutgoingLivechatBotpress(msg *webhookWhatsapp.IncomingDTO) error {
	additional := msg.Additional
	botResponse := msg.BotResponse

	msg.OutgoingResponse = make([]interface{}, 0)

	for _, outgoing := range botResponse.(map[string]interface{})["responses"].([]botpress.Response) {
		var response interface{}
		var err error
		result := make(map[string]interface{})

		switch outgoing.Type {
		case "text":
			payload := OutgoingTextSocioconnect{
				RecipientType:    "INDIVIDUAL",
				MessagingProduct: "WHATSAPP",
				To:               additional.UniqueId,
				Type:             "TEXT",
				Text:             Text{Body: outgoing.Text},
			}
			response, err = l.libsService.Text(additional.TenantId, additional.AccountId, payload)
		case "single-choice":
			if !outgoing.IsDropdown {
				payload := OutgoingButtonSocioconnect{
					RecipientType:    "INDIVIDUAL",
					MessagingProduct: "WHATSAPP",
					To:               additional.UniqueId,
					Type:             "INTERACTIVE",
					Interactive: Interactive{
						Type: "BUTTON",
						Body: Body{Text: outgoing.Text},
						Action: Action{
							Buttons: mapChoicesToButtons(outgoing.Choices),
						},
					},
				}
				response, err = l.libsService.Button(additional.TenantId, additional.AccountId, payload)
			} else {
				payload := OutgoingListSocioconnect{
					RecipientType:    "INDIVIDUAL",
					MessagingProduct: "WHATSAPP",
					To:               additional.UniqueId,
					Type:             "INTERACTIVE",
					Interactive: Interactive{
						Type: "LIST",
						Body: Body{Text: outgoing.Text},
						Action: Action{
							Button:   outgoing.Text,
							Sections: mapChoicesToSections(outgoing.Choices),
						},
					},
				}
				response, err = l.libsService.Button(additional.TenantId, additional.AccountId, payload)
			}
		case "carousel":
			payload := OutgoingButtonSocioconnect{
				RecipientType:    "INDIVIDUAL",
				MessagingProduct: "WHATSAPP",
				To:               additional.UniqueId,
				Type:             "INTERACTIVE",
				Interactive: Interactive{
					Type: "BUTTON",
					Header: Header{
						Type:  "IMAGE",
						Image: Image{Link: outgoing.Items[0].Image},
					},
					Body: Body{Text: outgoing.Items[0].SubTitle},
					Action: Action{
						Buttons: mapActionsToButtons(outgoing.Items),
					},
				},
			}
			response, err = l.libsService.Button(additional.TenantId, additional.AccountId, payload)

		}
		if err != nil {
			return err
		}

		result["response"] = response
		result["sent_date"] = time.Now().Format("2006-01-02 15:04:05")

		msg.OutgoingResponse = append(msg.OutgoingResponse, result)
	}

	return nil
}

func mapChoicesToButtons(choices []botpress.Choice) []Button {
	var buttons []Button
	for _, c := range choices {
		buttons = append(buttons, Button{Type: "reply", Reply: Reply{Title: c.Title, ID: c.Value}})
	}
	return buttons
}

func mapChoicesToSections(choices []botpress.Choice) []Section {
	var sections []Section
	for _, c := range choices {
		sections = append(sections, Section{Title: c.Title, Description: c.Title, ID: c.Value})
	}
	return sections
}

func mapActionsToButtons(actions []botpress.Carousel) []Button {
	var buttons []Button
	for _, a := range actions {
		buttons = append(buttons, Button{Type: "reply", Reply: Reply{Title: a.Title, ID: a.Actions[0].Payload}})
	}
	return buttons
}

func (l *WhatsappService) Finish(body []byte) error {
	var msg webhookWhatsapp.IncomingDTO
	if err := json.Unmarshal(body, &msg); err != nil {
		return util.HandleAppError(err, "Livechat Finish", "JSON Unmarshal", true)

	}

	sess, err := l.libsService.FindSessionByUniqueId(msg.Additional.UniqueId, msg.Additional.TenantId)
	if err != nil {
		return err
	}

	incomingDate, err := time.Parse(time.RFC3339, msg.Additional.DateTimestamp)
	if err != nil {
		return util.HandleAppError(err, "Livechat Finish", "Timestamp Parse", true)
	}

	botDate, err := time.Parse(time.RFC3339, msg.BotResponse.(map[string]string)["bot_date"])
	if err != nil {
		return util.HandleAppError(err, "Livechat Finish", "BotDate Parse", true)
	}

	session := &appSession.Session{
		Sid:              msg.Additional.Sid,
		TenantId:         msg.Additional.TenantId,
		UniqueId:         msg.Additional.UniqueId,
		BotPlatform:      msg.Additional.BotPlatform,
		State:            fmt.Sprintf("%v", msg.BotResponse.(map[string]interface{})["state"]),
		Stacktrace:       util.JSONstringify(msg.BotResponse.(map[string]interface{})["stacktrace"]),
		BotResponse:      util.JSONstringify(msg.BotResponse.(map[string]interface{})["responses"]),
		BotURL:           msg.Additional.BotEndpoint,
		ChannelSource:    msg.Additional.ChannelSources,
		ChannelPlatform:  msg.Additional.ChannelPlatform,
		ChannelId:        msg.Additional.ChannelId,
		Omnichannel:      msg.Additional.Omnichannel,
		BotDate:          botDate,
		OutgoingResponse: util.JSONstringify(msg.OutgoingResponse),
		BotAccount:       msg.Additional.BotAccount,
		ChannelAccount:   msg.Additional.AccountId,
		IncomingDate:     incomingDate,
	}

	switch msg.Additional.ChannelPlatform {

	case entities.SOCIOCONNECT:
		session.CustMessage = msg.Additional.CustName
		session.CustName = msg.Additional.CustName
		session.CustMessageType = "TEXT"

	}

	if sess == nil {
		if err := l.libsService.CreateSession(session); err != nil {
			return err
		}
	} else {
		session.Id = sess.Id
		if err := l.libsService.UpdateSession(session); err != nil {
			return err
		}
	}

	return nil
}

func (l *WhatsappService) End(body []byte) error {
	var msg map[string]interface{}
	if err := json.Unmarshal(body, &msg); err != nil {
		return util.HandleAppError(err, "Livechat End", "JSON Unmarshal", true)
	}

	payload := msg["payload"].(webhook.EndDTO)
	err := l.libsService.DeleteSessionByUniqueId(&payload)
	if err != nil {
		return err
	}

	return nil
}
