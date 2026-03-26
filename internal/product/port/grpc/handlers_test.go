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

func (suite *ProductHandlerTestSuite) TestGetProductByID() {
	tests := []struct {
		name      string
		setup     func()
		req       *pb.GetProductByIDReq
		expectNil bool
		expectErr bool
		validate  func(res *pb.GetProductByIDRes)
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockService.On("GetProductByID", mock.Anything, "productId1").
					Return(&model.Product{
						ID:          "productId1",
						Name:        "product",
						Description: "description",
						Price:       10.5,
					}, nil).Times(1)
			},
			req:       &pb.GetProductByIDReq{Id: "productId1"},
			expectNil: false,
			expectErr: false,
			validate: func(res *pb.GetProductByIDRes) {
				suite.Equal("productId1", res.Product.Id)
				suite.Equal("product", res.Product.Name)
				suite.Equal("description", res.Product.Description)
				suite.Equal(float32(10.5), res.Product.Price)
			},
		},
		{
			name:      "MissID",
			setup:     func() {},
			req:       &pb.GetProductByIDReq{},
			expectNil: true,
			expectErr: true,
		},
		{
			name: "Fail",
			setup: func() {
				suite.mockService.On("GetProductByID", mock.Anything, "productId1").
					Return(nil, errors.New("error")).Times(1)
			},
			req:       &pb.GetProductByIDReq{Id: "productId1"},
			expectNil: true,
			expectErr: true,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()

			res, err := suite.handler.GetProductByID(context.Background(), tc.req)

			if tc.expectErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
			if tc.expectNil {
				suite.Nil(res)
			} else {
				suite.NotNil(res)
				if tc.validate != nil {
					tc.validate(res)
				}
			}
		})
	}
}

// ListProducts
// =================================================================================================

func (suite *ProductHandlerTestSuite) TestListProducts() {
	tests := []struct {
		name      string
		setup     func()
		req       *pb.ListProductsReq
		expectNil bool
		expectErr bool
		validate  func(res *pb.ListProductsRes)
	}{
		{
			name: "Success",
			setup: func() {
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
			},
			req: &pb.ListProductsReq{
				Name:  "product",
				Page:  1,
				Limit: 10,
			},
			expectNil: false,
			expectErr: false,
			validate: func(res *pb.ListProductsRes) {
				suite.Equal(1, len(res.Products))
				suite.Equal("productId1", res.Products[0].Id)
				suite.Equal(int64(1), res.Total)
				suite.Equal(int64(1), res.CurrentPage)
				suite.Equal(int64(10), res.Limit)
			},
		},
		{
			name: "Fail",
			setup: func() {
				suite.mockService.On("ListProducts", mock.Anything, mock.Anything).
					Return(nil, nil, errors.New("error")).Times(1)
			},
			req:       &pb.ListProductsReq{},
			expectNil: true,
			expectErr: true,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()

			res, err := suite.handler.ListProducts(context.Background(), tc.req)

			if tc.expectErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
			if tc.expectNil {
				suite.Nil(res)
			} else {
				suite.NotNil(res)
				if tc.validate != nil {
					tc.validate(res)
				}
			}
		})
	}
}

// CreateProduct
// =================================================================================================

func (suite *ProductHandlerTestSuite) TestCreateProduct() {
	tests := []struct {
		name      string
		setup     func()
		req       *pb.CreateProductReq
		expectNil bool
		expectErr bool
		validate  func(res *pb.CreateProductRes)
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockService.On("Create", mock.Anything, mock.Anything).
					Return(&model.Product{
						ID:          "productId1",
						Name:        "product",
						Description: "description",
						Price:       10.5,
					}, nil).Times(1)
			},
			req: &pb.CreateProductReq{
				Name:        "product",
				Description: "description",
				Price:       10.5,
			},
			expectNil: false,
			expectErr: false,
			validate: func(res *pb.CreateProductRes) {
				suite.Equal("productId1", res.Product.Id)
				suite.Equal("product", res.Product.Name)
				suite.Equal("description", res.Product.Description)
			},
		},
		{
			name: "Fail",
			setup: func() {
				suite.mockService.On("Create", mock.Anything, mock.Anything).
					Return(nil, errors.New("error")).Times(1)
			},
			req: &pb.CreateProductReq{
				Name:        "product",
				Description: "description",
				Price:       10.5,
			},
			expectNil: true,
			expectErr: true,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()

			res, err := suite.handler.CreateProduct(context.Background(), tc.req)

			if tc.expectErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
			if tc.expectNil {
				suite.Nil(res)
			} else {
				suite.NotNil(res)
				if tc.validate != nil {
					tc.validate(res)
				}
			}
		})
	}
}

// UpdateProduct
// =================================================================================================

func (suite *ProductHandlerTestSuite) TestUpdateProduct() {
	tests := []struct {
		name      string
		setup     func()
		req       *pb.UpdateProductReq
		expectNil bool
		expectErr bool
		validate  func(res *pb.UpdateProductRes)
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockService.On("Update", mock.Anything, "productId1", mock.Anything).
					Return(&model.Product{
						ID:          "productId1",
						Name:        "updated",
						Description: "updated description",
						Price:       20.0,
					}, nil).Times(1)
			},
			req: &pb.UpdateProductReq{
				Id:          "productId1",
				Name:        "updated",
				Description: "updated description",
				Price:       20.0,
			},
			expectNil: false,
			expectErr: false,
			validate: func(res *pb.UpdateProductRes) {
				suite.Equal("productId1", res.Product.Id)
				suite.Equal("updated", res.Product.Name)
				suite.Equal(float32(20.0), res.Product.Price)
			},
		},
		{
			name:      "MissID",
			setup:     func() {},
			req:       &pb.UpdateProductReq{Name: "updated"},
			expectNil: true,
			expectErr: true,
		},
		{
			name: "Fail",
			setup: func() {
				suite.mockService.On("Update", mock.Anything, "productId1", mock.Anything).
					Return(nil, errors.New("error")).Times(1)
			},
			req: &pb.UpdateProductReq{
				Id:    "productId1",
				Name:  "updated",
				Price: 20.0,
			},
			expectNil: true,
			expectErr: true,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()

			res, err := suite.handler.UpdateProduct(context.Background(), tc.req)

			if tc.expectErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
			if tc.expectNil {
				suite.Nil(res)
			} else {
				suite.NotNil(res)
				if tc.validate != nil {
					tc.validate(res)
				}
			}
		})
	}
}
