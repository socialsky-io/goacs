package response

import (
	"github.com/gin-gonic/gin"
	"goacs/http/request"
	"goacs/repository"
)

func ResponseData(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, responseMap("Ok", data))
}

func ResponsePaginatior(ctx *gin.Context, response repository.PaginatorResponse) {
	ctx.JSON(200, response)
}

func ResponseError(ctx *gin.Context, code int, message string, data interface{}) {
	ctx.JSON(code, responseMap("Error", data))
}

func Response500(ctx *gin.Context, message string, data interface{}) {
	ResponseError(ctx, 500, message, responseMap("Server error", data))
}

func ResponseValidationErrors(ctx *gin.Context, validator *request.ApiValidator) {
	ctx.JSON(422, responseMap("Validation error", validator.Errors))
}

func responseMap(message string, data interface{}) map[string]interface{} {
	return map[string]interface{}{"message": message, "data": data}
}
