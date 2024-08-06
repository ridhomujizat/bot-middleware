package webhookWhatsapp

import (
	"bot-middleware/internal/pkg/messaging"

	"github.com/gin-gonic/gin"
)

func InitRouterWhatsapp(messagingGeneral messaging.MessagingGeneral, router *gin.RouterGroup) {
	whatsappController := NewWhatsappController(messagingGeneral)

	// Group Endpoint for socioconnect whatsapp
	routeGroupWhatsapp := router.Group("/socioconnect/whatsapp")
	routeGroupWhatsapp.POST("/:botplatform/:omnichannel/:tenantId/:account", whatsappController.Incoming)
	routeGroupWhatsapp.POST("/:botplatform/:omnichannel/:tenantId/:account/handover", whatsappController.Handover)
	routeGroupWhatsapp.POST("/:botplatform/:omnichannel/:tenantId/:account/end", whatsappController.End)
	routeGroupWhatsapp.POST("/:botplatform/:omnichannel/:tenantId/:account/commerce", whatsappController.Commerce)
	routeGroupWhatsapp.POST("/:botplatform/:omnichannel/:tenantId/:account/midtrans", whatsappController.Midtrans)
}
