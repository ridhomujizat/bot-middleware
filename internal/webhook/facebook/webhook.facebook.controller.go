package webhookFacebook

import (
	"bot-middleware/internal/pkg/messaging"
	"bot-middleware/internal/pkg/util"
	"bot-middleware/internal/webhook"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FacebookController struct {
	service *FacebookService
}

func NewFacebookController(messagingGeneral messaging.MessagingGeneral) *FacebookController {

	return &FacebookController{service: NewFacebookService(messagingGeneral)}
}

// Incoming godoc
// @Summary Incoming facebook
// @Description Incoming message for channel facebook
// @Tags facebook
// @Produce json
// @Param botplatform path string true "Bot Platform"
// @Param omnichannel path string true "Omni Channel"
// @Param tenantId path string true "Tenant"
// @Param account path string true "Account"
// @Param payload body IncomingDTO{} true "Payload"
// @Success 200 {object} util.Responses{data=interface{}}
// @Failure 500 {object} util.Responses{data=interface{}}
// @Router /facebook/{botplatform}/{omnichannel}/{tenantId}/{account} [post]
func (f *FacebookController) Incoming(ctx *gin.Context) {
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
		BotPlatform: webhook.BotPlatform(botplatform),
		Omnichannel: webhook.Omnichannel(omnichannel),
		TenantId:    tenantId,
		Account:     account,
	}

	if errStr := util.ValidatorErrorResponse(params); errStr != "" {
		util.APIResponse(ctx, errStr, http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	res, err := f.service.Incoming(params, payload)

	if err != nil {
		util.APIResponse(ctx, err.Error(), http.StatusInternalServerError, http.MethodPost, nil)
		return
	}

	util.APIResponse(ctx, "Incoming received", http.StatusOK, http.MethodPost, res)
}

// Handover godoc
// @Summary Handover facebook
// @Description Handover message for channel facebook
// @Tags facebook
// @Produce json
// @Param botplatform path string true "Bot Platform"
// @Param omnichannel path string true "Omni Channel"
// @Param tenantId path string true "Tenant"
// @Param account path string true "Account"
// @Param payload body webhook.HandoverDTO{} true "Payload"
// @Success 200 {object} util.Responses{data=interface{}}
// @Failure 500 {object} util.Responses{data=interface{}}
// @Router /facebook/{botplatform}/{omnichannel}/{tenantId}/{account}/handover [post]
func (f *FacebookController) Handover(ctx *gin.Context) {
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

	res, err := f.service.Handover(params, payload)

	if err != nil {
		util.APIResponse(ctx, err.Error(), http.StatusInternalServerError, http.MethodPost, nil)
		return
	}

	util.APIResponse(ctx, "Handover received", http.StatusOK, http.MethodPost, res)
}

// End godoc
// @Summary End facebook
// @Description End message for channel facebook
// @Tags facebook
// @Produce json
// @Param botplatform path string true "Bot Platform"
// @Param omnichannel path string true "Omni Channel"
// @Param tenantId path string true "Tenant"
// @Param account path string true "Account"
// @Param payload body EndDTO{} true "Payload"
// @Success 200 {object} util.Responses{data=interface{}}
// @Failure 500 {object} util.Responses{data=interface{}}
// @Router /facebook/{botplatform}/{omnichannel}/{tenantId}/{account}/end [post]
func (f *FacebookController) End(ctx *gin.Context) {
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

	res, err := f.service.End(params, payload)

	if err != nil {
		util.APIResponse(ctx, err.Error(), http.StatusInternalServerError, http.MethodPost, nil)
		return
	}

	util.APIResponse(ctx, "End received", http.StatusOK, http.MethodPost, res)
}
