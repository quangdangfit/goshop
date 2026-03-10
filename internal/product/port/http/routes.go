package http

import (
	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/validation"

	"goshop/internal/product/repository"
	"goshop/internal/product/service"
	"goshop/pkg/dbs"
	"goshop/pkg/middleware"
	"goshop/pkg/redis"
)

func Routes(r *gin.RouterGroup, db dbs.Database, validator validation.Validation, cache redis.Redis) {
	productRepo := repository.NewProductRepository(db)
	productSvc := service.NewProductService(validator, productRepo)
	productHandler := NewProductHandler(cache, productSvc)

	categoryRepo := repository.NewCategoryRepository(db)
	categorySvc := service.NewCategoryService(validator, categoryRepo)
	categoryHandler := NewCategoryHandler(categorySvc)

	reviewRepo := repository.NewReviewRepository(db)
	reviewSvc := service.NewReviewService(validator, reviewRepo, productRepo)
	reviewHandler := NewReviewHandler(reviewSvc)

	authMiddleware := middleware.JWTAuth()

	productRoute := r.Group("/products")
	{
		productRoute.GET("", productHandler.ListProducts)
		productRoute.POST("", authMiddleware, productHandler.CreateProduct)
		productRoute.PUT("/:id", authMiddleware, productHandler.UpdateProduct)
		productRoute.GET("/:id", productHandler.GetProductByID)
		productRoute.GET("/:id/reviews", reviewHandler.ListReviews)
		productRoute.POST("/:id/reviews", authMiddleware, reviewHandler.CreateReview)
		productRoute.PUT("/:id/reviews/:reviewId", authMiddleware, reviewHandler.UpdateReview)
		productRoute.DELETE("/:id/reviews/:reviewId", authMiddleware, reviewHandler.DeleteReview)
	}

	categoryRoute := r.Group("/categories")
	{
		categoryRoute.GET("", categoryHandler.ListCategories)
		categoryRoute.GET("/:id", categoryHandler.GetCategoryByID)
		categoryRoute.POST("", authMiddleware, categoryHandler.CreateCategory)
		categoryRoute.PUT("/:id", authMiddleware, categoryHandler.UpdateCategory)
		categoryRoute.DELETE("/:id", authMiddleware, categoryHandler.DeleteCategory)
	}
}
