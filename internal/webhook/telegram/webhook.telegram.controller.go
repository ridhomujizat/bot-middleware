package webHookTelegram

import (
	"bot-middleware/internal/application"
	"bot-middleware/internal/pkg/messaging"
	"bot-middleware/internal/pkg/util"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TelegramController struct {
	service     *TelegramService
	application *application.Services
}

func NewTelegramController(messagingGeneral messaging.MessagingGeneral, serviceApplication *application.Services) *TelegramController {
	return &TelegramController{service: NewTelegramService(messagingGeneral), application: serviceApplication}
}

func (t *TelegramController) IncomingHandler(ctx *gin.Context) {
	accountId := ctx.Param("account")
	fmt.Println(accountId)
	account, err := t.application.AccountService.GetUserByName(accountId)
	if err != nil {
		util.APIResponse(ctx, err.Error(), http.StatusInternalServerError, http.MethodPost, nil)
		return
	}

	var input Telegrampayload
	if err := ctx.ShouldBindJSON(&input); err != nil {
		util.APIResponse(ctx, err.Error(), http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	payload := PayloadQueue{}

	payload.AccountId = account.Account
	payload.TenantId = account.TenantId

	if input.CallbackQuery != nil {
		payload.Unique_id = input.CallbackQuery.From.ID
	} else {
		payload.Unique_id = input.Message.Chat.ID
	}

	util.APIResponse(ctx, "Message sent", http.StatusOK, http.MethodPost, payload)
}
