package http

import (
	"github.com/gin-gonic/gin"

	"goshop/internal/notification/repository"
	"goshop/internal/notification/service"
	"goshop/pkg/dbs"
	"goshop/pkg/middleware"
)

func Routes(r *gin.RouterGroup, db dbs.Database) {
	repo := repository.NewPreferenceRepository(db)
	svc := service.NewPreferenceService(repo)
	h := NewHandler(svc)

	g := r.Group("/me/notification-preferences", middleware.JWTAuth())
	g.GET("", h.ListPreferences)
	g.PUT("", h.SetPreference)
}
