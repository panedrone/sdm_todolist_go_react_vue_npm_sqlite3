package resp

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error   string            `json:"error"    binding:"required"`
	Details map[string]string `json:"details"  binding:"omitempty"`
}

func Abort400BadUri(ctx *gin.Context, err error) {
	Abort400hBadRequest(ctx, fmt.Sprintf("Invalid URI: %s", err.Error()))
}

func Abort400hBadRequest(ctx *gin.Context, message string) {
	AbortWithErrResp(ctx, http.StatusBadRequest, message)
}

func Abort400BadJson(ctx *gin.Context, err error) {
	Abort400hBadRequest(ctx, fmt.Sprintf("Invalid JSON: %s", err.Error()))
}

func Abort404NotFound(ctx *gin.Context, err error) {
	AbortWithErrResp(ctx, http.StatusNotFound, err.Error())
}

func Abort500(ctx *gin.Context, err error) {
	AbortWithErrResp(ctx, http.StatusInternalServerError, err.Error())
}

func AbortWithErrResp(ctx *gin.Context, httpStatusCode int, message string) {
	err := ErrorResponse{
		Error: message,
	}
	AbortWithJSON(ctx, httpStatusCode, err)
}

func AbortWithError(ctx *gin.Context, err *ErrorResponse) {
	AbortWithJSON(ctx, http.StatusBadRequest, err)
}

func AbortWithJSON(ctx *gin.Context, httpStatusCode int, jsonObject interface{}) {
	ctx.Header("Content-Type", "application/json")
	ctx.AbortWithStatusJSON(httpStatusCode, jsonObject)
}

func JSON(ctx *gin.Context, httpStatusCode int, jsonObject interface{}) {
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(httpStatusCode, jsonObject)
}
