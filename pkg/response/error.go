package response

import (
	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context, status int, err error, message string) {
	c.JSON(status, Response{Error: map[string]interface{}{
		"raw":     err.Error(),
		"message": message,
	}})
}
