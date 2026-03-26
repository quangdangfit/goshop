package service

import (
	"context"
	"errors"
	"testing"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"goshop/internal/product/dto"
	"goshop/internal/product/model"
	"goshop/internal/product/repository/mocks"
	"goshop/pkg/config"
	"goshop/pkg/paging"
)

type ProductServiceTestSuite struct {
	suite.Suite
	mockRepo *mocks.ProductRepository
	service  ProductService
}

func (suite *ProductServiceTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)

	validator := validation.New()
	suite.mockRepo = mocks.NewProductRepository(suite.T())
	suite.service = NewProductService(validator, suite.mockRepo)
}

func TestProductServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ProductServiceTestSuite))
}

func (suite *ProductServiceTestSuite) TestGetProductByID() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockRepo.On("GetProductByID", mock.Anything, "productID").
					Return(&model.Product{Name: "product", Description: "product description", Price: 1.1}, nil).Times(1)
			},
		},
		{
			name: "Not found",
			setup: func() {
				suite.mockRepo.On("GetProductByID", mock.Anything, "productID").
					Return(nil, errors.New("error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			product, err := suite.service.GetProductByID(context.Background(), "productID")
			if tc.wantErr {
				suite.Nil(product)
				suite.NotNil(err)
			} else {
				suite.NotNil(product)
				suite.Equal("product", product.Name)
				suite.Equal("product description", product.Description)
				suite.Equal(1.1, product.Price)
				suite.Nil(err)
			}
		})
	}
}

func (suite *ProductServiceTestSuite) TestListProducts() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockRepo.On("ListProducts", mock.Anything, mock.Anything).
					Return([]*model.Product{{Name: "product", Description: "product description", Price: 1.1}},
						&paging.Pagination{Total: 1, CurrentPage: 1, Limit: 10}, nil).Times(1)
			},
		},
		{
			name: "Fail",
			setup: func() {
				suite.mockRepo.On("ListProducts", mock.Anything, mock.Anything).
					Return(nil, nil, errors.New("error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			req := &dto.ListProductReq{Name: "product"}
			products, pagination, err := suite.service.ListProducts(context.Background(), req)
			if tc.wantErr {
				suite.Nil(products)
				suite.Nil(pagination)
				suite.NotNil(err)
			} else {
				suite.NotNil(products)
				suite.Equal(1, len(products))
				suite.NotNil(pagination)
				suite.Equal(int64(1), pagination.Total)
				suite.Nil(err)
			}
		})
	}
}

func (suite *ProductServiceTestSuite) TestCreate() {
	tests := []struct {
		name    string
		req     *dto.CreateProductReq
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			req:  &dto.CreateProductReq{Name: "product", Description: "product description", Price: 1.1},
			setup: func() {
				suite.mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name: "DB fail",
			req:  &dto.CreateProductReq{Name: "product", Description: "product description", Price: 1.1},
			setup: func() {
				suite.mockRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("error")).Times(1)
			},
			wantErr: true,
		},
		{
			name:    "Missing product name",
			req:     &dto.CreateProductReq{Description: "product description", Price: 1.1},
			setup:   func() {},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			product, err := suite.service.Create(context.Background(), tc.req)
			if tc.wantErr {
				suite.Nil(product)
				suite.NotNil(err)
			} else {
				suite.NotNil(product)
				suite.Equal(tc.req.Name, product.Name)
				suite.Equal(tc.req.Description, product.Description)
				suite.Equal(tc.req.Price, product.Price)
				suite.Nil(err)
			}
		})
	}
}

func (suite *ProductServiceTestSuite) TestUpdate() {
	tests := []struct {
		name    string
		req     *dto.UpdateProductReq
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			req:  &dto.UpdateProductReq{Name: "product", Description: "product description", Price: 1.1},
			setup: func() {
				suite.mockRepo.On("GetProductByID", mock.Anything, "productID").
					Return(&model.Product{Name: "product", Description: "product description", Price: 1.1}, nil).Times(1)
				suite.mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name: "Update DB fail",
			req:  &dto.UpdateProductReq{Name: "product", Description: "product description", Price: 1.1},
			setup: func() {
				suite.mockRepo.On("GetProductByID", mock.Anything, "productID").
					Return(&model.Product{Name: "product", Description: "product description", Price: 1.1}, nil).Times(1)
				suite.mockRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("error")).Times(1)
			},
			wantErr: true,
		},
		{
			name:    "Invalid price",
			req:     &dto.UpdateProductReq{Name: "product", Description: "product description", Price: -1.1},
			setup:   func() {},
			wantErr: true,
		},
		{
			name: "GetProductByID fail",
			req:  &dto.UpdateProductReq{Name: "product", Description: "product description", Price: 1.1},
			setup: func() {
				suite.mockRepo.On("GetProductByID", mock.Anything, "productID").
					Return(nil, errors.New("error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			product, err := suite.service.Update(context.Background(), "productID", tc.req)
			if tc.wantErr {
				suite.Nil(product)
				suite.NotNil(err)
			} else {
				suite.NotNil(product)
				suite.Equal(tc.req.Name, product.Name)
				suite.Equal(tc.req.Description, product.Description)
				suite.Equal(tc.req.Price, product.Price)
				suite.Nil(err)
			}
		})
	}
}
