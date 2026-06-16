package middleware

import (
	"context"
	"net/http"
	"time"

	"viscraft-backend/constant"
	"viscraft-backend/model/response"

	"github.com/gin-gonic/gin"
)

// Timeout wraps the request context with a deadline. If the handler does not
// complete within the given duration, the middleware responds with HTTP 504
// and a standard error body. Background goroutines spawned by handlers are
// unaffected because they create their own context (context.Background()).
func Timeout(duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a child context with the specified timeout.
		ctx, cancel := context.WithTimeout(c.Request.Context(), duration)
		defer cancel()

		// Replace the request context so downstream handlers observe the deadline.
		c.Request = c.Request.WithContext(ctx)

		// Channel signals that the handler finished before the deadline.
		done := make(chan struct{}, 1)

		go func() {
			c.Next()
			done <- struct{}{}
		}()

		select {
		case <-done:
			// Handler completed in time — response already written.
			return
		case <-ctx.Done():
			// Deadline exceeded — abort and return 504.
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
