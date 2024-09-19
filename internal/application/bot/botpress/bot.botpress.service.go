package botpress

import (
	appAccount "bot-middleware/internal/application/account"
	"bot-middleware/internal/pkg/util"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/pterm/pterm"
	"gorm.io/gorm"
)

type BotpressService struct {
	db             *gorm.DB
	accountService *appAccount.AccountService
}

func NewBotpressService(db *gorm.DB) *BotpressService {
	return &BotpressService{db: db, accountService: appAccount.NewAccountService(db)}
}

type loginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *BotpressService) Login(botAccount, tenantId string, refreshToken *bool) (*LoginRespon, error) {
	account, err := a.accountService.GetAccount(botAccount, tenantId)
	if err != nil {
		return nil, err
	}
	// Set default value for refreshToken if it is nil
	if refreshToken == nil {
		defaultRefreshToken := false
		refreshToken = &defaultRefreshToken
	}

	today := time.Now()
	dayUpdate := account.UpdatedAt
	diffDays := today.Sub(dayUpdate).Hours() / 24

	if diffDays >= 1 || account.Token == "" || *refreshToken {
		payloadLogin := loginPayload{
			Email:    account.Username,
			Password: account.Password,
		}
		jsonData, err := json.Marshal(payloadLogin)
		if err != nil {
			return nil, util.HandleAppError(err, "Botpress Login", "JSON Marshal", true)
		}

		respon, errCode, errResponse := util.HttpPost(account.AuthURL, jsonData, map[string]string{
			"Content-Type": "application/json",
		})

		if errResponse != nil {
			util.HandleAppError(errResponse, "HttpPost", "Login", false)
			return nil, errResponse
		}
		if errCode >= 400 {
			util.HandleAppError(errResponse, "HttpPost", "Login", false)
			return nil, fmt.Errorf("error code: %d, with response: %+v", errCode, respon)
		}

		var loginResponse LoginBotPressDTO
		errDecode := json.Unmarshal([]byte(respon), &loginResponse)

		if errDecode != nil {
			return nil, util.HandleAppError(err, "Botpress Login", "JSON Unmarshal", true)
		}
		token := loginResponse.Payload.Jwt
		account.Token = token
		if err := a.accountService.SaveAccount(account); err != nil {
			return nil, err
		}
		return &LoginRespon{
			Token:   token,
			BaseURL: account.BaseURL,
		}, nil
	} else {
		return &LoginRespon{
			Token:   account.Token,
			BaseURL: account.BaseURL,
		}, nil
	}
}

func (a *BotpressService) AskBotpress(uniqueId string, token string, baseURL string, botP *AskPayloadBotpresDTO, RefreshToken *RefreshToken) (*BotpressRespon, error) {
	url := fmt.Sprintf("%s/converse/%s/secured?include=state,suggestions,decision,nlu", baseURL, uniqueId)

	jsonData, err := json.Marshal(botP)
	if err != nil {
		return nil, util.HandleAppError(err, "AskBotpress", "JSON Marshal", true)
	}

	body, statusCode, errRespon := util.HttpPost(url, []byte(jsonData),
		map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", token),
			"Content-Type":  "application/json",
		})

	if errRespon != nil {
		return nil, errRespon
	}

	// fmt.Println("Response statusCode:", statusCode)
	// fmt.Println("Response Body:", body)
	pterm.Info.Printfln("Received a message %s: %s", strconv.Itoa(statusCode), body)

	if statusCode >= 400 {
		// refresh token
		if statusCode == 401 {
			refreshToken := true
			a.Login(RefreshToken.BotAccount, RefreshToken.TenantId, &refreshToken)
		}
		return nil, fmt.Errorf("error code: %d, with response: %s", statusCode, body)
	}

	responBotpress, err := UnmarshalBotpressRespon([]byte(body))
	if err != nil {
		return nil, util.HandleAppError(err, "AskBotpress", "JSON Unmarshal", true)
	}

	return &responBotpress, nil
}

// func (a *BotpressService) BPLCOC(payload webhookLivechat.IncomingDTO) (*AskPayloadBotpresDTO, error) {
// 	var botPayload AskPayloadBotpresDTO
// 	switch payload.Action {
// 	case "clientReplyText":
// 		botPayload.Type = BotpressMessageType(TEXT)
// 		botPayload.Text = payload.Message
// 		botPayload.Metadata = payload.Additional
// 	case "clientReplyButton":
// 		botPayload.Type = BotpressMessageType(SINGLE_CHOICE)
// 		botPayload.Text = payload.Message
// 		botPayload.Metadata = payload.Additional
// 	case "clientReplyCarousel":
// 		botPayload.Type = BotpressMessageType(POSTBACK)
// 		botPayload.Payload = payload.Message
// 		botPayload.Metadata = payload.Additional
// 	default:
// 		return nil, fmt.Errorf("unsupported action: %s", payload.Action)
// 	}
// 	return &botPayload, nil
// }

// func (a *BotpressService) BPTLGOF(payload *webhookTelegram.IncomingDTO) (*AskPayloadBotpresDTO, error) {
// 	if payload == nil {
// 		return nil, errors.New("payload is nil")
// 	}

// 	var botPayload AskPayloadBotpresDTO

// 	if payload.Message != nil {
// 		botPayload.Type = TEXT
// 		botPayload.Text = payload.Message.Text
// 		botPayload.Metadata = *payload.Additional
// 		return &botPayload, nil
// 	} else if payload.CallbackQuery != nil {
// 		botPayload.Type = POSTBACK
// 		botPayload.Text = payload.CallbackQuery.Data
// 		botPayload.Payload = payload.CallbackQuery.Data
// 		botPayload.Metadata = *payload.Additional
// 		return &botPayload, nil
// 	}

// 	return nil, errors.New("no valid message or callback query found in payload")
// }

// func (a *BotpressService) BPWASC(payload webhookWhatsapp.IncomingDTO) (*AskPayloadBotpresDTO, error) {
// 	var botPayload AskPayloadBotpresDTO

// 	switch payload.Messages[0].Type {
// 	case "text":
// 		botPayload.Type = BotpressMessageType(TEXT)
// 		botPayload.Text = payload.Messages[0].Text.Body
// 		botPayload.Metadata = payload.Additional
// 		botPayload.Metadata.CustMessage = botPayload.Text

// 	case "interactive":
// 		switch payload.Messages[0].Interactive.Type {
// 		case "list_reply":
// 			botPayload.Type = BotpressMessageType(SINGLE_CHOICE)
// 			botPayload.Text = payload.Messages[0].Interactive.List_reply.Id
// 			botPayload.Metadata = payload.Additional
// 			botPayload.Metadata.CustMessage = botPayload.Text
// 		case "button_reply":
// 			botPayload.Type = BotpressMessageType(SINGLE_CHOICE)
// 			botPayload.Text = payload.Messages[0].Interactive.Button_reply.Id
// 			botPayload.Metadata = payload.Additional
// 			botPayload.Metadata.CustMessage = botPayload.Text
// 		default:
// 			return nil, fmt.Errorf("unsupported interactive type: %s", payload.Messages[0].Interactive.Type)
// 		}

// 	default:
// 		return nil, fmt.Errorf("unsupported action: %s", payload.Messages[0].Type)
// 	}
// 	return &botPayload, nil

// }

// func (a *BotpressService) ParsingPayloadTelegram(payload webhookTelegram.IncomingDTO) (*AskPayloadBotpresDTO, error) {
// 	var botPayload AskPayloadBotpresDTO

// 	if payload.CallbackQuery != nil {
// 		botPayload.Type = BotpressMessageType(POSTBACK)
// 		botPayload.Text = payload.CallbackQuery.Message.Text
// 		botPayload.Metadata = *payload.Additional
// 		botPayload.Payload = payload.CallbackQuery.Data
// 	} else {
// 		botPayload.Type = BotpressMessageType(TEXT)
// 		botPayload.Text = payload.Message.Text
// 		botPayload.Metadata = *payload.Additional
// 	}

// 	return &botPayload, nil
// }
