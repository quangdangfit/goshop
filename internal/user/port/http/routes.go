package http

import (
	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/validation"

	"goshop/internal/user/repository"
	"goshop/internal/user/service"
	"goshop/pkg/dbs"
	"goshop/pkg/middleware"
)

func Routes(r *gin.RouterGroup, sqlDB dbs.Database, validator validation.Validation) {
	userRepo := repository.NewUserRepository(sqlDB)
	userSvc := service.NewUserService(validator, userRepo)
	userHandler := NewUserHandler(userSvc)

	addressRepo := repository.NewAddressRepository(sqlDB)
	addressSvc := service.NewAddressService(validator, addressRepo)
	addressHandler := NewAddressHandler(addressSvc)

	wishlistRepo := repository.NewWishlistRepository(sqlDB)
	wishlistSvc := service.NewWishlistService(wishlistRepo)
	wishlistHandler := NewWishlistHandler(wishlistSvc)

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

	addressRoute := r.Group("/addresses", authMiddleware)
	{
		addressRoute.GET("", addressHandler.ListAddresses)
		addressRoute.POST("", addressHandler.CreateAddress)
		addressRoute.GET("/:id", addressHandler.GetAddressByID)
		addressRoute.PUT("/:id", addressHandler.UpdateAddress)
		addressRoute.DELETE("/:id", addressHandler.DeleteAddress)
		addressRoute.PUT("/:id/default", addressHandler.SetDefaultAddress)
	}

	wishlistRoute := r.Group("/wishlist", authMiddleware)
	{
		wishlistRoute.GET("", wishlistHandler.GetWishlist)
		wishlistRoute.POST("", wishlistHandler.AddProduct)
		wishlistRoute.DELETE("/:productId", wishlistHandler.RemoveProduct)
	}
}
