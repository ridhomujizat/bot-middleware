package webhookTole

import (
	"bot-middleware/internal/pkg/messaging"
	"bot-middleware/internal/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ToleController struct {
	service *ToleService
}

func NewToleController(messagingGeneral messaging.MessagingGeneral) *ToleController {
	return &ToleController{service: NewToleService(messagingGeneral)}
}

// SendMessageHandler godoc
// @Summary Send a message to tole queue
// @Description Send a message to rabbit tole queue using tenant channel
// @Tags tole
// @Produce json
// @Param queueName path string true "Queue Name"
// @Param payload body interface{} true "Payload"
// @Success 200 {object} util.Responses{data=interface{}}
// @Failure 500 {object} util.Responses{data=interface{}}
// @Router /tole/{queueName} [post]
func (t *ToleController) SendMessageHandler(ctx *gin.Context) {
	queueName := ctx.Param("queueName")
	var message interface{}
	if err := ctx.ShouldBindJSON(&message); err != nil {
		util.APIResponse(ctx, err.Error(), http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	if err := t.service.Send(queueName, message); err != nil {
		util.APIResponse(ctx, err.Error(), http.StatusInternalServerError, http.MethodPost, nil)
		return
	}

	util.APIResponse(ctx, "Message sent", http.StatusOK, http.MethodPost, nil)
}
