package webhookTelegram

import (
	"bot-middleware/internal/pkg/messaging"
	"bot-middleware/internal/pkg/util"
	"bot-middleware/internal/webhook"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TelegramController struct {
	service *TelegramService
}

func NewTelegramController(messagingGeneral messaging.MessagingGeneral) *TelegramController {

	return &TelegramController{service: NewTelegramService(messagingGeneral)}
}

// Incoming godoc
// @Summary Incoming telegram
// @Description Incoming message for channel telegram
// @Tags telegram
// @Produce json
// @Param botplatform path string true "Bot Platform"
// @Param omnichannel path string true "Omni Channel"
// @Param tenantId path string true "Tenant"
// @Param account path string true "Account"
// @Param payload body map[string]interface{} true "Payload"
// @Success 200 {object} util.Responses{data=interface{}}
// @Failure 500 {object} util.Responses{data=interface{}}
// @Router /telegram/{botplatform}/{omnichannel}/{tenantId}/{account} [post]
func (t *TelegramController) Incoming(ctx *gin.Context) {
	botplatform := ctx.Param("botplatform")
	omnichannel := ctx.Param("omnichannel")
	tenantId := ctx.Param("tenantId")
	account := ctx.Param("account")

	var payload map[string]interface{}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		util.APIResponse(ctx, err.Error(), http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	if errStr := util.ValidatorErrorResponse(payload); errStr != "" {
		util.APIResponse(ctx, errStr, http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	params := webhook.ParamsDTO{
		BotPlatform: webhook.BotPlatform(botplatform),
		Omnichannel: webhook.Omnichannel(omnichannel),
		TenantId:    tenantId,
		Account:     account,
	}

	if errStr := util.ValidatorErrorResponse(params); errStr != "" {
		util.APIResponse(ctx, errStr, http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	res, err := t.service.Incoming(params, payload)

	if err != nil {
		util.APIResponse(ctx, err.Error(), http.StatusInternalServerError, http.MethodPost, nil)
		return
	}

	util.APIResponse(ctx, "Incoming received", http.StatusOK, http.MethodPost, res)
}

// Handover godoc
// @Summary Handover telegram
// @Description Handover message for channel telegram
// @Tags telegram
// @Produce json
// @Param botplatform path string true "Bot Platform"
// @Param omnichannel path string true "Omni Channel"
// @Param tenantId path string true "Tenant"
// @Param account path string true "Account"
// @Param payload body webhook.HandoverDTO{} true "Payload"
// @Success 200 {object} util.Responses{data=interface{}}
// @Failure 500 {object} util.Responses{data=interface{}}
// @Router /telegram/{botplatform}/{omnichannel}/{tenantId}/{account}/handover [post]
func (t *TelegramController) Handover(ctx *gin.Context) {
	botplatform := ctx.Param("botplatform")
	omnichannel := ctx.Param("omnichannel")
	tenantId := ctx.Param("tenantId")
	account := ctx.Param("account")

	var payload webhook.HandoverDTO
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		util.APIResponse(ctx, err.Error(), http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	if errStr := util.ValidatorErrorResponse(payload); errStr != "" {
		util.APIResponse(ctx, errStr, http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	params := webhook.ParamsDTO{
		BotPlatform: webhook.BotPlatform(botplatform),
		Omnichannel: webhook.Omnichannel(omnichannel),
		TenantId:    tenantId,
		Account:     account,
	}

	if errStr := util.ValidatorErrorResponse(params); errStr != "" {
		util.APIResponse(ctx, errStr, http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	res, err := t.service.Handover(params, payload)

	if err != nil {
		util.APIResponse(ctx, err.Error(), http.StatusInternalServerError, http.MethodPost, nil)
		return
	}

	util.APIResponse(ctx, "Handover received", http.StatusOK, http.MethodPost, res)
}

// End godoc
// @Summary End telegram
// @Description End message for channel telegram
// @Tags telegram
// @Produce json
// @Param botplatform path string true "Bot Platform"
// @Param omnichannel path string true "Omni Channel"
// @Param tenantId path string true "Tenant"
// @Param account path string true "Account"
// @Param payload body EndDTO{} true "Payload"
// @Success 200 {object} util.Responses{data=interface{}}
// @Failure 500 {object} util.Responses{data=interface{}}
// @Router /telegram/{botplatform}/{omnichannel}/{tenantId}/{account}/end [post]
func (t *TelegramController) End(ctx *gin.Context) {
	botplatform := ctx.Param("botplatform")
	omnichannel := ctx.Param("omnichannel")
	tenantId := ctx.Param("tenantId")
	account := ctx.Param("account")

	var payload EndDTO
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		util.APIResponse(ctx, err.Error(), http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	if errStr := util.ValidatorErrorResponse(payload); errStr != "" {
		util.APIResponse(ctx, errStr, http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	params := webhook.ParamsDTO{
		BotPlatform: webhook.BotPlatform(botplatform),
		Omnichannel: webhook.Omnichannel(omnichannel),
		TenantId:    tenantId,
		Account:     account,
	}

	if errStr := util.ValidatorErrorResponse(params); errStr != "" {
		util.APIResponse(ctx, errStr, http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	res, err := t.service.End(params, payload)

	if err != nil {
		util.APIResponse(ctx, err.Error(), http.StatusInternalServerError, http.MethodPost, nil)
		return
	}

	util.APIResponse(ctx, "End received", http.StatusOK, http.MethodPost, res)
}
