package util

import "github.com/gin-gonic/gin"

type Responses struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Method  string      `json:"method"`
	Data    interface{} `json:"data"`
}

func ValidatorErrorResponse(ctx *gin.Context, Status int, Message string, Method string, Data interface{}) {
	errResponse := Responses{
		Status:  Status,
		Method:  Method,
		Message: Message,
		Data:    Data,
	}

	ctx.JSON(Status, errResponse)
	defer ctx.AbortWithStatus(Status)
}

func APIResponse(ctx *gin.Context, Message string, Status int, Method string, Data interface{}) {
	jsonResponse := Responses{
		Status:  Status,
		Method:  Method,
		Message: Message,
		Data:    Data,
	}

	if Status >= 400 {
		ctx.IndentedJSON(Status, jsonResponse)
		defer ctx.AbortWithStatus(Status)
	} else {
		ctx.IndentedJSON(Status, jsonResponse)
	}
}
