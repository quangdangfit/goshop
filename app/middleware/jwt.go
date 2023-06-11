package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"goshop/pkg/jtoken"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, nil)
			c.Abort()
			return
		}

		payload, err := jtoken.ValidateToken(token)
		if err != nil || payload == nil || payload["type"] != jtoken.AccessTokenType {
			c.JSON(http.StatusUnauthorized, nil)
			c.Abort()
			return
		}
		c.Set("userId", payload["id"])
		c.Set("role", payload["role"])
		c.Next()
	}
}

func JWTRefresh() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, nil)
			c.Abort()
			return
		}

		payload, err := jtoken.ValidateToken(token)
		if err != nil || payload == nil || payload["type"] != jtoken.RefreshTokenType {
			c.JSON(http.StatusUnauthorized, nil)
			c.Abort()
			return
		}
		c.Set("userId", payload["id"])
		c.Set("role", payload["role"])
		c.Next()
	}
}
