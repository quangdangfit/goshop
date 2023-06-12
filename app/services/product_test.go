package services

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/stretchr/testify/suite"

	"goshop/app/models"
	"goshop/app/serializers"
	"goshop/config"
	"goshop/mocks"
	"goshop/pkg/paging"
)

type ProductServiceTestSuite struct {
	suite.Suite
	mockRepo *mocks.MockIProductRepository
	service  IProductService
}

func (suite *ProductServiceTestSuite) SetupTest() {
	logger.Initialize(config.TestEnv)

	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()
	suite.mockRepo = mocks.NewMockIProductRepository(mockCtrl)
	suite.service = NewProductService(suite.mockRepo)
}

func TestProductServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ProductServiceTestSuite))
}

// GetProductByID
// =================================================================

func (suite *ProductServiceTestSuite) TestGetProductByIDSuccess() {
	productID := "productID"
	suite.mockRepo.EXPECT().GetProductByID(gomock.Any(), productID).Return(&models.Product{
		Name:        "product",
		Description: "product description",
		Price:       1.1,
	}, nil).Times(1)

	product, err := suite.service.GetProductByID(context.Background(), productID)
	suite.NotNil(product)
	suite.Equal("product", product.Name)
	suite.Equal("product description", product.Description)
	suite.Equal(1.1, product.Price)
	suite.Nil(err)
}

func (suite *ProductServiceTestSuite) TestGetProductByIDFail() {
	productID := "productID"
	suite.mockRepo.EXPECT().GetProductByID(gomock.Any(), productID).Return(nil, errors.New("error")).Times(1)

	product, err := suite.service.GetProductByID(context.Background(), productID)
	suite.Nil(product)
	suite.NotNil(err)
}

// ListProducts
// =================================================================

func (suite *ProductServiceTestSuite) TestListProductsSuccess() {
	req := &serializers.ListProductReq{
		Name: "product",
	}

	suite.mockRepo.EXPECT().ListProducts(gomock.Any(), req).Return(
		[]*models.Product{
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
		nil).Times(1)

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
	req := &serializers.ListProductReq{
		Name: "product",
	}

	suite.mockRepo.EXPECT().ListProducts(gomock.Any(), req).Return(nil, nil, errors.New("error")).Times(1)

	products, pagination, err := suite.service.ListProducts(context.Background(), req)
	suite.Nil(products)
	suite.Nil(pagination)
	suite.NotNil(err)
}

// Create
// =================================================================

func (suite *ProductServiceTestSuite) TestCreateSuccess() {
	req := &serializers.CreateProductReq{
		Name:        "product",
		Description: "product description",
		Price:       1.1,
	}

	suite.mockRepo.EXPECT().Create(gomock.Any(), &models.Product{
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
	req := &serializers.CreateProductReq{
		Name:        "product",
		Description: "product description",
		Price:       1.1,
	}

	suite.mockRepo.EXPECT().Create(gomock.Any(), &models.Product{
		Name:        "product",
		Description: "product description",
		Price:       1.1,
	}).Return(errors.New("error")).Times(1)

	product, err := suite.service.Create(context.Background(), req)
	suite.Nil(product)
	suite.NotNil(err)
}

// Update
// =================================================================

func (suite *ProductServiceTestSuite) TestUpdateSuccess() {
	productID := "productID"
	req := &serializers.UpdateProductReq{
		Name:        "product",
		Description: "product description",
		Price:       1.1,
	}

	suite.mockRepo.EXPECT().GetProductByID(gomock.Any(), productID).Return(
		&models.Product{
			Name:        "product",
			Description: "product description",
			Price:       1.1,
		},
		nil).Times(1)
	suite.mockRepo.EXPECT().Update(gomock.Any(), &models.Product{
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
	req := &serializers.UpdateProductReq{
		Name:        "product",
		Description: "product description",
		Price:       1.1,
	}

	suite.mockRepo.EXPECT().GetProductByID(gomock.Any(), productID).Return(
		&models.Product{
			Name:        "product",
			Description: "product description",
			Price:       1.1,
		},
		nil).Times(1)
	suite.mockRepo.EXPECT().Update(gomock.Any(), &models.Product{
		Name:        "product",
		Description: "product description",
		Price:       1.1,
	}).Return(errors.New("error")).Times(1)

	product, err := suite.service.Update(context.Background(), productID, req)
	suite.Nil(product)
	suite.NotNil(err)
}

func (suite *ProductServiceTestSuite) TestUpdateGetProductByIDFail() {
	productID := "productID"
	req := &serializers.UpdateProductReq{
		Name:        "product",
		Description: "product description",
		Price:       1.1,
	}

	suite.mockRepo.EXPECT().GetProductByID(gomock.Any(), productID).Return(nil, errors.New("error")).Times(1)

	product, err := suite.service.Update(context.Background(), productID, req)
	suite.Nil(product)
	suite.NotNil(err)
}
