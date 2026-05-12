package service

import (
	"context"
	"testing"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"goshop/internal/product/domain"
	"goshop/internal/product/model"
	"goshop/internal/product/repository/mocks"
	"goshop/pkg/config"
)

func newProductSvc(t *testing.T) (ProductService, *mocks.ProductRepository) {
	logger.Initialize(config.ProductionEnv)
	repo := mocks.NewProductRepository(t)
	return NewProductService(validation.New(), repo), repo
}

func TestProductService_Create_WithCategory(t *testing.T) {
	svc, repo := newProductSvc(t)
	repo.On("Create", mock.Anything, mock.MatchedBy(func(p *model.Product) bool {
		return p.CategoryID != nil && *p.CategoryID == "cat1"
	})).Return(nil).Once()

	_, err := svc.Create(context.Background(), &domain.CreateProductReq{
		Name: "x", Description: "x", Price: 10, StockQuantity: 1, CategoryID: "cat1",
	})
	require.NoError(t, err)
}

func TestProductService_Update_AllFieldsMutated(t *testing.T) {
	svc, repo := newProductSvc(t)
	repo.On("GetProductByID", mock.Anything, "p1").
		Return(&model.Product{ID: "p1", Name: "old", Price: 1}, nil).Once()
	repo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()

	qty := 50
	cid := "cat-new"
	_, err := svc.Update(context.Background(), "p1", &domain.UpdateProductReq{
		Name:          "new",
		Description:   "new",
		Price:         99,
		StockQuantity: &qty,
		Images:        []string{"img"},
		CategoryID:    cid,
	})
	require.NoError(t, err)
}

func TestCategoryService_Update_AllFieldsMutated(t *testing.T) {
	repo := mocks.NewCategoryRepository(t)
	svc := NewCategoryService(validation.New(), repo)
	repo.On("GetByID", mock.Anything, "c1").Return(&model.Category{ID: "c1"}, nil).Once()
	repo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
	_, err := svc.Update(context.Background(), "c1", &domain.UpdateCategoryReq{
		Name: "n", Slug: "s", Description: "d",
	})
	require.NoError(t, err)
}

func TestReviewService_Update_AllFieldsMutated(t *testing.T) {
	rrepo := mocks.NewReviewRepository(t)
	prepo := mocks.NewProductRepository(t)
	svc := NewReviewService(validation.New(), rrepo, prepo)
	rrepo.On("GetByID", mock.Anything, "r1").
		Return(&model.Review{ID: "r1", UserID: "u1", ProductID: "p1"}, nil).Once()
	rrepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
	rrepo.On("GetAggregates", mock.Anything, "p1").Return(float64(4), 1, nil).Maybe()
	prepo.On("UpdateRating", mock.Anything, "p1", float64(4), 1).Return(nil).Maybe()

	_, err := svc.UpdateReview(context.Background(), "r1", "u1", &domain.UpdateReviewReq{Rating: 5, Comment: "good"})
	require.NoError(t, err)
}
