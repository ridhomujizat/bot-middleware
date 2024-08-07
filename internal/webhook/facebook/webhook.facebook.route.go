package webhookFacebook

import (
	"bot-middleware/internal/pkg/messaging"

	"github.com/gin-gonic/gin"
)

func InitRouterFacebook(messagingGeneral messaging.MessagingGeneral, router *gin.RouterGroup) {
	facebookController := NewFacebookController(messagingGeneral)

	// Group Endpoint for facebook
	routeGroupFacebook := router.Group("/facebook")
	routeGroupFacebook.POST("/:botplatform/:omnichannel/:tenantId/:account", facebookController.Incoming)
	routeGroupFacebook.POST("/:botplatform/:omnichannel/:tenantId/:account/handover", facebookController.Handover)
	routeGroupFacebook.POST("/:botplatform/:omnichannel/:tenantId/:account/end", facebookController.End)
}
