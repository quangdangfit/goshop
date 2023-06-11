package response

import (
	"github.com/gin-gonic/gin"
)

// Response schema
type Response struct {
	Result interface{} `json:"result"`
	Error  interface{} `json:"error"`
}

func JSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, Response{
		Result: data,
	})
}
