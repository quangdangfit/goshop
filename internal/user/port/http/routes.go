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

	// Password-based endpoints are always mounted — both auth modes accept
	// email/password sign-in and self-service registration. OIDC mode adds
	// SSO on top of that rather than replacing it.
	authRoute.POST("/register", userHandler.Register)
	authRoute.POST("/login", userHandler.Login)
	authRoute.POST("/refresh", refreshAuthMiddleware, userHandler.RefreshToken)
	authRoute.GET("/me", authMiddleware, userHandler.GetMe)
	authRoute.PUT("/change-password", authMiddleware, userHandler.ChangePassword)

	if cfg.AuthMode == config.AuthModeOIDC {
		// OIDC mode additionally exposes GET /auth/login (redirects to
		// Authentik with an optional ?provider hint) and /auth/callback
		// (exchanges the code, upserts the user, mints a GoShop JWT pair,
		// then bounces the browser to the FE). The POST /auth/login route
		// above coexists thanks to method-based routing.
		oidcHandler := NewOIDCHandler(userSvc, cfg)
		authRoute.GET("/login", oidcHandler.Login)
		authRoute.GET("/callback", oidcHandler.Callback)
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
