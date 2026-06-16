package middleware

import (
	"net/http"
	"sync"
	"time"

	"viscraft-backend/constant"
	"viscraft-backend/model/response"

	"github.com/gin-gonic/gin"
)

// RateLimit returns a Gin middleware that enforces a per-user request rate limit.
// It allows at most maxRequests requests within the given window duration per user.
// Requests exceeding the limit receive HTTP 429 with error code ERR_02.
func RateLimit(maxRequests int, window time.Duration) gin.HandlerFunc {
	var (
		limiter = make(map[string][]time.Time)
		mu      sync.Mutex
	)

	return func(c *gin.Context) {
		userId := c.GetString("userId")
		now := time.Now()
		windowStart := now.Add(-1 * window)

		mu.Lock()
		defer mu.Unlock()

		// Prune expired entries
		recent := make([]time.Time, 0)
		for _, t := range limiter[userId] {
			if t.After(windowStart) {
				recent = append(recent, t)
			}
		}

		// Check limit
		if len(recent) >= maxRequests {
			requestId := c.GetString("requestId")
			c.JSON(http.StatusTooManyRequests, response.BaseResponse{
				RequestId: requestId,
				Success:   false,
				ErrorCode: constant.ErrTooManyRequest.Code,
				Message:   constant.ErrTooManyRequest.Message,
			})
			c.Abort()
			return
		}

		// Record this request
		limiter[userId] = append(recent, now)
		c.Next()
	}
}
