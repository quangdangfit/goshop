package http

import (
	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/validation"

	"goshop/internal/user/repository"
	"goshop/internal/user/service"
	"goshop/pkg/dbs"
	"goshop/pkg/middleware"
)

func Routes(r *gin.RouterGroup, sqlDB dbs.IDatabase, validator validation.Validation) {
	userRepo := repository.NewUserRepository(sqlDB)
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
