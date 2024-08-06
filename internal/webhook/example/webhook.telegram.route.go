package webHookTelegram

import (
	"bot-middleware/internal/application"
	"bot-middleware/internal/pkg/messaging"

	"github.com/gin-gonic/gin"
)

func InitRouterTelegram(messagingGeneral messaging.MessagingGeneral, router *gin.RouterGroup, serviceApplication *application.Services) {
	telegramController := NewTelegramController(messagingGeneral, serviceApplication)

	// Group Endpoint for tole
	routeGroupTole := router.Group("/telegram")
	routeGroupTole.POST("/incoming/:account", telegramController.IncomingHandler)
}
