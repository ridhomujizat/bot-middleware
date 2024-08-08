package webhookLivechat

import (
	"bot-middleware/internal/entities"
	"bot-middleware/internal/pkg/messaging"
	"bot-middleware/internal/pkg/util"
	"bot-middleware/internal/webhook"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LivechatController struct {
	service *LivechatService
}

func NewLivechatController(messagingGeneral messaging.MessagingGeneral) *LivechatController {

	return &LivechatController{service: NewLivechatService(messagingGeneral)}
}

// Incoming godoc
// @Summary Incoming livechat
// @Description Incoming message for channel livechat
// @Tags livechat
// @Produce json
// @Param botplatform path string true "Bot Platform"
// @Param omnichannel path string true "Omni Channel"
// @Param tenantId path string true "Tenant"
// @Param account path string true "Account"
// @Param payload body IncomingDTO{} true "Payload"
// @Success 200 {object} util.Responses{data=interface{}}
// @Failure 500 {object} util.Responses{data=interface{}}
// @Router /octopushchat/livechat/{botplatform}/{omnichannel}/{tenantId}/{account} [post]
func (l *LivechatController) Incoming(ctx *gin.Context) {
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

	res, err := l.service.Incoming(params, payload)

	if err != nil {
		util.APIResponse(ctx, err.Error(), http.StatusInternalServerError, http.MethodPost, nil)
		return
	}

	util.APIResponse(ctx, "Incoming received", http.StatusOK, http.MethodPost, res)
}

// Handover godoc
// @Summary Handover livechat
// @Description Handover message for channel livechat
// @Tags livechat
// @Produce json
// @Param botplatform path string true "Bot Platform"
// @Param omnichannel path string true "Omni Channel"
// @Param tenantId path string true "Tenant"
// @Param account path string true "Account"
// @Param payload body webhook.HandoverDTO{} true "Payload"
// @Success 200 {object} util.Responses{data=interface{}}
// @Failure 500 {object} util.Responses{data=interface{}}
// @Router /octopushchat/livechat/{botplatform}/{omnichannel}/{tenantId}/{account}/handover [post]
func (l *LivechatController) Handover(ctx *gin.Context) {
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

	res, err := l.service.Handover(params, payload)

	if err != nil {
		util.APIResponse(ctx, err.Error(), http.StatusInternalServerError, http.MethodPost, nil)
		return
	}

	util.APIResponse(ctx, "Handover received", http.StatusOK, http.MethodPost, res)
}

// End godoc
// @Summary End livechat
// @Description End message for channel livechat
// @Tags livechat
// @Produce json
// @Param botplatform path string true "Bot Platform"
// @Param omnichannel path string true "Omni Channel"
// @Param tenantId path string true "Tenant"
// @Param account path string true "Account"
// @Param payload body EndDTO{} true "Payload"
// @Success 200 {object} util.Responses{data=interface{}}
// @Failure 500 {object} util.Responses{data=interface{}}
// @Router /octopushchat/livechat/{botplatform}/{omnichannel}/{tenantId}/{account}/end [post]
func (l *LivechatController) End(ctx *gin.Context) {
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

	res, err := l.service.End(params, payload)

	if err != nil {
		util.APIResponse(ctx, err.Error(), http.StatusInternalServerError, http.MethodPost, nil)
		return
	}

	util.APIResponse(ctx, "End received", http.StatusOK, http.MethodPost, res)
}
