package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"goshop/internal/order/dto"
	"goshop/internal/order/model"
	orderMocks "goshop/internal/order/repository/mocks"
	"goshop/pkg/config"
)

type CouponServiceTestSuite struct {
	suite.Suite
	mockRepo *orderMocks.CouponRepository
	service  CouponService
}

func (suite *CouponServiceTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)
	validator := validation.New()
	suite.mockRepo = orderMocks.NewCouponRepository(suite.T())
	suite.service = NewCouponService(validator, suite.mockRepo)
}

func TestCouponServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CouponServiceTestSuite))
}

// GetByCode
// =================================================================================================

func (suite *CouponServiceTestSuite) TestGetByCodeSuccess() {
	suite.mockRepo.On("GetByCode", mock.Anything, "SAVE10").
		Return(&model.Coupon{ID: "c1", Code: "SAVE10"}, nil).Times(1)

	coupon, err := suite.service.GetByCode(context.Background(), "SAVE10")
	suite.Nil(err)
	suite.Equal("SAVE10", coupon.Code)
}

func (suite *CouponServiceTestSuite) TestGetByCodeFail() {
	suite.mockRepo.On("GetByCode", mock.Anything, "INVALID").
		Return(nil, errors.New("not found")).Times(1)

	coupon, err := suite.service.GetByCode(context.Background(), "INVALID")
	suite.NotNil(err)
	suite.Nil(coupon)
}

// Create
// =================================================================================================

func (suite *CouponServiceTestSuite) TestCreateSuccess() {
	req := &dto.CreateCouponReq{
		Code:          "SAVE10",
		DiscountType:  "fixed",
		DiscountValue: 10,
	}
	suite.mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Times(1)

	coupon, err := suite.service.Create(context.Background(), req)
	suite.Nil(err)
	suite.NotNil(coupon)
	suite.Equal("SAVE10", coupon.Code)
}

func (suite *CouponServiceTestSuite) TestCreateValidationFail() {
	req := &dto.CreateCouponReq{} // missing required fields

	coupon, err := suite.service.Create(context.Background(), req)
	suite.NotNil(err)
	suite.Nil(coupon)
}

func (suite *CouponServiceTestSuite) TestCreateDBFail() {
	req := &dto.CreateCouponReq{
		Code:          "SAVE10",
		DiscountType:  "fixed",
		DiscountValue: 10,
	}
	suite.mockRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("duplicate")).Times(1)

	coupon, err := suite.service.Create(context.Background(), req)
	suite.NotNil(err)
	suite.Nil(coupon)
}

// Apply
// =================================================================================================

func (suite *CouponServiceTestSuite) TestApplyFixedDiscountSuccess() {
	suite.mockRepo.On("GetByCode", mock.Anything, "SAVE10").
		Return(&model.Coupon{
			ID:            "c1",
			Code:          "SAVE10",
			DiscountType:  model.DiscountTypeFixed,
			DiscountValue: 10,
		}, nil).Times(1)

	discount, coupon, err := suite.service.Apply(context.Background(), "SAVE10", 100)
	suite.Nil(err)
	suite.NotNil(coupon)
	suite.Equal(float64(10), discount)
}

func (suite *CouponServiceTestSuite) TestApplyPercentageDiscountSuccess() {
	suite.mockRepo.On("GetByCode", mock.Anything, "SAVE10PCT").
		Return(&model.Coupon{
			ID:            "c1",
			Code:          "SAVE10PCT",
			DiscountType:  model.DiscountTypePercentage,
			DiscountValue: 10,
		}, nil).Times(1)

	discount, coupon, err := suite.service.Apply(context.Background(), "SAVE10PCT", 200)
	suite.Nil(err)
	suite.NotNil(coupon)
	suite.Equal(float64(20), discount)
}

func (suite *CouponServiceTestSuite) TestApplyFixedDiscountExceedsTotal() {
	suite.mockRepo.On("GetByCode", mock.Anything, "SAVE100").
		Return(&model.Coupon{
			ID:            "c1",
			Code:          "SAVE100",
			DiscountType:  model.DiscountTypeFixed,
			DiscountValue: 200,
		}, nil).Times(1)

	discount, coupon, err := suite.service.Apply(context.Background(), "SAVE100", 50)
	suite.Nil(err)
	suite.NotNil(coupon)
	suite.Equal(float64(50), discount) // capped at total price
}

func (suite *CouponServiceTestSuite) TestApplyCouponNotFound() {
	suite.mockRepo.On("GetByCode", mock.Anything, "INVALID").
		Return(nil, errors.New("not found")).Times(1)

	discount, coupon, err := suite.service.Apply(context.Background(), "INVALID", 100)
	suite.NotNil(err)
	suite.Nil(coupon)
	suite.Equal(float64(0), discount)
}

func (suite *CouponServiceTestSuite) TestApplyCouponExpired() {
	past := time.Now().Add(-24 * time.Hour)
	suite.mockRepo.On("GetByCode", mock.Anything, "EXPIRED").
		Return(&model.Coupon{
			ID:        "c1",
			Code:      "EXPIRED",
			ExpiresAt: &past,
		}, nil).Times(1)

	discount, coupon, err := suite.service.Apply(context.Background(), "EXPIRED", 100)
	suite.NotNil(err)
	suite.Nil(coupon)
	suite.Equal(float64(0), discount)
}

func (suite *CouponServiceTestSuite) TestApplyCouponMaxUsageReached() {
	suite.mockRepo.On("GetByCode", mock.Anything, "MAXED").
		Return(&model.Coupon{
			ID:        "c1",
			Code:      "MAXED",
			MaxUsage:  5,
			UsedCount: 5,
		}, nil).Times(1)

	discount, coupon, err := suite.service.Apply(context.Background(), "MAXED", 100)
	suite.NotNil(err)
	suite.Nil(coupon)
	suite.Equal(float64(0), discount)
}

func (suite *CouponServiceTestSuite) TestApplyBelowMinOrderAmount() {
	suite.mockRepo.On("GetByCode", mock.Anything, "MIN50").
		Return(&model.Coupon{
			ID:             "c1",
			Code:           "MIN50",
			DiscountType:   model.DiscountTypeFixed,
			DiscountValue:  10,
			MinOrderAmount: 50,
		}, nil).Times(1)

	discount, coupon, err := suite.service.Apply(context.Background(), "MIN50", 30)
	suite.NotNil(err)
	suite.Nil(coupon)
	suite.Equal(float64(0), discount)
}

// IncrUsedCount
// =================================================================================================

func (suite *CouponServiceTestSuite) TestIncrUsedCountSuccess() {
	suite.mockRepo.On("IncrUsedCount", mock.Anything, "c1").Return(nil).Times(1)

	err := suite.service.IncrUsedCount(context.Background(), "c1")
	suite.Nil(err)
}

func (suite *CouponServiceTestSuite) TestIncrUsedCountFail() {
	suite.mockRepo.On("IncrUsedCount", mock.Anything, "c1").Return(errors.New("db error")).Times(1)

	err := suite.service.IncrUsedCount(context.Background(), "c1")
	suite.NotNil(err)
}
