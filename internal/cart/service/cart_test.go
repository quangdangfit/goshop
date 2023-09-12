package service

import (
	"context"
	"errors"
	"testing"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"goshop/internal/cart/dto"
	"goshop/internal/cart/model"
	"goshop/internal/cart/repository/mocks"
	"goshop/pkg/config"
)

type CartServiceTestSuite struct {
	suite.Suite
	mockRepo *mocks.ICartRepository
	service  ICartService
}

func (suite *CartServiceTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)

	validator := validation.New()
	suite.mockRepo = mocks.NewICartRepository(suite.T())
	suite.service = NewCartService(validator, suite.mockRepo)
}

func TestCartServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CartServiceTestSuite))
}

// GetCartByUserID
// =================================================================

func (suite *CartServiceTestSuite) TestGetCartByUserIDSuccessfully() {
	userID := "userID"

	suite.mockRepo.On("GetCartByUserID", mock.Anything, userID).
		Return(
			&model.Cart{
				ID:     "cartId1",
				UserID: "userID",
				Lines: []*model.CartLine{
					{
						ProductID: "productID1",
						Quantity:  4,
					},
					{
						ProductID: "productID2",
						Quantity:  3,
					},
				},
			},
			nil,
		).Times(1)

	cart, err := suite.service.GetCartByUserID(context.Background(), userID)
	suite.NotNil(cart)
	suite.Equal("userID", cart.UserID)
	suite.Equal(2, len(cart.Lines))
	suite.Nil(err)
}

func (suite *CartServiceTestSuite) TestGetCartByUserIDFail() {
	userID := "userID"
	suite.mockRepo.On("GetCartByUserID", mock.Anything, userID).
		Return(nil, errors.New("error")).Times(1)

	suite.mockRepo.On("Create", mock.Anything, &model.Cart{
		UserID: "userID",
	}).Return(nil).Times(1)

	cart, err := suite.service.GetCartByUserID(context.Background(), userID)
	suite.NotNil(cart)
	suite.Equal("userID", cart.UserID)
	suite.Equal(0, len(cart.Lines))
	suite.Nil(err)
}

func (suite *CartServiceTestSuite) TestGetCartByUserIDCreateFail() {
	userID := "userID"
	suite.mockRepo.On("GetCartByUserID", mock.Anything, userID).
		Return(nil, errors.New("error")).Times(1)

	suite.mockRepo.On("Create", mock.Anything, &model.Cart{
		UserID: "userID",
	}).Return(errors.New("error")).Times(1)

	cart, err := suite.service.GetCartByUserID(context.Background(), userID)
	suite.Nil(cart)
	suite.NotNil(err)
}

// AddProduct
// =================================================================

func (suite *CartServiceTestSuite) TestAddProductSuccessfully() {
	req := &dto.AddProductReq{
		UserID: "userID",
		Line: &dto.CartLineReq{
			ProductID: "productID2",
			Quantity:  3,
		},
	}

	suite.mockRepo.On("GetCartByUserID", mock.Anything, "userID").
		Return(
			&model.Cart{
				ID:     "cartId1",
				UserID: "userID",
				Lines: []*model.CartLine{
					{
						ProductID: "productID1",
						Quantity:  4,
					},
				},
			},
			nil,
		).Times(1)

	suite.mockRepo.On("Update", mock.Anything, &model.Cart{
		ID:     "cartId1",
		UserID: "userID",
		Lines: []*model.CartLine{
			{
				ProductID: "productID1",
				Quantity:  4,
			},
			{
				ProductID: "productID2",
				Quantity:  3,
			},
		},
	}).Return(nil).Times(1)

	cart, err := suite.service.AddProduct(context.Background(), req)
	suite.NotNil(cart)
	suite.Equal("userID", cart.UserID)
	suite.Equal(2, len(cart.Lines))
	suite.Nil(err)
}

func (suite *CartServiceTestSuite) TestAddProductMissUserID() {
	req := &dto.AddProductReq{
		Line: &dto.CartLineReq{
			ProductID: "productID2",
			Quantity:  3,
		},
	}

	cart, err := suite.service.AddProduct(context.Background(), req)
	suite.Nil(cart)
	suite.NotNil(err)
}

func (suite *CartServiceTestSuite) TestAddProductMissProductID() {
	req := &dto.AddProductReq{
		UserID: "userID",
		Line: &dto.CartLineReq{
			Quantity: 3,
		},
	}

	cart, err := suite.service.AddProduct(context.Background(), req)
	suite.Nil(cart)
	suite.NotNil(err)
}

func (suite *CartServiceTestSuite) TestAddProductMissQuantity() {
	req := &dto.AddProductReq{
		UserID: "userID",
		Line: &dto.CartLineReq{
			ProductID: "productID2",
		},
	}

	cart, err := suite.service.AddProduct(context.Background(), req)
	suite.Nil(cart)
	suite.NotNil(err)
}

func (suite *CartServiceTestSuite) TestAddProductCartNotFound() {
	req := &dto.AddProductReq{
		UserID: "userID",
		Line: &dto.CartLineReq{
			ProductID: "productID2",
			Quantity:  3,
		},
	}

	suite.mockRepo.On("GetCartByUserID", mock.Anything, "userID").
		Return(nil, errors.New("error")).Times(1)

	suite.mockRepo.On("Create", mock.Anything, &model.Cart{
		UserID: "userID",
		Lines: []*model.CartLine{
			{
				ProductID: "productID2",
				Quantity:  3,
			},
		},
	}).Return(nil).Times(1)

	cart, err := suite.service.AddProduct(context.Background(), req)
	suite.NotNil(cart)
	suite.Equal("userID", cart.UserID)
	suite.Equal(1, len(cart.Lines))
	suite.Nil(err)
}

func (suite *CartServiceTestSuite) TestAddProductCartNotFoundCreateFail() {
	req := &dto.AddProductReq{
		UserID: "userID",
		Line: &dto.CartLineReq{
			ProductID: "productID2",
			Quantity:  3,
		},
	}

	suite.mockRepo.On("GetCartByUserID", mock.Anything, "userID").
		Return(nil, errors.New("error")).Times(1)

	suite.mockRepo.On("Create", mock.Anything, &model.Cart{
		UserID: "userID",
		Lines: []*model.CartLine{
			{
				ProductID: "productID2",
				Quantity:  3,
			},
		},
	}).Return(errors.New("error")).Times(1)

	cart, err := suite.service.AddProduct(context.Background(), req)
	suite.Nil(cart)
	suite.NotNil(err)
}

func (suite *CartServiceTestSuite) TestAddProductAlreadyExistInCart() {
	req := &dto.AddProductReq{
		UserID: "userID",
		Line: &dto.CartLineReq{
			ProductID: "productID2",
			Quantity:  3,
		},
	}

	suite.mockRepo.On("GetCartByUserID", mock.Anything, "userID").
		Return(
			&model.Cart{
				ID:     "cartId1",
				UserID: "userID",
				Lines: []*model.CartLine{
					{
						ProductID: "productID1",
						Quantity:  4,
					},
					{
						ProductID: "productID2",
						Quantity:  3,
					},
				},
			},
			nil,
		).Times(1)

	cart, err := suite.service.AddProduct(context.Background(), req)
	suite.NotNil(cart)
	suite.Equal("userID", cart.UserID)
	suite.Equal(2, len(cart.Lines))
	suite.Nil(err)
}

func (suite *CartServiceTestSuite) TestAddProductUpdateFail() {
	req := &dto.AddProductReq{
		UserID: "userID",
		Line: &dto.CartLineReq{
			ProductID: "productID2",
			Quantity:  3,
		},
	}

	suite.mockRepo.On("GetCartByUserID", mock.Anything, "userID").
		Return(
			&model.Cart{
				ID:     "cartId1",
				UserID: "userID",
				Lines: []*model.CartLine{
					{
						ProductID: "productID1",
						Quantity:  4,
					},
				},
			},
			nil,
		).Times(1)

	suite.mockRepo.On("Update", mock.Anything, &model.Cart{
		ID:     "cartId1",
		UserID: "userID",
		Lines: []*model.CartLine{
			{
				ProductID: "productID1",
				Quantity:  4,
			},
			{
				ProductID: "productID2",
				Quantity:  3,
			},
		},
	}).Return(errors.New("error")).Times(1)

	cart, err := suite.service.AddProduct(context.Background(), req)
	suite.Nil(cart)
	suite.NotNil(err)
}

// RemoveProduct
// =================================================================

func (suite *CartServiceTestSuite) TestRemoveProductSuccessfully() {
	req := &dto.RemoveProductReq{
		UserID:    "userID",
		ProductID: "productID1",
	}

	suite.mockRepo.On("GetCartByUserID", mock.Anything, "userID").
		Return(
			&model.Cart{
				ID:     "cartId1",
				UserID: "userID",
				Lines: []*model.CartLine{
					{
						ProductID: "productID1",
						Quantity:  4,
					},
				},
			},
			nil,
		).Times(1)

	suite.mockRepo.On("Update", mock.Anything, &model.Cart{
		ID:     "cartId1",
		UserID: "userID",
		Lines:  []*model.CartLine{},
	}).Return(nil).Times(1)

	cart, err := suite.service.RemoveProduct(context.Background(), req)
	suite.NotNil(cart)
	suite.Equal("userID", cart.UserID)
	suite.Equal(0, len(cart.Lines))
	suite.Nil(err)
}

func (suite *CartServiceTestSuite) TestRemoveProductMissUserID() {
	req := &dto.RemoveProductReq{
		ProductID: "productID1",
	}

	cart, err := suite.service.RemoveProduct(context.Background(), req)
	suite.Nil(cart)
	suite.NotNil(err)
}

func (suite *CartServiceTestSuite) TestRemoveProductMissProductID() {
	req := &dto.RemoveProductReq{
		UserID: "userID",
	}

	cart, err := suite.service.RemoveProduct(context.Background(), req)
	suite.Nil(cart)
	suite.NotNil(err)
}

func (suite *CartServiceTestSuite) TestRemoveProductCartNotFound() {
	req := &dto.RemoveProductReq{
		UserID:    "userID",
		ProductID: "productID1",
	}

	suite.mockRepo.On("GetCartByUserID", mock.Anything, "userID").
		Return(nil, errors.New("error")).Times(1)

	suite.mockRepo.On("Create", mock.Anything, &model.Cart{UserID: "userID"}).Return(nil).Times(1)

	cart, err := suite.service.RemoveProduct(context.Background(), req)
	suite.NotNil(cart)
	suite.Equal("userID", cart.UserID)
	suite.Equal(0, len(cart.Lines))
	suite.Nil(err)
}

func (suite *CartServiceTestSuite) TestRemoveProductCartNotFoundCreateFail() {
	req := &dto.RemoveProductReq{
		UserID:    "userID",
		ProductID: "productID1",
	}

	suite.mockRepo.On("GetCartByUserID", mock.Anything, "userID").
		Return(nil, errors.New("error")).Times(1)

	suite.mockRepo.On("Create", mock.Anything, &model.Cart{UserID: "userID"}).Return(errors.New("error")).Times(1)

	cart, err := suite.service.RemoveProduct(context.Background(), req)
	suite.Nil(cart)
	suite.NotNil(err)
}

func (suite *CartServiceTestSuite) TestRemoveProductUpdateFail() {
	req := &dto.RemoveProductReq{
		UserID:    "userID",
		ProductID: "productID1",
	}

	suite.mockRepo.On("GetCartByUserID", mock.Anything, "userID").
		Return(
			&model.Cart{
				ID:     "cartId1",
				UserID: "userID",
				Lines: []*model.CartLine{
					{
						ProductID: "productID1",
						Quantity:  4,
					},
				},
			},
			nil,
		).Times(1)

	suite.mockRepo.On("Update", mock.Anything, &model.Cart{
		ID:     "cartId1",
		UserID: "userID",
		Lines:  []*model.CartLine{},
	}).Return(errors.New("error")).Times(1)

	cart, err := suite.service.RemoveProduct(context.Background(), req)
	suite.Nil(cart)
	suite.NotNil(err)
}
