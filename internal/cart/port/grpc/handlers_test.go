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
	mockService *mocks.ICartService
	handler     *CartHandler
}

func (suite *CartHandlerTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)

	suite.mockService = mocks.NewICartService(suite.T())
	suite.handler = NewCartHandler(suite.mockService)
}

func TestCartHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(CartHandlerTestSuite))
}

// AddProduct
// =================================================================================================

func (suite *CartHandlerTestSuite) TestCartAPI_AddProductSuccess() {
	req := &pb.AddProductReq{
		ProductId: "productId",
		Quantity:  2,
	}

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

	ctx := context.WithValue(context.Background(), "userId", "userID")
	res, err := suite.handler.AddProduct(ctx, req)

	suite.Nil(err)
	suite.Equal("userID", res.Cart.User.Id)
	suite.Equal(1, len(res.Cart.Lines))
}

func (suite *CartHandlerTestSuite) TestCartAPI_AddProductFail() {
	req := &pb.AddProductReq{
		ProductId: "productId",
		Quantity:  2,
	}

	suite.mockService.On("AddProduct", mock.Anything, &dto.AddProductReq{
		UserID: "userID",
		Line: &dto.CartLineReq{
			ProductID: "productId",
			Quantity:  2,
		},
	}).Return(nil, errors.New("error")).Times(1)

	ctx := context.WithValue(context.Background(), "userId", "userID")
	res, err := suite.handler.AddProduct(ctx, req)
	suite.Nil(res)
	suite.NotNil(err)
}

func (suite *CartHandlerTestSuite) TestCartAPI_AddProductUnauthorized() {
	req := &pb.AddProductReq{
		ProductId: "productId",
		Quantity:  2,
	}

	res, err := suite.handler.AddProduct(context.Background(), req)
	suite.Nil(res)
	suite.NotNil(err)
}

// RemoveProduct
// =================================================================================================

func (suite *CartHandlerTestSuite) TestCartAPI_RemoveProductSuccess() {
	req := &pb.RemoveProductReq{
		ProductId: "productId",
	}

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

	ctx := context.WithValue(context.Background(), "userId", "userID")
	res, err := suite.handler.RemoveProduct(ctx, req)

	suite.Nil(err)
	suite.Equal("userID", res.Cart.User.Id)
	suite.Equal(1, len(res.Cart.Lines))
}

func (suite *CartHandlerTestSuite) TestCartAPI_RemoveProductFail() {
	req := &pb.RemoveProductReq{
		ProductId: "productId",
	}

	suite.mockService.On("RemoveProduct", mock.Anything, &dto.RemoveProductReq{
		UserID:    "userID",
		ProductID: "productId",
	}).Return(nil, errors.New("error")).Times(1)

	ctx := context.WithValue(context.Background(), "userId", "userID")
	res, err := suite.handler.RemoveProduct(ctx, req)
	suite.Nil(res)
	suite.NotNil(err)
}

func (suite *CartHandlerTestSuite) TestCartAPI_RemoveProductUnauthorized() {
	req := &pb.RemoveProductReq{
		ProductId: "productId",
	}

	res, err := suite.handler.RemoveProduct(context.Background(), req)
	suite.Nil(res)
	suite.NotNil(err)
}

// GetCart
// =================================================================================================

func (suite *CartHandlerTestSuite) TestCartAPI_GetCartSuccess() {
	req := &pb.GetCartReq{}

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

	ctx := context.WithValue(context.Background(), "userId", "userID")
	res, err := suite.handler.GetCart(ctx, req)

	suite.Nil(err)
	suite.Equal("userID", res.Cart.User.Id)
	suite.Equal(1, len(res.Cart.Lines))
}

func (suite *CartHandlerTestSuite) TestCartAPI_GetCartFail() {
	req := &pb.GetCartReq{}

	suite.mockService.On("GetCartByUserID", mock.Anything, "userID").Return(nil, errors.New("error")).Times(1)

	ctx := context.WithValue(context.Background(), "userId", "userID")
	res, err := suite.handler.GetCart(ctx, req)
	suite.Nil(res)
	suite.NotNil(err)
}

func (suite *CartHandlerTestSuite) TestCartAPI_GetCartUnauthorized() {
	req := &pb.GetCartReq{}

	res, err := suite.handler.GetCart(context.Background(), req)
	suite.Nil(res)
	suite.NotNil(err)
}
