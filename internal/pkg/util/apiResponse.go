package util

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Responses struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Method  string      `json:"method"`
	Data    interface{} `json:"data"`
}

func ValidatorErrorResponse(payload interface{}) string {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	if err := validate.Struct(payload); err != nil {
		var errorMessages []string

		validationErrors := validationError(err)
		if validationErrors != "" {
			errorMessages = append(errorMessages, validationErrors)
		}
		message := "Validation failed: " + strings.Join(errorMessages, ", ")
		return message
	}

	return ""
}

func validationError(err error) string {
	if errs, ok := err.(validator.ValidationErrors); ok {
		var sb strings.Builder
		for _, e := range errs {
			fieldName := e.Field()
			tag := e.Tag()
			param := e.Param()
			msg, ok := validationMessages[tag]
			if !ok {
				msg = tag
			} else {
				msg = fmt.Sprintf(msg, param)
			}
			sb.WriteString(fmt.Sprintf("%s: %s", fieldName, msg))
			sb.WriteString(", ")
		}
		return strings.TrimSuffix(sb.String(), ", ")
	}
	return err.Error()
}

var validationMessages = map[string]string{
	"required":    "is required",
	"url":         "must be a valid URL",
	"datetime":    "must be a valid date-time format (2006-01-02T15:04:05Z07:00)",
	"number":      "must be a number",
	"oneof":       "must be one of the allowed values: %s",
	"email":       "must be a valid email address",
	"min":         "must be greater than or equal to %s",
	"max":         "must be less than or equal to %s",
	"len":         "must have the exact length of %s",
	"alpha":       "must contain only alphabetic characters",
	"alphanum":    "must contain only alphanumeric characters",
	"eqfield":     "must be equal to the value of the %s field",
	"nefield":     "must not be equal to the value of the %s field",
	"gt":          "must be greater than %s",
	"gte":         "must be greater than or equal to %s",
	"lt":          "must be less than %s",
	"lte":         "must be less than or equal to %s",
	"excludes":    "must not contain the value %s",
	"excludesall": "must not contain any of the values: %s",
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
