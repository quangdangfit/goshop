package http

import (
	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/validation"
	"gorm.io/gorm"

	"goshop/internal/user/repository"
	"goshop/internal/user/service"
	"goshop/pkg/middleware"
)

func Routes(r *gin.RouterGroup, db *gorm.DB, validator validation.Validation) {
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(validator, userRepo)
	userHandler := NewUserHandler(userSvc)

	authMiddleware := middleware.JWTAuth()
	refreshAuthMiddleware := middleware.JWTRefresh()
	authRoute := r.Group("/auth")
	{
		authRoute.POST("/register", userHandler.Register)
		authRoute.POST("/login", userHandler.Login)
		authRoute.POST("/refresh", refreshAuthMiddleware, userHandler.RefreshToken)
		authRoute.GET("/me", authMiddleware, userHandler.GetMe)
		authRoute.PUT("/change-password", authMiddleware, userHandler.ChangePassword)
	}
}
