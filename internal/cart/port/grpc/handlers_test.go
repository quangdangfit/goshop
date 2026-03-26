package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"goshop/internal/cart/dto"
	"goshop/internal/cart/model"
	"goshop/internal/cart/service/mocks"
	"goshop/pkg/config"
	pb "goshop/proto/gen/go/cart"
)

type CartHandlerTestSuite struct {
	suite.Suite
	mockService *mocks.CartService
	handler     *CartHandler
}

func (suite *CartHandlerTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)

	suite.mockService = mocks.NewCartService(suite.T())
	suite.handler = NewCartHandler(suite.mockService)
}

func TestCartHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(CartHandlerTestSuite))
}

// AddProduct
// =================================================================================================

func (suite *CartHandlerTestSuite) TestAddProduct() {
	tests := []struct {
		name      string
		setup     func()
		req       *pb.AddProductReq
		ctx       context.Context
		expectNil bool
		expectErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockService.On("AddProduct", mock.Anything, &dto.AddProductReq{
					UserID: "userID",
					Line: &dto.CartLineReq{
						ProductID: "productId",
						Quantity:  2,
					},
				}).Return(
					&model.Cart{
						UserID: "userID",
						User: &model.User{
							ID: "userID",
						},
						Lines: []*model.CartLine{
							{
								ProductID: "productId",
								Quantity:  2,
							},
						},
					},
					nil,
				).Times(1)
			},
			req: &pb.AddProductReq{
				ProductId: "productId",
				Quantity:  2,
			},
			ctx:       context.WithValue(context.Background(), "userId", "userID"),
			expectNil: false,
			expectErr: false,
		},
		{
			name: "Fail",
			setup: func() {
				suite.mockService.On("AddProduct", mock.Anything, &dto.AddProductReq{
					UserID: "userID",
					Line: &dto.CartLineReq{
						ProductID: "productId",
						Quantity:  2,
					},
				}).Return(nil, errors.New("error")).Times(1)
			},
			req: &pb.AddProductReq{
				ProductId: "productId",
				Quantity:  2,
			},
			ctx:       context.WithValue(context.Background(), "userId", "userID"),
			expectNil: true,
			expectErr: true,
		},
		{
			name:  "Unauthorized",
			setup: func() {},
			req: &pb.AddProductReq{
				ProductId: "productId",
				Quantity:  2,
			},
			ctx:       context.Background(),
			expectNil: true,
			expectErr: true,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()

			res, err := suite.handler.AddProduct(tc.ctx, tc.req)

			if tc.expectErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
			if tc.expectNil {
				suite.Nil(res)
			} else {
				suite.NotNil(res)
				suite.Equal("userID", res.Cart.User.Id)
				suite.Equal(1, len(res.Cart.Lines))
			}
		})
	}
}

// RemoveProduct
// =================================================================================================

func (suite *CartHandlerTestSuite) TestRemoveProduct() {
	tests := []struct {
		name      string
		setup     func()
		req       *pb.RemoveProductReq
		ctx       context.Context
		expectNil bool
		expectErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockService.On("RemoveProduct", mock.Anything, &dto.RemoveProductReq{
					UserID:    "userID",
					ProductID: "productId",
				}).Return(
					&model.Cart{
						UserID: "userID",
						User: &model.User{
							ID: "userID",
						},
						Lines: []*model.CartLine{
							{
								ProductID: "productId1",
								Quantity:  2,
							},
						},
					},
					nil,
				).Times(1)
			},
			req:       &pb.RemoveProductReq{ProductId: "productId"},
			ctx:       context.WithValue(context.Background(), "userId", "userID"),
			expectNil: false,
			expectErr: false,
		},
		{
			name: "Fail",
			setup: func() {
				suite.mockService.On("RemoveProduct", mock.Anything, &dto.RemoveProductReq{
					UserID:    "userID",
					ProductID: "productId",
				}).Return(nil, errors.New("error")).Times(1)
			},
			req:       &pb.RemoveProductReq{ProductId: "productId"},
			ctx:       context.WithValue(context.Background(), "userId", "userID"),
			expectNil: true,
			expectErr: true,
		},
		{
			name:      "Unauthorized",
			setup:     func() {},
			req:       &pb.RemoveProductReq{ProductId: "productId"},
			ctx:       context.Background(),
			expectNil: true,
			expectErr: true,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()

			res, err := suite.handler.RemoveProduct(tc.ctx, tc.req)

			if tc.expectErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
			if tc.expectNil {
				suite.Nil(res)
			} else {
				suite.NotNil(res)
				suite.Equal("userID", res.Cart.User.Id)
				suite.Equal(1, len(res.Cart.Lines))
			}
		})
	}
}

// GetCart
// =================================================================================================

func (suite *CartHandlerTestSuite) TestGetCart() {
	tests := []struct {
		name      string
		setup     func()
		ctx       context.Context
		expectNil bool
		expectErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockService.On("GetCartByUserID", mock.Anything, "userID").Return(
					&model.Cart{
						UserID: "userID",
						User: &model.User{
							ID: "userID",
						},
						Lines: []*model.CartLine{
							{
								ProductID: "productId1",
								Quantity:  2,
							},
						},
					},
					nil,
				).Times(1)
			},
			ctx:       context.WithValue(context.Background(), "userId", "userID"),
			expectNil: false,
			expectErr: false,
		},
		{
			name: "Fail",
			setup: func() {
				suite.mockService.On("GetCartByUserID", mock.Anything, "userID").Return(nil, errors.New("error")).Times(1)
			},
			ctx:       context.WithValue(context.Background(), "userId", "userID"),
			expectNil: true,
			expectErr: true,
		},
		{
			name:      "Unauthorized",
			setup:     func() {},
			ctx:       context.Background(),
			expectNil: true,
			expectErr: true,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()

			res, err := suite.handler.GetCart(tc.ctx, &pb.GetCartReq{})

			if tc.expectErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
			if tc.expectNil {
				suite.Nil(res)
			} else {
				suite.NotNil(res)
				suite.Equal("userID", res.Cart.User.Id)
				suite.Equal(1, len(res.Cart.Lines))
			}
		})
	}
}
