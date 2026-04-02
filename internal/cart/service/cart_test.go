package service

import (
	"context"
	"errors"
	"testing"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"goshop/internal/cart/domain"
	"goshop/internal/cart/model"
	"goshop/internal/cart/repository/mocks"
	"goshop/pkg/config"
)

type CartServiceTestSuite struct {
	suite.Suite
	mockRepo *mocks.CartRepository
	service  CartService
}

func (suite *CartServiceTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)

	validator := validation.New()
	suite.mockRepo = mocks.NewCartRepository(suite.T())
	suite.service = NewCartService(validator, suite.mockRepo)
}

func TestCartServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CartServiceTestSuite))
}

func (suite *CartServiceTestSuite) TestGetCartByUserID() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
		wantLen int
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockRepo.On("GetCartByUserID", mock.Anything, "userID").
					Return(&model.Cart{
						ID: "cartId1", UserID: "userID",
						Lines: []*model.CartLine{
							{ProductID: "productID1", Quantity: 4},
							{ProductID: "productID2", Quantity: 3},
						},
					}, nil).Times(1)
			},
			wantLen: 2,
		},
		{
			name: "Not found, creates new",
			setup: func() {
				suite.mockRepo.On("GetCartByUserID", mock.Anything, "userID").
					Return(nil, errors.New("error")).Times(1)
				suite.mockRepo.On("Create", mock.Anything, &model.Cart{UserID: "userID"}).Return(nil).Times(1)
			},
			wantLen: 0,
		},
		{
			name: "Not found, create fails",
			setup: func() {
				suite.mockRepo.On("GetCartByUserID", mock.Anything, "userID").
					Return(nil, errors.New("error")).Times(1)
				suite.mockRepo.On("Create", mock.Anything, &model.Cart{UserID: "userID"}).Return(errors.New("error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			cart, err := suite.service.GetCartByUserID(context.Background(), "userID")
			if tc.wantErr {
				suite.Nil(cart)
				suite.NotNil(err)
			} else {
				suite.NotNil(cart)
				suite.Equal("userID", cart.UserID)
				suite.Equal(tc.wantLen, len(cart.Lines))
				suite.Nil(err)
			}
		})
	}
}

func (suite *CartServiceTestSuite) TestAddProduct() {
	tests := []struct {
		name    string
		req     *domain.AddProductReq
		setup   func()
		wantErr bool
		wantLen int
		wantQty uint
	}{
		{
			name: "Success - new product",
			req: &domain.AddProductReq{
				UserID: "userID",
				Line:   &domain.CartLineReq{ProductID: "productID2", Quantity: 3},
			},
			setup: func() {
				suite.mockRepo.On("GetCartByUserID", mock.Anything, "userID").
					Return(&model.Cart{
						ID: "cartId1", UserID: "userID",
						Lines: []*model.CartLine{{ProductID: "productID1", Quantity: 4}},
					}, nil).Times(1)
				suite.mockRepo.On("Update", mock.Anything, &model.Cart{
					ID: "cartId1", UserID: "userID",
					Lines: []*model.CartLine{
						{ProductID: "productID1", Quantity: 4},
						{ProductID: "productID2", Quantity: 3},
					},
				}).Return(nil).Times(1)
			},
			wantLen: 2,
		},
		{
			name: "Missing UserID",
			req: &domain.AddProductReq{
				Line: &domain.CartLineReq{ProductID: "productID2", Quantity: 3},
			},
			setup:   func() {},
			wantErr: true,
		},
		{
			name: "Missing ProductID",
			req: &domain.AddProductReq{
				UserID: "userID",
				Line:   &domain.CartLineReq{Quantity: 3},
			},
			setup:   func() {},
			wantErr: true,
		},
		{
			name: "Missing Quantity",
			req: &domain.AddProductReq{
				UserID: "userID",
				Line:   &domain.CartLineReq{ProductID: "productID2"},
			},
			setup:   func() {},
			wantErr: true,
		},
		{
			name: "Cart not found, creates new",
			req: &domain.AddProductReq{
				UserID: "userID",
				Line:   &domain.CartLineReq{ProductID: "productID2", Quantity: 3},
			},
			setup: func() {
				suite.mockRepo.On("GetCartByUserID", mock.Anything, "userID").
					Return(nil, errors.New("error")).Times(1)
				suite.mockRepo.On("Create", mock.Anything, &model.Cart{
					UserID: "userID",
					Lines:  []*model.CartLine{{ProductID: "productID2", Quantity: 3}},
				}).Return(nil).Times(1)
			},
			wantLen: 1,
		},
		{
			name: "Cart not found, create fails",
			req: &domain.AddProductReq{
				UserID: "userID",
				Line:   &domain.CartLineReq{ProductID: "productID2", Quantity: 3},
			},
			setup: func() {
				suite.mockRepo.On("GetCartByUserID", mock.Anything, "userID").
					Return(nil, errors.New("error")).Times(1)
				suite.mockRepo.On("Create", mock.Anything, &model.Cart{
					UserID: "userID",
					Lines:  []*model.CartLine{{ProductID: "productID2", Quantity: 3}},
				}).Return(errors.New("error")).Times(1)
			},
			wantErr: true,
		},
		{
			name: "Product already in cart - updates quantity",
			req: &domain.AddProductReq{
				UserID: "userID",
				Line:   &domain.CartLineReq{ProductID: "productID2", Quantity: 5},
			},
			setup: func() {
				suite.mockRepo.On("GetCartByUserID", mock.Anything, "userID").
					Return(&model.Cart{
						ID: "cartId1", UserID: "userID",
						Lines: []*model.CartLine{
							{ProductID: "productID1", Quantity: 4},
							{ProductID: "productID2", Quantity: 3},
						},
					}, nil).Times(1)
				suite.mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Times(1)
			},
			wantLen: 2,
			wantQty: 5,
		},
		{
			name: "Product already in cart - update fails",
			req: &domain.AddProductReq{
				UserID: "userID",
				Line:   &domain.CartLineReq{ProductID: "productID2", Quantity: 5},
			},
			setup: func() {
				suite.mockRepo.On("GetCartByUserID", mock.Anything, "userID").
					Return(&model.Cart{
						ID: "cartId1", UserID: "userID",
						Lines: []*model.CartLine{{ProductID: "productID2", Quantity: 3}},
					}, nil).Times(1)
				suite.mockRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("error")).Times(1)
			},
			wantErr: true,
		},
		{
			name: "New product - update fails",
			req: &domain.AddProductReq{
				UserID: "userID",
				Line:   &domain.CartLineReq{ProductID: "productID2", Quantity: 3},
			},
			setup: func() {
				suite.mockRepo.On("GetCartByUserID", mock.Anything, "userID").
					Return(&model.Cart{
						ID: "cartId1", UserID: "userID",
						Lines: []*model.CartLine{{ProductID: "productID1", Quantity: 4}},
					}, nil).Times(1)
				suite.mockRepo.On("Update", mock.Anything, &model.Cart{
					ID: "cartId1", UserID: "userID",
					Lines: []*model.CartLine{
						{ProductID: "productID1", Quantity: 4},
						{ProductID: "productID2", Quantity: 3},
					},
				}).Return(errors.New("error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			cart, err := suite.service.AddProduct(context.Background(), tc.req)
			if tc.wantErr {
				suite.Nil(cart)
				suite.NotNil(err)
			} else {
				suite.NotNil(cart)
				suite.Equal("userID", cart.UserID)
				suite.Equal(tc.wantLen, len(cart.Lines))
				if tc.wantQty > 0 {
					suite.Equal(tc.wantQty, cart.Lines[1].Quantity)
				}
				suite.Nil(err)
			}
		})
	}
}

func (suite *CartServiceTestSuite) TestRemoveProduct() {
	tests := []struct {
		name    string
		req     *domain.RemoveProductReq
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			req:  &domain.RemoveProductReq{UserID: "userID", ProductID: "productID1"},
			setup: func() {
				suite.mockRepo.On("GetCartByUserID", mock.Anything, "userID").
					Return(&model.Cart{
						ID: "cartId1", UserID: "userID",
						Lines: []*model.CartLine{{ProductID: "productID1", Quantity: 4}},
					}, nil).Times(1)
				suite.mockRepo.On("DeleteLine", mock.Anything, "cartId1", "productID1").Return(nil).Times(1)
			},
		},
		{
			name:    "Missing UserID",
			req:     &domain.RemoveProductReq{ProductID: "productID1"},
			setup:   func() {},
			wantErr: true,
		},
		{
			name:    "Missing ProductID",
			req:     &domain.RemoveProductReq{UserID: "userID"},
			setup:   func() {},
			wantErr: true,
		},
		{
			name: "Cart not found, creates new",
			req:  &domain.RemoveProductReq{UserID: "userID", ProductID: "productID1"},
			setup: func() {
				suite.mockRepo.On("GetCartByUserID", mock.Anything, "userID").
					Return(nil, errors.New("error")).Times(1)
				suite.mockRepo.On("Create", mock.Anything, &model.Cart{UserID: "userID"}).Return(nil).Times(1)
			},
		},
		{
			name: "Cart not found, create fails",
			req:  &domain.RemoveProductReq{UserID: "userID", ProductID: "productID1"},
			setup: func() {
				suite.mockRepo.On("GetCartByUserID", mock.Anything, "userID").
					Return(nil, errors.New("error")).Times(1)
				suite.mockRepo.On("Create", mock.Anything, &model.Cart{UserID: "userID"}).Return(errors.New("error")).Times(1)
			},
			wantErr: true,
		},
		{
			name: "DeleteLine fails",
			req:  &domain.RemoveProductReq{UserID: "userID", ProductID: "productID1"},
			setup: func() {
				suite.mockRepo.On("GetCartByUserID", mock.Anything, "userID").
					Return(&model.Cart{
						ID: "cartId1", UserID: "userID",
						Lines: []*model.CartLine{{ProductID: "productID1", Quantity: 4}},
					}, nil).Times(1)
				suite.mockRepo.On("DeleteLine", mock.Anything, "cartId1", "productID1").Return(errors.New("error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			cart, err := suite.service.RemoveProduct(context.Background(), tc.req)
			if tc.wantErr {
				suite.Nil(cart)
				suite.NotNil(err)
			} else {
				suite.NotNil(cart)
				suite.Equal("userID", cart.UserID)
				suite.Nil(err)
			}
		})
	}
}
