package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID generates a new UUID v4 for every incoming request,
// stores it in the Gin context, and sets the X-Request-ID response header.
// Any client-supplied X-Request-ID header or requestId in the body is ignored.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := uuid.New().String()
		c.Set("requestId", id)
		c.Header("X-Request-ID", id)
		c.Next()
	}
}
