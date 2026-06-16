package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"viscraft-backend/constant"
	"viscraft-backend/model/response"
	"viscraft-backend/pkg/logger"
)

// Recovery returns a Gin middleware that catches panics during request handling,
// logs the panic value and stack trace with the requestId, and returns a generic
// HTTP 500 error response without exposing internal details.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				requestId := c.GetString("requestId")
				stack := string(debug.Stack())

				logger.Error(requestId, "panic recovered",
					"panic", fmt.Sprintf("%v", r),
					"stack", stack,
				)

				c.AbortWithStatusJSON(http.StatusInternalServerError, response.BaseResponse{
					RequestId: requestId,
					Success:   false,
					ErrorCode: constant.ErrInternalServer.Code,
					Message:   constant.ErrInternalServer.Message,
				})
			}
		}()

		c.Next()
	}
}
