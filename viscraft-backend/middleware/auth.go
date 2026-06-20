package middleware

import (
	"net/http"
	"strings"

	"viscraft-backend/constant"
	"viscraft-backend/model/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			abortWithUnauthorized(c)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			abortWithUnauthorized(c)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			abortWithUnauthorized(c)
			return
		}

		userId, ok := claims["userId"].(string)
		if !ok || userId == "" {
			abortWithUnauthorized(c)
			return
		}

		c.Set("userId", userId)
		c.Next()
	}
}

// abortWithUnauthorized sends a 401 response with ERR_09 and aborts the request chain.
func abortWithUnauthorized(c *gin.Context) {
	appErr := constant.ErrUnauthorized
	c.JSON(http.StatusUnauthorized, response.BaseResponse{
		RequestId: c.GetString("requestId"),
		Success:   false,
		ErrorCode: appErr.Code,
		Message:   appErr.Message,
	})
	c.Abort()
}
