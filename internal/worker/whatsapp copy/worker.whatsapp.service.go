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

	responses, ok := botResponse.(map[string]interface{})["responses"].([]interface{})
	if !ok || responses == nil {
		return fmt.Errorf("invalid type or nil for responses field")
	}

	for _, outgoing := range botResponse.(map[string]interface{})["responses"].([]interface{}) {
		var response interface{}
		var err error
		outgoingMap, ok := outgoing.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid type for outgoing response")
		}

		// Assert the type field to string
		responseType, ok := outgoingMap["type"].(string)
		if !ok {
			return fmt.Errorf("invalid type for response type")
		}

		var result map[string]interface{}

		switch responseType {
		case "text":
			// Assert the text field to string
			text, ok := outgoingMap["text"].(string)
			if !ok {
				return fmt.Errorf("invalid type for text field")
			}
			payload := OutgoingTextSocioconnect{
				RecipientType:    "INDIVIDUAL",
				MessagingProduct: "WHATSAPP",
				To:               additional.UniqueId,
				Type:             "TEXT",
				Text:             Text{Body: text},
			}
			fmt.Println("x", payload)
			// response, err = l.libsService.Text(additional.AccountId, additional.TenantId, payload)
			if err != nil {
				return fmt.Errorf("failed to send text message: %v", err)
			}
		case "single-choice":
			isDropdown, ok := outgoingMap["isDropdown"].(bool)
			if !ok {
				return fmt.Errorf("invalid type for isDropdown field")
			}

			if !isDropdown {
				payload := OutgoingButtonSocioconnect{
					RecipientType:    "INDIVIDUAL",
					MessagingProduct: "WHATSAPP",
					To:               additional.UniqueId,
					Type:             "INTERACTIVE",
					Interactive: Interactive{
						Type: "BUTTON",
						Body: Body{Text: outgoingMap["text"].(string)},
						Action: Action{
							Buttons: mapChoicesToButtons(outgoingMap["choices"].([]botpress.Choice)),
						},
					},
				}
				fmt.Println("payload", payload)

				response, err = l.libsService.Button(additional.AccountId, additional.TenantId, payload)
				if err != nil {
					return fmt.Errorf("error di bukan dropdown")
				}
			} else {
				// Ensure outgoingMap["choices"] is not nil and is of the correct type
				choices, ok := outgoingMap["choices"].([]botpress.Choice)
				if !ok {
					return fmt.Errorf("invalid type or nil for choices field")
				}

				payload := OutgoingListSocioconnect{
					RecipientType:    "INDIVIDUAL",
					MessagingProduct: "WHATSAPP",
					To:               additional.UniqueId,
					Type:             "INTERACTIVE",
					Interactive: Interactive{
						Type: "LIST",
						Body: Body{Text: outgoingMap["text"].(string)},
						Action: Action{
							Button:   outgoingMap["text"].(string),
							Sections: mapChoicesToSections(choices),
						},
					},
				}
				response, err = l.libsService.Button(additional.AccountId, additional.TenantId, payload)
				if err != nil {
					return fmt.Errorf("error di dropdown")
				}
			}
		case "carousel":
			outgoingStruct := outgoing.(map[string]interface{})
			payload := OutgoingButtonSocioconnect{
				RecipientType:    "INDIVIDUAL",
				MessagingProduct: "WHATSAPP",
				To:               additional.UniqueId,
				Type:             "INTERACTIVE",
				Interactive: Interactive{
					Type: "BUTTON",
					Header: Header{
						Type:  "IMAGE",
						Image: Image{Link: outgoingStruct["items"].([]botpress.Carousel)[0].Image},
					},
					Body: Body{Text: outgoingStruct["items"].([]botpress.Carousel)[0].SubTitle},
					Action: Action{
						Buttons: mapActionsToButtons(outgoingStruct["items"].([]botpress.Carousel)),
					},
				},
			}
			response, err = l.libsService.Button(additional.AccountId, additional.TenantId, payload)

		}
		if err != nil {
			return err
		}

		// Initialize the result map
		result = make(map[string]interface{})

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
