package webhookTole

import (
	"bot-middleware/internal/pkg/messaging"

	"github.com/gin-gonic/gin"
)

func InitRouterTole(messagingGeneral messaging.MessagingGeneral, router *gin.RouterGroup) {
	toleController := NewToleController(messagingGeneral)

	// Group Endpoint for tole
	routeGroupTole := router.Group("/tole")
	routeGroupTole.POST("/:queueName", toleController.SendMessageHandler)
}
