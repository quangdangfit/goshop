package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"goshop/internal/product/model"
	"goshop/internal/product/service/mocks"
	"goshop/pkg/config"
	"goshop/pkg/paging"
	pb "goshop/proto/gen/go/product"
)

type ProductHandlerTestSuite struct {
	suite.Suite
	mockService *mocks.ProductService
	handler     *ProductHandler
}

func (suite *ProductHandlerTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)

	suite.mockService = mocks.NewProductService(suite.T())
	suite.handler = NewProductHandler(suite.mockService)
}

func TestProductHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ProductHandlerTestSuite))
}

// GetProductByID
// =================================================================================================

func (suite *ProductHandlerTestSuite) TestGetProductByIDSuccess() {
	req := &pb.GetProductByIDReq{Id: "productId1"}

	suite.mockService.On("GetProductByID", mock.Anything, "productId1").
		Return(&model.Product{
			ID:          "productId1",
			Name:        "product",
			Description: "description",
			Price:       10.5,
		}, nil).Times(1)

	res, err := suite.handler.GetProductByID(context.Background(), req)

	suite.Nil(err)
	suite.NotNil(res)
	suite.Equal("productId1", res.Product.Id)
	suite.Equal("product", res.Product.Name)
	suite.Equal("description", res.Product.Description)
	suite.Equal(float32(10.5), res.Product.Price)
}

func (suite *ProductHandlerTestSuite) TestGetProductByIDMissID() {
	req := &pb.GetProductByIDReq{}

	res, err := suite.handler.GetProductByID(context.Background(), req)

	suite.Nil(res)
	suite.NotNil(err)
}

func (suite *ProductHandlerTestSuite) TestGetProductByIDFail() {
	req := &pb.GetProductByIDReq{Id: "productId1"}

	suite.mockService.On("GetProductByID", mock.Anything, "productId1").
		Return(nil, errors.New("error")).Times(1)

	res, err := suite.handler.GetProductByID(context.Background(), req)

	suite.Nil(res)
	suite.NotNil(err)
}

// ListProducts
// =================================================================================================

func (suite *ProductHandlerTestSuite) TestListProductsSuccess() {
	req := &pb.ListProductsReq{
		Name:  "product",
		Page:  1,
		Limit: 10,
	}

	suite.mockService.On("ListProducts", mock.Anything, mock.Anything).
		Return(
			[]*model.Product{
				{
					ID:          "productId1",
					Name:        "product",
					Description: "description",
					Price:       10.5,
				},
			},
			&paging.Pagination{
				Total:       1,
				CurrentPage: 1,
				Limit:       10,
			},
			nil,
		).Times(1)

	res, err := suite.handler.ListProducts(context.Background(), req)

	suite.Nil(err)
	suite.NotNil(res)
	suite.Equal(1, len(res.Products))
	suite.Equal("productId1", res.Products[0].Id)
	suite.Equal(int64(1), res.Total)
	suite.Equal(int64(1), res.CurrentPage)
	suite.Equal(int64(10), res.Limit)
}

func (suite *ProductHandlerTestSuite) TestListProductsFail() {
	req := &pb.ListProductsReq{}

	suite.mockService.On("ListProducts", mock.Anything, mock.Anything).
		Return(nil, nil, errors.New("error")).Times(1)

	res, err := suite.handler.ListProducts(context.Background(), req)

	suite.Nil(res)
	suite.NotNil(err)
}

// CreateProduct
// =================================================================================================

func (suite *ProductHandlerTestSuite) TestCreateProductSuccess() {
	req := &pb.CreateProductReq{
		Name:        "product",
		Description: "description",
		Price:       10.5,
	}

	suite.mockService.On("Create", mock.Anything, mock.Anything).
		Return(&model.Product{
			ID:          "productId1",
			Name:        "product",
			Description: "description",
			Price:       10.5,
		}, nil).Times(1)

	res, err := suite.handler.CreateProduct(context.Background(), req)

	suite.Nil(err)
	suite.NotNil(res)
	suite.Equal("productId1", res.Product.Id)
	suite.Equal("product", res.Product.Name)
	suite.Equal("description", res.Product.Description)
}

func (suite *ProductHandlerTestSuite) TestCreateProductFail() {
	req := &pb.CreateProductReq{
		Name:        "product",
		Description: "description",
		Price:       10.5,
	}

	suite.mockService.On("Create", mock.Anything, mock.Anything).
		Return(nil, errors.New("error")).Times(1)

	res, err := suite.handler.CreateProduct(context.Background(), req)

	suite.Nil(res)
	suite.NotNil(err)
}

// UpdateProduct
// =================================================================================================

func (suite *ProductHandlerTestSuite) TestUpdateProductSuccess() {
	req := &pb.UpdateProductReq{
		Id:          "productId1",
		Name:        "updated",
		Description: "updated description",
		Price:       20.0,
	}

	suite.mockService.On("Update", mock.Anything, "productId1", mock.Anything).
		Return(&model.Product{
			ID:          "productId1",
			Name:        "updated",
			Description: "updated description",
			Price:       20.0,
		}, nil).Times(1)

	res, err := suite.handler.UpdateProduct(context.Background(), req)

	suite.Nil(err)
	suite.NotNil(res)
	suite.Equal("productId1", res.Product.Id)
	suite.Equal("updated", res.Product.Name)
	suite.Equal(float32(20.0), res.Product.Price)
}

func (suite *ProductHandlerTestSuite) TestUpdateProductMissID() {
	req := &pb.UpdateProductReq{Name: "updated"}

	res, err := suite.handler.UpdateProduct(context.Background(), req)

	suite.Nil(res)
	suite.NotNil(err)
}

func (suite *ProductHandlerTestSuite) TestUpdateProductFail() {
	req := &pb.UpdateProductReq{
		Id:    "productId1",
		Name:  "updated",
		Price: 20.0,
	}

	suite.mockService.On("Update", mock.Anything, "productId1", mock.Anything).
		Return(nil, errors.New("error")).Times(1)

	res, err := suite.handler.UpdateProduct(context.Background(), req)

	suite.Nil(res)
	suite.NotNil(err)
}
