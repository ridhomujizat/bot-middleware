package webhookLivechat

import (
	"bot-middleware/internal/pkg/messaging"

	"github.com/gin-gonic/gin"
)

func InitRouterLivechat(messagingGeneral messaging.MessagingGeneral, router *gin.RouterGroup) {
	livechatController := NewLivechatController(messagingGeneral)

	// Group Endpoint for livechat
	routeGroupLivechat := router.Group("/octopushchat/livechat")
	routeGroupLivechat.POST("/:botplatform/:omnichannel/:tenantId/:account", livechatController.Incoming)
	routeGroupLivechat.POST("/:botplatform/:omnichannel/:tenantId/:account/handover", livechatController.Handover)
	routeGroupLivechat.POST("/:botplatform/:omnichannel/:tenantId/:account/end", livechatController.End)
}
