package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"goshop/pkg/config"
	"goshop/pkg/redis"
	"goshop/pkg/response"
)

func RateLimit(cache redis.Redis) gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := config.GetConfig()
		maxRequests := cfg.RateLimitRequests
		window := time.Duration(cfg.RateLimitWindowSeconds) * time.Second

		if maxRequests <= 0 {
			c.Next()
			return
		}

		key := fmt.Sprintf("rate_limit:%s", c.ClientIP())
		count, err := cache.Incr(key, window)
		if err != nil {
			c.Next()
			return
		}

		if int(count) > maxRequests {
			response.Error(c, http.StatusTooManyRequests, fmt.Errorf("rate limit exceeded"), "Too many requests")
			c.Abort()
			return
		}

		c.Next()
	}
}
