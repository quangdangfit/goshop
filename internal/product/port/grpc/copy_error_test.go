package grpc

import (
	"context"
	"math"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"goshop/internal/product/model"
	srvMocks "goshop/internal/product/service/mocks"
	"goshop/pkg/paging"
	pb "goshop/proto/gen/go/product"
)

func nanProduct() *model.Product {
	return &model.Product{ID: "p1", Price: math.NaN()}
}

func TestProductGRPC_GetProductByID_CopyError(t *testing.T) {
	svc := srvMocks.NewProductService(t)
	svc.On("GetProductByID", mock.Anything, "p1").Return(nanProduct(), nil).Once()
	h := NewProductHandler(svc)
	_, err := h.GetProductByID(context.Background(), &pb.GetProductByIDReq{Id: "p1"})
	require.Error(t, err)
}

func TestProductGRPC_ListProducts_CopyError(t *testing.T) {
	svc := srvMocks.NewProductService(t)
	svc.On("ListProducts", mock.Anything, mock.Anything).
		Return([]*model.Product{nanProduct()}, &paging.Pagination{}, nil).Once()
	h := NewProductHandler(svc)
	_, err := h.ListProducts(context.Background(), &pb.ListProductsReq{})
	require.Error(t, err)
}

func TestProductGRPC_CreateProduct_CopyError(t *testing.T) {
	svc := srvMocks.NewProductService(t)
	svc.On("Create", mock.Anything, mock.Anything).Return(nanProduct(), nil).Once()
	h := NewProductHandler(svc)
	_, err := h.CreateProduct(context.Background(), &pb.CreateProductReq{
		Name: "x", Description: "x", Price: 1,
	})
	require.Error(t, err)
}

func TestProductGRPC_UpdateProduct_CopyError(t *testing.T) {
	svc := srvMocks.NewProductService(t)
	svc.On("Update", mock.Anything, "p1", mock.Anything).Return(nanProduct(), nil).Once()
	h := NewProductHandler(svc)
	_, err := h.UpdateProduct(context.Background(), &pb.UpdateProductReq{Id: "p1", Name: "x"})
	require.Error(t, err)
}
