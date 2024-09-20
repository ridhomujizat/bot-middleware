package workerWhatsapp

import (
	"bot-middleware/internal/application"
	"bot-middleware/internal/application/bot/botpress"
	appSession "bot-middleware/internal/application/session"
	"bot-middleware/internal/entities"
	"bot-middleware/internal/pkg/libs"
	"bot-middleware/internal/pkg/messaging"
	"bot-middleware/internal/pkg/util"
	"bot-middleware/internal/worker"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
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
	var msg IncomingDTO

	if err := json.Unmarshal(body, &msg); err != nil {
		return util.HandleAppError(err, "Unmarshal process", "Unmarshal", false)
	}

	result, err := l.application.BotService.GetAndUpdateBotServer()
	if err != nil {
		util.HandleAppError(err, "WA process", "GetAndUpdateBotServer", false)
		return err
	}

	session, err := l.application.SessionService.InitAndCheckSession(msg.AccountId, string(entities.SOCIOCONNECT), string(entities.WHATSAPP), msg.TenantId)
	if err != nil {
		util.HandleAppError(err, "Whatsapp process", "FindSession", false)
		return err
	}

	payload := PayloadDTO{
		Incoming: msg,
		MetaData: worker.MetaData{
			BotEndpoint:     result.ServerAddress,
			BotAccount:      result.ServerAccount,
			AccountId:       msg.AccountId,
			UniqueId:        msg.Contacts[0].WaId,
			ChannelPlatform: entities.SOCIOCONNECT,
			ChannelSources:  entities.WHATSAPP,
			ChannelId:       entities.WHATSAPP_ID,
			DateTimestamp:   msg.Messages[0].Timestamp,
			Sid:             session.Sid,
			NewSession:      session.NewSession,
		},
	}

	queueName := fmt.Sprintf("%s:%s:%s:bot", msg.TenantId, util.GodotEnv("WHATSAPP_QUEUE_NAME"), msg.AccountId)
	return l.messagingGeneral.Publish(queueName, payload)

}

func (l *WhatsappService) ProcessBot(body []byte) error {
	var payload PayloadDTO

	if err := json.Unmarshal(body, &payload); err != nil {
		util.HandleAppError(err, "initiateBot", "Unmarshal", false)
		return err
	}

	botPayload, err := l.ParsingIncoming(*&payload)
	if err != nil {
		return err
	}
	refreshToken := false
	loginResult, err := l.application.BotService.Botpress.Login(payload.MetaData.BotAccount, payload.Incoming.TenantId, &refreshToken)
	if err != nil {
		util.HandleAppError(err, "processBotpress", "Login", false)
		return err
	}

	result, err := l.application.BotService.Botpress.AskBotpress(payload.MetaData.UniqueId, loginResult.Token, loginResult.BaseURL, botPayload, &botpress.RefreshToken{
		BotAccount: payload.MetaData.BotAccount,
		TenantId:   payload.Incoming.TenantId,
	})
	if err != nil {
		util.HandleAppError(err, "processBotpress", "AskBotpress", false)
		return err
	}

	payload.BotResponse = botpress.AnswarPayloadBotpresDTO{
		Responses:  result.Responses,
		State:      result.State.Context.CurrentNode,
		Stacktrace: result.State.Stacktrace,
		BotDate:    time.Now().Format("2006-01-02 15:04:05"),
	}

	queueName := fmt.Sprintf("%s:%s:%s:outgoing", payload.Incoming.TenantId, util.GodotEnv("WHATSAPP_QUEUE_NAME"), payload.MetaData.AccountId)
	return l.messagingGeneral.Publish(queueName, payload)

}

func (a *WhatsappService) ParsingIncoming(payload PayloadDTO) (*botpress.AskPayloadBotpresDTO, error) {
	var botPayload botpress.AskPayloadBotpresDTO

	switch payload.Incoming.Messages[0].Type {
	case "text":
		botPayload.Type = botpress.BotpressMessageType(botpress.TEXT)
		botPayload.Text = payload.Incoming.Messages[0].Text.Body
		botPayload.Metadata = payload.MetaData
		botPayload.Metadata.CustMessage = botPayload.Text

	case "interactive":
		switch payload.Incoming.Messages[0].Interactive.Type {
		case "list_reply":
			botPayload.Type = botpress.BotpressMessageType(botpress.SINGLE_CHOICE)
			botPayload.Text = payload.Incoming.Messages[0].Interactive.List_reply.Id
			botPayload.Metadata = payload.MetaData
			botPayload.Metadata.CustMessage = botPayload.Text
		case "button_reply":
			botPayload.Type = botpress.BotpressMessageType(botpress.SINGLE_CHOICE)
			botPayload.Text = payload.Incoming.Messages[0].Interactive.Button_reply.Id
			botPayload.Metadata = payload.MetaData
			botPayload.Metadata.CustMessage = botPayload.Text
		default:
			return nil, fmt.Errorf("unsupported interactive type: %s", payload.Incoming.Messages[0].Interactive.Type)
		}

	default:
		return nil, fmt.Errorf("unsupported action: %s", payload.Incoming.Messages[0].Type)
	}
	return &botPayload, nil

}

func (l *WhatsappService) ProcessOutgoing(body []byte) error {
	var payload PayloadDTO
	if err := json.Unmarshal(body, &payload); err != nil {
		util.HandleAppError(err, "initiateBot", "Unmarshal", false)
		return err
	}
	MetaData := payload.MetaData
	botResponse := payload.BotResponse

	payload.OutgoingResponse = make([]interface{}, 0)

	for _, outgoing := range botResponse.Responses {
		outPayload := OutgoingText{
			RecipientType:    "INDIVIDUAL",
			MessagingProduct: "WHATSAPP",
			To:               "6285156824825",
			Type:             "TEXT",
			Text:             Text{Body: "Hello, this is a test message"},
		}
		l.libsService.Text(MetaData.AccountId, payload.Incoming.TenantId, outPayload)
		var response interface{}
		var err error
		result := make(map[string]interface{})

		switch outgoing.Type {
		case string(botpress.TEXT):
			if err != nil {
				return fmt.Errorf("failed to send text message: %v", err)
			}
			result["payload"] = outPayload
			result["response"] = response
			result["sent_date"] = time.Now().Format("2006-01-02 15:04:05")
		case string(botpress.SINGLE_CHOICE):

			if !outgoing.IsDropdown {
				outPayload := OutgoingButton{
					RecipientType:    "INDIVIDUAL",
					MessagingProduct: "WHATSAPP",
					To:               MetaData.UniqueId,
					Type:             "INTERACTIVE",
					Interactive: Interactive{
						Type: "button",
						Body: Body{Text: outgoing.Text},
						Action: Action{
							Buttons: mapChoicesToButtons(outgoing.Choices),
						},
					},
				}
				response, err = l.libsService.Text(MetaData.AccountId, payload.Incoming.TenantId, outPayload)
				if err != nil {
					return fmt.Errorf("failed to send button message: %v", err)
				}
				result["payload"] = outPayload
				result["response"] = response
				result["sent_date"] = time.Now().Format("2006-01-02 15:04:05")
			} else {
				Button := outgoing.DropdownPlaceholder

				if Button == "" {
					Button = "Select an option"
				}
				outPayload := OutgoingList{
					RecipientType:    "INDIVIDUAL",
					MessagingProduct: "WHATSAPP",
					To:               MetaData.UniqueId,
					Type:             "INTERACTIVE",
					Interactive: Interactive{
						Type: "list",
						Body: Body{Text: outgoing.Text},
						Action: Action{
							Button: Button,
							Sections: []Section{{
								Title: Button,
								Rows:  mapChoicesToSections(outgoing.Choices)}},
						},
					},
				}
				response, err = l.libsService.Text(MetaData.AccountId, payload.Incoming.TenantId, outPayload)
				if err != nil {
					return fmt.Errorf("failed to send button message: %v", err)
				}
				result["payload"] = outPayload
				result["response"] = response
				result["sent_date"] = time.Now().Format("2006-01-02 15:04:05")
			}
		case string(botpress.CAROUSEL):
			Button := outgoing.DropdownPlaceholder

			if Button == "" {
				Button = "Select an option"
			}
			outPayload := OutgoingCarousel{
				RecipientType:    "INDIVIDUAL",
				MessagingProduct: "WHATSAPP",
				To:               MetaData.UniqueId,
				Type:             "INTERACTIVE",
				Interactive: Interactive{
					Type: "carousel",
					Body: Body{Text: outgoing.Text},
					Action: Action{
						Button:  Button,
						Buttons: mapChoicesToButtons(outgoing.Choices),
					},
				},
			}
			response, err = l.libsService.Text(MetaData.AccountId, payload.Incoming.TenantId, outPayload)
			if err != nil {
				return fmt.Errorf("failed to send button message: %v", err)
			}

			result["payload"] = outPayload
			result["response"] = response
			result["sent_date"] = time.Now().Format("2006-01-02 15:04:05")

		}

		payload.OutgoingResponse = append(payload.OutgoingResponse, result)
	}
	queueName := fmt.Sprintf("%s:%s:%s:finish", payload.Incoming.TenantId, util.GodotEnv("WHATSAPP_QUEUE_NAME"), payload.MetaData.AccountId)
	return l.messagingGeneral.Publish(queueName, payload)
}

func mapChoicesToButtons(choices []botpress.Choice) []Button {
	var buttons []Button
	for _, c := range choices {
		buttons = append(buttons, Button{Type: "reply", Reply: Reply{Title: c.Title, ID: c.Value}})
	}
	return buttons
}

func mapChoicesToSections(choices []botpress.Choice) []Rows {
	var sections []Rows
	for _, c := range choices {
		parts := strings.Split(c.Title, "|")
		title := strings.TrimSpace(parts[0])
		description := ""
		if len(parts) > 1 {
			description = strings.TrimSpace(parts[1])
			if len(description) > 71 {
				description = description[:71]
			}
		}
		sections = append(sections, Rows{Title: title, Description: description, ID: c.Value})
	}
	return sections
}

func (l *WhatsappService) ProcessFinish(body []byte) error {
	var payload PayloadDTO
	if err := json.Unmarshal(body, &payload); err != nil {
		util.HandleAppError(err, "Process Finish WA", "Unmarshal", false)
		return err
	}

	BotResponse := payload.BotResponse
	MetaData := payload.MetaData
	Incoming := payload.Incoming

	sess, err := l.libsService.FindSessionByUniqueId(MetaData.UniqueId, Incoming.TenantId)
	if err != nil {
		util.HandleAppError(err, "WA  Finish", "Get Sess Parse", false)
		return err
	}

	incomingTimestamp, err := strconv.ParseInt(MetaData.DateTimestamp, 10, 64)
	if err != nil {
		util.HandleAppError(err, "WA Finish", "Timestamp Parse", false)
		return err
	}
	incomingDate := time.Unix(incomingTimestamp, 0)

	// Parse BotResponse.BotDate with the correct layout
	botDate, err := time.Parse("2006-01-02 15:04:05", BotResponse.BotDate)
	if err != nil {
		util.HandleAppError(err, "WA Finish", "BotDate Parse", false)
		return err
	}

	session := &appSession.Session{
		Sid:              MetaData.Sid,
		TenantId:         Incoming.TenantId,
		UniqueId:         MetaData.UniqueId,
		BotPlatform:      "botpress",
		State:            fmt.Sprintf("%v", BotResponse.State),
		Stacktrace:       util.JSONstringify(BotResponse.Stacktrace),
		BotResponse:      util.JSONstringify(BotResponse.Responses),
		BotURL:           MetaData.BotEndpoint,
		ChannelSource:    MetaData.ChannelSources,
		ChannelPlatform:  MetaData.ChannelPlatform,
		ChannelId:        MetaData.ChannelId,
		Omnichannel:      "onx",
		BotDate:          botDate,
		OutgoingResponse: util.JSONstringify(payload.OutgoingResponse),
		BotAccount:       MetaData.BotAccount,
		ChannelAccount:   MetaData.AccountId,
		IncomingDate:     incomingDate,
	}

	session.CustMessage = Incoming.Messages[0].Text.Body
	session.CustName = Incoming.Contacts[0].Profile.Name
	session.CustMessageType = string(Incoming.Messages[0].Type)

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
