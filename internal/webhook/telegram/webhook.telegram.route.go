package webhookTelegram

import (
	"bot-middleware/internal/application"
	"bot-middleware/internal/pkg/messaging"

	"github.com/gin-gonic/gin"
)

func InitRouterTelegram(messagingGeneral messaging.MessagingGeneral, router *gin.RouterGroup, serviceApplication *application.Services) {
	telegramController := NewTelegramController(messagingGeneral, serviceApplication)

	// Group Endpoint for telegram
	routeGroupTelegram := router.Group("/telegram")
	routeGroupTelegram.POST("/:botplatform/:omnichannel/:tenantId/:account", telegramController.Incoming)
	routeGroupTelegram.POST("/:botplatform/:omnichannel/:tenantId/:account/handover", telegramController.Handover)
	routeGroupTelegram.POST("/:botplatform/:omnichannel/:tenantId/:account/end", telegramController.End)
}
