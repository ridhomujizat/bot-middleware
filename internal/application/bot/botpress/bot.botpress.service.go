package botpress

import (
	appAccount "bot-middleware/internal/application/account"
	"bot-middleware/internal/entities"
	"bot-middleware/internal/pkg/util"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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

func (a *BotpressService) Login() (*LoginRespon, error) {

	account, err := a.accountService.GetUserByAccountId("libra_onx")
	if err != nil {
		return nil, err
	}

	payloadLogin := loginPayload{
		Email:    account.Username,
		Password: account.Password,
	}
	// Marshal the payload into JSON
	jsonData, err := json.Marshal(payloadLogin)

	respon, statusCode, errResponse := util.HttpPost(account.AuthURL, jsonData, map[string]string{
		"Content-Type": "application/json",
	})
	if errResponse != nil {
		util.HandleAppError(errResponse, "HttpPost", "Login", false)
		return nil, errResponse
	}

	if statusCode == http.StatusOK {
		var loginResponse LoginBotPressDTO
		errDecode := json.Unmarshal([]byte(respon), &loginResponse)
		if errDecode != nil {

			fmt.Println("Error decoding JSON:", errDecode)
			return nil, errDecode
		}
		return &LoginRespon{
			Token:   loginResponse.Payload.Jwt,
			BaseURL: account.BaseURL,
		}, nil
	} else {
		fmt.Println("Request failed with status:", statusCode)
		return nil, err
	}
}

func (a *BotpressService) AskBotpress(uniqueId string, token string, baseURL string, botP *AskPayloadBotpresDTO) ([]Response, error) {
	url := fmt.Sprintf("%s/converse/%s/secured?include=state,suggestions,decision,nlu", baseURL, uniqueId)

	jsonData, err := json.Marshal(botP)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	}

	body, statusCode, errRespon := util.HttpPost(url, []byte(jsonData),
		map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", token),
			"Content-Type":  "application/json",
		})

	if errRespon != nil {
		util.HandleAppError(errRespon, "POST", "util.HttpPost", false)
		return nil, errRespon
	}

	// fmt.Println("Response statusCode:", statusCode)
	// fmt.Println("Response Body:", body)
	pterm.Info.Printfln("Received a message %s: %s", strconv.Itoa(statusCode), body)

	responBotpress, err := UnmarshalBotpressRespon([]byte(body))
	if err != nil {
		util.HandleAppError(err, "UnmarshalBotpressRespon", "AskBotpress", false)
		return nil, err
	}

	var Responses = responBotpress.Responses
	return Responses, nil
}

func (a *BotpressService) ParsingPayloadTelegram(payload entities.IncomingTelegramDTO) (*AskPayloadBotpresDTO, error) {
	var botPayload AskPayloadBotpresDTO

	if payload.CallbackQuery != nil {
		botPayload.Type = BotpressMessageType(POSTBACK)
		botPayload.Text = payload.CallbackQuery.Message.Text
		botPayload.Metadata = *payload.Additional
		botPayload.Payload = payload.CallbackQuery.Data
	} else {
		botPayload.Type = BotpressMessageType(TEXT)
		botPayload.Text = payload.Message.Text
		botPayload.Metadata = *payload.Additional
	}

	return &botPayload, nil
}
