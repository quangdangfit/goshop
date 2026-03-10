package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"goshop/pkg/config"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := config.GetConfig()
		origin := c.Request.Header.Get("Origin")
		allowedOrigins := cfg.CORSAllowedOrigins
		if allowedOrigins == "*" || strings.Contains(allowedOrigins, origin) {
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			c.Header("Access-Control-Allow-Origin", allowedOrigins)
		}
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin,Content-Type,Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
