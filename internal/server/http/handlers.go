package http

import (
	"github.com/gin-gonic/gin"

	userHttp "goshop/internal/user/port/http"
)

func (s Server) MapRoutes(e *gin.Engine) error {
	v1 := e.Group("/api/v1")
	userHttp.Routes(v1)
	return nil
}
