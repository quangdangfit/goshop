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
	mockRepo *mocks.IProductRepository
	service  IProductService
}

func (suite *ProductServiceTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)

	validator := validation.New()
	suite.mockRepo = mocks.NewIProductRepository(suite.T())
	suite.service = NewProductService(validator, suite.mockRepo)
}

func TestProductServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ProductServiceTestSuite))
}

// GetProductByID
// =================================================================

func (suite *ProductServiceTestSuite) TestGetProductByIDSuccess() {
	productID := "productID"

	suite.mockRepo.On("GetProductByID", mock.Anything, productID).
		Return(
			&model.Product{
				Name:        "product",
				Description: "product description",
				Price:       1.1,
			},
			nil,
		).Times(1)

	product, err := suite.service.GetProductByID(context.Background(), productID)
	suite.NotNil(product)
	suite.Equal("product", product.Name)
	suite.Equal("product description", product.Description)
	suite.Equal(1.1, product.Price)
	suite.Nil(err)
}

func (suite *ProductServiceTestSuite) TestGetProductByIDFail() {
	productID := "productID"
	suite.mockRepo.On("GetProductByID", mock.Anything, productID).
		Return(nil, errors.New("error")).Times(1)

	product, err := suite.service.GetProductByID(context.Background(), productID)
	suite.Nil(product)
	suite.NotNil(err)
}

// ListProducts
// =================================================================

func (suite *ProductServiceTestSuite) TestListProductsSuccess() {
	req := &dto.ListProductReq{
		Name: "product",
	}

	suite.mockRepo.On("ListProducts", mock.Anything, req).
		Return(
			[]*model.Product{
				{
					Name:        "product",
					Description: "product description",
					Price:       1.1,
				},
			},
			&paging.Pagination{
				Total:       1,
				CurrentPage: 1,
				Limit:       10,
			},
			nil,
		).Times(1)

	products, pagination, err := suite.service.ListProducts(context.Background(), req)
	suite.NotNil(products)
	suite.Equal(1, len(products))
	suite.Equal("product", products[0].Name)
	suite.Equal("product description", products[0].Description)
	suite.Equal(1.1, products[0].Price)
	suite.NotNil(pagination)
	suite.Equal(int64(1), pagination.Total)
	suite.Equal(int64(1), pagination.CurrentPage)
	suite.Equal(int64(10), pagination.Limit)
	suite.Nil(err)
}

func (suite *ProductServiceTestSuite) TestListProductsFail() {
	req := &dto.ListProductReq{
		Name: "product",
	}

	suite.mockRepo.On("ListProducts", mock.Anything, req).
		Return(nil, nil, errors.New("error")).Times(1)

	products, pagination, err := suite.service.ListProducts(context.Background(), req)
	suite.Nil(products)
	suite.Nil(pagination)
	suite.NotNil(err)
}

// Create
// =================================================================

func (suite *ProductServiceTestSuite) TestCreateSuccess() {
	req := &dto.CreateProductReq{
		Name:        "product",
		Description: "product description",
		Price:       1.1,
	}

	suite.mockRepo.On("Create", mock.Anything, &model.Product{
		Name:        "product",
		Description: "product description",
		Price:       1.1,
	}).Return(nil).Times(1)

	product, err := suite.service.Create(context.Background(), req)
	suite.NotNil(product)
	suite.Equal(req.Name, product.Name)
	suite.Equal(req.Description, product.Description)
	suite.Equal(req.Price, product.Price)
	suite.Nil(err)
}

func (suite *ProductServiceTestSuite) TestCreateFail() {
	req := &dto.CreateProductReq{
		Name:        "product",
		Description: "product description",
		Price:       1.1,
	}

	suite.mockRepo.On("Create", mock.Anything, &model.Product{
		Name:        "product",
		Description: "product description",
		Price:       1.1,
	}).Return(errors.New("error")).Times(1)

	product, err := suite.service.Create(context.Background(), req)
	suite.Nil(product)
	suite.NotNil(err)
}

func (suite *ProductServiceTestSuite) TestCreateMissProductName() {
	req := &dto.CreateProductReq{
		Description: "product description",
		Price:       1.1,
	}

	product, err := suite.service.Create(context.Background(), req)
	suite.Nil(product)
	suite.NotNil(err)
}

// Update
// =================================================================

func (suite *ProductServiceTestSuite) TestUpdateSuccess() {
	productID := "productID"
	req := &dto.UpdateProductReq{
		Name:        "product",
		Description: "product description",
		Price:       1.1,
	}

	suite.mockRepo.On("GetProductByID", mock.Anything, productID).
		Return(&model.Product{
			Name:        "product",
			Description: "product description",
			Price:       1.1,
		},
			nil).Times(1)

	suite.mockRepo.On("Update", mock.Anything, &model.Product{
		Name:        "product",
		Description: "product description",
		Price:       1.1,
	}).Return(nil).Times(1)

	product, err := suite.service.Update(context.Background(), productID, req)
	suite.NotNil(product)
	suite.Equal(req.Name, product.Name)
	suite.Equal(req.Description, product.Description)
	suite.Equal(req.Price, product.Price)
	suite.Nil(err)
}

func (suite *ProductServiceTestSuite) TestUpdateFail() {
	productID := "productID"
	req := &dto.UpdateProductReq{
		Name:        "product",
		Description: "product description",
		Price:       1.1,
	}

	suite.mockRepo.On("GetProductByID", mock.Anything, productID).
		Return(&model.Product{
			Name:        "product",
			Description: "product description",
			Price:       1.1,
		},
			nil).Times(1)

	suite.mockRepo.On("Update", mock.Anything, &model.Product{
		Name:        "product",
		Description: "product description",
		Price:       1.1,
	}).Return(errors.New("error")).Times(1)

	product, err := suite.service.Update(context.Background(), productID, req)
	suite.Nil(product)
	suite.NotNil(err)
}

func (suite *ProductServiceTestSuite) TestUpdateInvalidPrice() {
	productID := "productID"
	req := &dto.UpdateProductReq{
		Name:        "product",
		Description: "product description",
		Price:       -1.1,
	}

	product, err := suite.service.Update(context.Background(), productID, req)
	suite.Nil(product)
	suite.NotNil(err)
}

func (suite *ProductServiceTestSuite) TestUpdateGetProductByIDFail() {
	productID := "productID"
	req := &dto.UpdateProductReq{
		Name:        "product",
		Description: "product description",
		Price:       1.1,
	}

	suite.mockRepo.On("GetProductByID", mock.Anything, productID).
		Return(nil, errors.New("error")).Times(1)

	product, err := suite.service.Update(context.Background(), productID, req)
	suite.Nil(product)
	suite.NotNil(err)
}
