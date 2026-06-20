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
