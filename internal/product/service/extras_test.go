package service

import (
	"context"
	"errors"
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

func newSvcExtras(t *testing.T) (ProductService, *mocks.ProductRepository) {
	logger.Initialize(config.ProductionEnv)
	repo := mocks.NewProductRepository(t)
	return NewProductService(validation.New(), repo), repo
}

func TestAddStock_RejectsZeroOrNegative(t *testing.T) {
	svc, _ := newSvcExtras(t)
	for _, q := range []int{0, -1, -100} {
		_, err := svc.AddStock(context.Background(), "p1", q, "admin")
		require.ErrorIs(t, err, errInvalidStockQty)
	}
}

func TestAddStock_RepoError(t *testing.T) {
	svc, repo := newSvcExtras(t)
	repo.On("AddStock", mock.Anything, "p1", 10).Return(errors.New("db")).Once()
	_, err := svc.AddStock(context.Background(), "p1", 10, "admin")
	require.Error(t, err)
}

func TestAddStock_Success(t *testing.T) {
	svc, repo := newSvcExtras(t)
	repo.On("AddStock", mock.Anything, "p1", 10).Return(nil).Once()
	repo.On("GetProductByID", mock.Anything, "p1").
		Return(&model.Product{ID: "p1", StockQuantity: 110}, nil).Once()
	got, err := svc.AddStock(context.Background(), "p1", 10, "admin-id")
	require.NoError(t, err)
	require.Equal(t, 110, got.StockQuantity)
}

func TestUpdate_GetError(t *testing.T) {
	svc, repo := newSvcExtras(t)
	repo.On("GetProductByID", mock.Anything, "p1").Return(nil, errors.New("not found")).Once()
	_, err := svc.Update(context.Background(), "p1", &domain.UpdateProductReq{Name: "x"})
	require.Error(t, err)
}
