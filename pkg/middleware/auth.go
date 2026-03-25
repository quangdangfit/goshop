package middleware

import (
	"github.com/gin-gonic/gin"

	"goshop/pkg/apperror"
	"goshop/pkg/jtoken"
)

func JWTAuth() gin.HandlerFunc {
	return JWT(jtoken.AccessTokenType)
}

func JWTRefresh() gin.HandlerFunc {
	return JWT(jtoken.RefreshTokenType)
}

func JWT(tokenType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			apperror.ErrUnauthorized.HTTPError(c)
			c.Abort()
			return
		}

		payload, err := jtoken.ValidateToken(token)
		if err != nil || payload == nil || payload["type"] != tokenType {
			apperror.ErrUnauthorized.HTTPError(c)
			c.Abort()
			return
		}
		c.Set("userId", payload["id"])
		c.Set("role", payload["role"])
		c.Next()
	}
}
