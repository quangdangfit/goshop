package cache

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/quangdangfit/gocommon/utils/logger"
)

type Cache interface {
	IsConnected() bool
	Get(key string, data interface{}) error
	Set(key string, val []byte) error
	Remove(key string) error
}

// Setup Initialize the Cache instance
func New() Cache {
	return NewRedis()
}

var cache = New()

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func Cached() gin.HandlerFunc {
	return func(c *gin.Context) {
		if cache == nil || !cache.IsConnected() {
			logger.Warn("Cache cache is not available")
			c.Next()
			return
		}

		key := c.Request.URL.RequestURI()
		if c.Request.Method != "GET" {
			c.Next()

			if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "DELETE" {
				cache.Remove(key)
			}
			return
		}

		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		var data map[string]interface{}
		cache.Get(key, &data)

		if data != nil {
			c.JSON(http.StatusOK, data)
			c.Abort()
			return
		}

		c.Next()

		statusCode := w.Status()
		if statusCode == http.StatusOK {
			cache.Set(key, w.body.Bytes())
		}
	}
}
