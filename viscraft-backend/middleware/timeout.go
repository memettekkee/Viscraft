package middleware

import (
	"context"
	"net/http"
	"time"

	"viscraft-backend/constant"
	"viscraft-backend/model/response"

	"github.com/gin-gonic/gin"
)


func Timeout(duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), duration)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		done := make(chan struct{}, 1)

		go func() {
			c.Next()
			done <- struct{}{}
		}()

		select {
		case <-done:
			return
		case <-ctx.Done():
			requestId := c.GetString("requestId")
			c.Abort()
			c.JSON(http.StatusGatewayTimeout, response.BaseResponse{
				RequestId: requestId,
				Success:   false,
				ErrorCode: constant.ErrGeminiTimeout.Code,
				Message:   constant.ErrGeminiTimeout.Message,
			})
		}
	}
}
