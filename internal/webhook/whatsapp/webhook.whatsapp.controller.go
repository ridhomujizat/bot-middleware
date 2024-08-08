package webhookWhatsapp

import (
	"bot-middleware/internal/entities"
	"bot-middleware/internal/pkg/messaging"
	"bot-middleware/internal/pkg/util"
	"bot-middleware/internal/webhook"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WhatsappController struct {
	service *WhatsappService
}

func NewWhatsappController(messagingGeneral messaging.MessagingGeneral) *WhatsappController {
	return &WhatsappController{service: NewWhatsappService(messagingGeneral)}
}

// Incoming godoc
// @Summary Incoming WhatsApp
// @Description Incoming message for channel WhatsApp
// @Tags whatsapp
// @Produce json
// @Param botplatform path string true "Bot Platform"
// @Param omnichannel path string true "Omni Channel"
// @Param tenantId path string true "Tenant"
// @Param account path string true "Account"
// @Param payload body IncomingDTO{} true "Payload"
// @Success 200 {object} util.Responses{data=interface{}}
// @Failure 500 {object} util.Responses{data=interface{}}
// @Router /socioconnect/whatsapp/{botplatform}/{omnichannel}/{tenantId}/{account} [post]
func (w *WhatsappController) Incoming(ctx *gin.Context) {
	botplatform := ctx.Param("botplatform")
	omnichannel := ctx.Param("omnichannel")
	tenantId := ctx.Param("tenantId")
	account := ctx.Param("account")

	var payload IncomingDTO
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		util.APIResponse(ctx, err.Error(), http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	if errStr := util.ValidatorErrorResponse(payload); errStr != "" {
		util.APIResponse(ctx, errStr, http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	params := webhook.ParamsDTO{
		BotPlatform: entities.BotPlatform(botplatform),
		Omnichannel: entities.Omnichannel(omnichannel),
		TenantId:    tenantId,
		Account:     account,
	}

	if errStr := util.ValidatorErrorResponse(params); errStr != "" {
		util.APIResponse(ctx, errStr, http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	res, err := w.service.Incoming(params, payload)

	if err != nil {
		util.APIResponse(ctx, err.Error(), http.StatusInternalServerError, http.MethodPost, nil)
		return
	}

	util.APIResponse(ctx, "Incoming received", http.StatusOK, http.MethodPost, res)
}

// Handover godoc
// @Summary Handover WhatsApp
// @Description Handover message for channel WhatsApp
// @Tags whatsapp
// @Produce json
// @Param botplatform path string true "Bot Platform"
// @Param omnichannel path string true "Omni Channel"
// @Param tenantId path string true "Tenant"
// @Param account path string true "Account"
// @Param payload body webhook.HandoverDTO{} true "Payload"
// @Success 200 {object} util.Responses{data=interface{}}
// @Failure 500 {object} util.Responses{data=interface{}}
// @Router /socioconnect/whatsapp/{botplatform}/{omnichannel}/{tenantId}/{account}/handover [post]
func (w *WhatsappController) Handover(ctx *gin.Context) {
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
		BotPlatform: entities.BotPlatform(botplatform),
		Omnichannel: entities.Omnichannel(omnichannel),
		TenantId:    tenantId,
		Account:     account,
	}

	if errStr := util.ValidatorErrorResponse(params); errStr != "" {
		util.APIResponse(ctx, errStr, http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	res, err := w.service.Handover(params, payload)

	if err != nil {
		util.APIResponse(ctx, err.Error(), http.StatusInternalServerError, http.MethodPost, nil)
		return
	}

	util.APIResponse(ctx, "Handover received", http.StatusOK, http.MethodPost, res)
}

// End godoc
// @Summary End WhatsApp
// @Description End message for channel WhatsApp
// @Tags whatsapp
// @Produce json
// @Param botplatform path string true "Bot Platform"
// @Param omnichannel path string true "Omni Channel"
// @Param tenantId path string true "Tenant"
// @Param account path string true "Account"
// @Param payload body EndDTO{} true "Payload"
// @Success 200 {object} util.Responses{data=interface{}}
// @Failure 500 {object} util.Responses{data=interface{}}
// @Router /socioconnect/whatsapp/{botplatform}/{omnichannel}/{tenantId}/{account}/end [post]
func (w *WhatsappController) End(ctx *gin.Context) {
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
		BotPlatform: entities.BotPlatform(botplatform),
		Omnichannel: entities.Omnichannel(omnichannel),
		TenantId:    tenantId,
		Account:     account,
	}

	if errStr := util.ValidatorErrorResponse(params); errStr != "" {
		util.APIResponse(ctx, errStr, http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	res, err := w.service.End(params, payload)

	if err != nil {
		util.APIResponse(ctx, err.Error(), http.StatusInternalServerError, http.MethodPost, nil)
		return
	}

	util.APIResponse(ctx, "End received", http.StatusOK, http.MethodPost, res)
}

// Commerce godoc
// @Summary Commerce WhatsApp
// @Description Commerce message for channel WhatsApp
// @Tags whatsapp
// @Produce json
// @Param botplatform path string true "Bot Platform"
// @Param omnichannel path string true "Omni Channel"
// @Param tenantId path string true "Tenant"
// @Param account path string true "Account"
// @Param payload body IncomingDTO{} true "Payload"
// @Success 200 {object} util.Responses{data=interface{}}
// @Failure 500 {object} util.Responses{data=interface{}}
// @Router /socioconnect/whatsapp/{botplatform}/{omnichannel}/{tenantId}/{account}/commerce [post]
func (w *WhatsappController) Commerce(ctx *gin.Context) {
	botplatform := ctx.Param("botplatform")
	omnichannel := ctx.Param("omnichannel")
	tenantId := ctx.Param("tenantId")
	account := ctx.Param("account")

	var payload IncomingDTO
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		util.APIResponse(ctx, err.Error(), http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	if errStr := util.ValidatorErrorResponse(payload); errStr != "" {
		util.APIResponse(ctx, errStr, http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	params := webhook.ParamsDTO{
		BotPlatform: entities.BotPlatform(botplatform),
		Omnichannel: entities.Omnichannel(omnichannel),
		TenantId:    tenantId,
		Account:     account,
	}

	if errStr := util.ValidatorErrorResponse(params); errStr != "" {
		util.APIResponse(ctx, errStr, http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	res, err := w.service.Commerce(params, payload)

	if err != nil {
		util.APIResponse(ctx, err.Error(), http.StatusInternalServerError, http.MethodPost, nil)
		return
	}

	util.APIResponse(ctx, "Commerce received", http.StatusOK, http.MethodPost, res)
}

// Midtrans godoc
// @Summary Midtrans WhatsApp
// @Description Midtrans message for channel WhatsApp
// @Tags whatsapp
// @Produce json
// @Param botplatform path string true "Bot Platform"
// @Param omnichannel path string true "Omni Channel"
// @Param tenantId path string true "Tenant"
// @Param account path string true "Account"
// @Param payload body interface{} true "Payload"
// @Success 200 {object} util.Responses{data=interface{}}
// @Failure 500 {object} util.Responses{data=interface{}}
// @Router /socioconnect/whatsapp/{botplatform}/{omnichannel}/{tenantId}/{account}/midtrans [post]
func (w *WhatsappController) Midtrans(ctx *gin.Context) {
	botplatform := ctx.Param("botplatform")
	omnichannel := ctx.Param("omnichannel")
	tenantId := ctx.Param("tenantId")
	account := ctx.Param("account")

	var payload interface{}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		util.APIResponse(ctx, err.Error(), http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	params := webhook.ParamsDTO{
		BotPlatform: entities.BotPlatform(botplatform),
		Omnichannel: entities.Omnichannel(omnichannel),
		TenantId:    tenantId,
		Account:     account,
	}

	if errStr := util.ValidatorErrorResponse(params); errStr != "" {
		util.APIResponse(ctx, errStr, http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	res, err := w.service.Midtrans(params, payload)

	if err != nil {
		util.APIResponse(ctx, err.Error(), http.StatusInternalServerError, http.MethodPost, nil)
		return
	}

	util.APIResponse(ctx, "Midtrans received", http.StatusOK, http.MethodPost, res)
}
