package http

import (
	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/validation"

	"goshop/internal/user/repository"
	"goshop/internal/user/service"
	"goshop/pkg/config"
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

	cfg := config.GetConfig()
	authRoute := r.Group("/auth")

	var authMiddleware gin.HandlerFunc
	var refreshAuthMiddleware gin.HandlerFunc

	// Protected routes always use JWTAuth — both auth modes mint GoShop JWT
	// session tokens (the OIDC callback exchanges the Authentik ID token for
	// our own JWT pair before redirecting back to the FE).
	authMiddleware = middleware.JWTAuth()
	refreshAuthMiddleware = middleware.JWTRefresh()

	// /auth/refresh, /auth/me and /auth/change-password are mounted in both
	// modes so the FE can refresh tokens and load the profile regardless of
	// how the user signed in.
	authRoute.POST("/refresh", refreshAuthMiddleware, userHandler.RefreshToken)
	authRoute.GET("/me", authMiddleware, userHandler.GetMe)

	if cfg.AuthMode == config.AuthModeOIDC {
		// OIDC mode: Authentik owns credentials, so password-based endpoints
		// are disabled. /auth/login redirects to Authentik; /auth/callback
		// trades the code for tokens and bounces the browser back to the FE.
		oidcHandler := NewOIDCHandler(userSvc, cfg)
		authRoute.GET("/login", oidcHandler.Login)
		authRoute.GET("/callback", oidcHandler.Callback)
	} else {
		// Legacy JWT mode: email + bcrypt password.
		authRoute.POST("/register", userHandler.Register)
		authRoute.POST("/login", userHandler.Login)
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

	r.PUT("/me/cart-snapshot", authMiddleware, PutCartSnapshot)
}
