package workerLivechat

import (
	"bot-middleware/internal/application"
	"bot-middleware/internal/application/bot/botpress"
	appSession "bot-middleware/internal/application/session"
	"bot-middleware/internal/entities"
	"bot-middleware/internal/pkg/libs"
	"bot-middleware/internal/pkg/messaging"
	"bot-middleware/internal/pkg/util"
	"bot-middleware/internal/webhook"
	webhookLivechat "bot-middleware/internal/webhook/livechat"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pterm/pterm"
)

type LivechatService struct {
	messagingGeneral messaging.MessagingGeneral
	application      *application.Services
	libsService      *libs.LibsService
}

func NewLivechatService(messagingGeneral messaging.MessagingGeneral, application *application.Services, libsService *libs.LibsService) *LivechatService {
	return &LivechatService{
		messagingGeneral: messagingGeneral,
		application:      application,
		libsService:      libsService,
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
		return util.HandleAppError(err, "Livechat process", "Unmarshal", true)
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

	msg.Additional.Sid = sid
	msg.Additional.NewSession = true
	pterm.Info.Printfln("msg => %+v", msg)
	return l.messagingGeneral.Publish(util.GodotEnv("QUEUE_LIVECHAT_BOT"), *msg)
}

func (l *LivechatService) handleExistingSession(msg *webhookLivechat.IncomingDTO, session *appSession.Session) error {
	msg.Additional.BotEndpoint = session.BotURL
	msg.Additional.BotAccount = session.BotAccount
	msg.Additional.Sid = session.Sid
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

	if err := l.processBotOctopushchat(&msg); err != nil {
		util.HandleAppError(err, "initiateBot", "processBotOctopushchat", false)
		return err
	}

	return nil
}

func (l *LivechatService) processBotOctopushchat(msg *webhookLivechat.IncomingDTO) error {
	additional := msg.Additional
	switch additional.BotPlatform {
	case entities.BOTPRESS:
		// botPayload, err := l.application.BotService.Botpress.BPLCOC(*msg)
		// if err != nil {
		// 	return err
		// }
		// if botPayload != nil {
		// 	return l.processBotpress(msg, botPayload)
		// }
	}
	return nil
}

func (l *LivechatService) processBotpress(msg *webhookLivechat.IncomingDTO, botPayload *botpress.AskPayloadBotpresDTO) error {
	additional := msg.Additional

	loginResult, err := l.application.BotService.Botpress.Login(additional.BotAccount, additional.TenantId)
	if err != nil {
		util.HandleAppError(err, "processBotpress", "Login", false)
		return err
	}

	result, err := l.application.BotService.Botpress.AskBotpress(additional.UniqueId, loginResult.Token, loginResult.BaseURL, botPayload)
	if err != nil {
		return err
	}

	msg.BotResponse = map[string]interface{}{
		"responses":  result.Responses,
		"state":      result.State.Context.CurrentNode,
		"stacktrace": result.State.Stacktrace,
		"bot_date":   time.Now().Format("2006-01-02 15:04:05"),
	}

	return l.messagingGeneral.Publish(util.GodotEnv("QUEUE_LIVECHAT_OUTGOING"), *msg)
}

func (l *LivechatService) Outgoing(body []byte) error {
	var msg webhookLivechat.IncomingDTO
	if err := json.Unmarshal(body, &msg); err != nil {
		return util.HandleAppError(err, "Livechat Outgoing", "JSON Unmarshal", true)
	}

	additional := msg.Additional

	switch additional.BotPlatform {
	case entities.BOTPRESS:
		switch additional.ChannelPlatform {
		case entities.OCTOPUSHCHAT:
			return l.processOutgoingLivechatBotpress(&msg)
		}
	}

	return nil
}

func (l *LivechatService) processOutgoingLivechatBotpress(msg *webhookLivechat.IncomingDTO) error {
	additional := msg.Additional
	botResponse := msg.BotResponse

	msg.OutgoingResponse = make([]interface{}, 0)

	for _, outgoing := range botResponse.(map[string]interface{})["responses"].([]BotResponse) {
		var response interface{}
		var err error
		result := make(map[string]interface{})

		switch outgoing.Type {
		case entities.BotpressMessageTypeText:
			payload := Outgoing{
				FromName: "Quitina",
				Message:  outgoing.Text,
				Token:    additional.UniqueId,
			}

			response, err = l.libsService.Text(additional.TenantId, additional.AccountId, payload)

			result["fromName"] = payload.FromName
			result["message"] = payload.Message
			result["token"] = payload.Token

		case entities.BotpressMessageTypeSingleChoice:
			payload := OutgoingButton{
				FromName: "Quitina",
				Message: MessageButton{
					Title:  outgoing.Text,
					Button: l.mapChoices(outgoing.Choices),
				},
				Token: additional.UniqueId,
			}

			response, err = l.libsService.Button(additional.TenantId, additional.AccountId, payload)

			result["fromName"] = payload.FromName
			result["message"] = payload.Message
			result["token"] = payload.Token

		case entities.BotpressMessageTypeCarousel:
			payload := OutgoingCarousel{
				FromName: "Quitina",
				Message: MessageCarousel{
					MessageType: "carousel",
					Slider:      l.mapCarouselItems(outgoing.Items),
				},
				Token: additional.UniqueId,
			}

			response, err = l.libsService.Carousel(additional.TenantId, additional.AccountId, payload)

			result["fromName"] = payload.FromName
			result["message"] = payload.Message
			result["token"] = payload.Token
		}

		if err != nil {
			return err
		}

		result["response"] = response
		result["sent_date"] = time.Now().Format("2006-01-02 15:04:05")

		msg.OutgoingResponse = append(msg.OutgoingResponse, result)
	}

	return l.messagingGeneral.Publish(util.GodotEnv("QUEUE_LIVECHAT_FINISH"), *msg)
}

func (l *LivechatService) mapChoices(choices []Choices) []Button {
	buttons := make([]Button, len(choices))
	for i, c := range choices {
		buttons[i] = Button{
			Label: c.Title,
			Value: c.Value,
		}
	}
	return buttons
}

func (l *LivechatService) mapCarouselItems(items []Carousel) []MessageCarouselSlider {
	sliders := make([]MessageCarouselSlider, len(items))
	for i, c := range items {
		sliders[i] = MessageCarouselSlider{
			Title:    c.Title,
			Subtitle: c.Subtitle,
			Picture:  c.Image,
			Menu:     l.mapActions(c.Actions),
		}
	}
	return sliders
}

func (l *LivechatService) mapActions(actions []CarouselActions) []MessageCarouselSliderMenu {
	menus := make([]MessageCarouselSliderMenu, len(actions))
	for i, a := range actions {
		menus[i] = MessageCarouselSliderMenu{
			Label: a.Title,
			Value: a.Payload,
		}
	}
	return menus
}

func (l *LivechatService) End(body []byte) error {
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

func (l *LivechatService) Handover(body []byte) error {
	var msg map[string]interface{}
	if err := json.Unmarshal(body, &msg); err != nil {
		return util.HandleAppError(err, "Livechat Handover", "JSON Unmarshal", true)
	}

	payload := msg["payload"].(webhook.HandoverDTO)
	sid := payload.Sid
	session, err := l.libsService.FindSessionBySid(sid)
	if err != nil {
		return err
	}

	if session != nil {
		session.State = "handover"
		if err := l.libsService.UpdateSession(session); err != nil {
			return err
		}

		data, ok := payload.CustMessage.(map[string]interface{})
		if !ok {
			return util.HandleAppError(err, "Livechat Handover", "Message origin", true)
		}

		messageOrigin, ok := data["message_origin"].(map[string]interface{})
		if !ok {
			return util.HandleAppError(err, "Livechat Handover", "Message origin", true)
		}

		token, ok := messageOrigin["token"].(string)
		if !ok {
			return util.HandleAppError(err, "Livechat Handover", "Message origin token", true)
		}
		pterm.Info.Printfln("data!@# HANDOVER: %+v", data)

		payloadToSend := map[string]string{
			"email":       payload.Email,
			"username":    payload.Name,
			"mobilePhone": "+628888888888",
			"token":       token,
		}

		payloadData, err := json.Marshal(payloadToSend)
		if err != nil {
			return util.HandleAppError(err, "Livechat Handover", "JSON Marsjal", true)
		}

		res, _, errRespon := util.HttpPost("https://be-livechat.on5.co.id/onx_quitline/handover/createSession", []byte(payloadData),
			map[string]string{
				"Content-Type": "application/json",
			})

		if errRespon != nil {
			return errRespon
		}

		pterm.Success.Printfln("responseHandover livechat: %+v", res)
	}

	return nil
}

func (l *LivechatService) Finish(body []byte) error {
	var msg webhookLivechat.IncomingDTO
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
	case entities.OCTOPUSHCHAT:
		switch msg.Action {
		case "clientReplyText":
			session.CustMessage = msg.Message
			session.CustName = ""
			session.CustMessageType = "text"
		case "clientReplyButton":
			session.CustMessage = msg.Message
			session.CustName = ""
			session.CustMessageType = "button"
		}
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

func (l *LivechatService) Forward(body []byte) error {
	var msg webhookLivechat.IncomingDTO
	if err := json.Unmarshal(body, &msg); err != nil {
		return util.HandleAppError(err, "Livechat Forward", "JSON Unmarshal", true)
	}

	var url string
	switch msg.Action {
	case "clientReplyText":
		url = "https://be-livechat.on5.co.id/onx_quitline/handover/reply/text"
	case "clientReplyMedia":
		url = "https://be-livechat.on5.co.id/onx_quitline/handover/reply/media"
	case "clientReplyLocation":
		url = "https://be-livechat.on5.co.id/onx_quitline/handover/reply/location"
	default:
		return nil
	}

	payloadBytes, err := json.Marshal(msg.MessageOrigin)
	if err != nil {
		return util.HandleAppError(err, "Livechat Forward", "JSON Marshal", true)
	}

	res, _, _ := util.HttpPost(url, payloadBytes,
		map[string]string{
			"Content-Type": "application/json",
		})

	pterm.Info.Printfln("responseHandover: %+v", res)

	return nil
}
