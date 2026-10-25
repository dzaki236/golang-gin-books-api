package helpers

import (
	"books-api/app/structs"

	"github.com/gin-gonic/gin"
)

func Respond(context *gin.Context, status int, message string, data any) {
	context.JSON(status, structs.APIResponse{Message: message, Data: data})
}
