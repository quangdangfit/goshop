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

	"goshop/internal/order/domain"
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

func (suite *CouponServiceTestSuite) TestGetByCode() {
	tests := []struct {
		name    string
		code    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			code: "SAVE10",
			setup: func() {
				suite.mockRepo.On("GetByCode", mock.Anything, "SAVE10").
					Return(&model.Coupon{ID: "c1", Code: "SAVE10"}, nil).Times(1)
			},
		},
		{
			name: "Not found",
			code: "INVALID",
			setup: func() {
				suite.mockRepo.On("GetByCode", mock.Anything, "INVALID").
					Return(nil, errors.New("not found")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			coupon, err := suite.service.GetByCode(context.Background(), tc.code)
			if tc.wantErr {
				suite.NotNil(err)
				suite.Nil(coupon)
			} else {
				suite.Nil(err)
				suite.Equal(tc.code, coupon.Code)
			}
		})
	}
}

func (suite *CouponServiceTestSuite) TestCreate() {
	tests := []struct {
		name    string
		req     *domain.CreateCouponReq
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			req:  &domain.CreateCouponReq{Code: "SAVE10", DiscountType: "fixed", DiscountValue: 10},
			setup: func() {
				suite.mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name:    "Validation fail",
			req:     &domain.CreateCouponReq{},
			setup:   func() {},
			wantErr: true,
		},
		{
			name: "DB fail",
			req:  &domain.CreateCouponReq{Code: "SAVE10", DiscountType: "fixed", DiscountValue: 10},
			setup: func() {
				suite.mockRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("duplicate")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			coupon, err := suite.service.Create(context.Background(), tc.req)
			if tc.wantErr {
				suite.NotNil(err)
				suite.Nil(coupon)
			} else {
				suite.Nil(err)
				suite.NotNil(coupon)
				suite.Equal("SAVE10", coupon.Code)
			}
		})
	}
}

func (suite *CouponServiceTestSuite) TestApply() {
	past := time.Now().Add(-24 * time.Hour)
	tests := []struct {
		name         string
		code         string
		totalPrice   float64
		setup        func()
		wantErr      bool
		wantDiscount float64
	}{
		{
			name:       "Fixed discount",
			code:       "SAVE10",
			totalPrice: 100,
			setup: func() {
				suite.mockRepo.On("GetByCode", mock.Anything, "SAVE10").
					Return(&model.Coupon{ID: "c1", Code: "SAVE10", DiscountType: model.DiscountTypeFixed, DiscountValue: 10}, nil).Times(1)
			},
			wantDiscount: 10,
		},
		{
			name:       "Percentage discount",
			code:       "SAVE10PCT",
			totalPrice: 200,
			setup: func() {
				suite.mockRepo.On("GetByCode", mock.Anything, "SAVE10PCT").
					Return(&model.Coupon{ID: "c1", Code: "SAVE10PCT", DiscountType: model.DiscountTypePercentage, DiscountValue: 10}, nil).Times(1)
			},
			wantDiscount: 20,
		},
		{
			name:       "Fixed discount exceeds total",
			code:       "SAVE100",
			totalPrice: 50,
			setup: func() {
				suite.mockRepo.On("GetByCode", mock.Anything, "SAVE100").
					Return(&model.Coupon{ID: "c1", Code: "SAVE100", DiscountType: model.DiscountTypeFixed, DiscountValue: 200}, nil).Times(1)
			},
			wantDiscount: 50,
		},
		{
			name:       "Not found",
			code:       "INVALID",
			totalPrice: 100,
			setup: func() {
				suite.mockRepo.On("GetByCode", mock.Anything, "INVALID").
					Return(nil, errors.New("not found")).Times(1)
			},
			wantErr: true,
		},
		{
			name:       "Expired",
			code:       "EXPIRED",
			totalPrice: 100,
			setup: func() {
				suite.mockRepo.On("GetByCode", mock.Anything, "EXPIRED").
					Return(&model.Coupon{ID: "c1", Code: "EXPIRED", ExpiresAt: &past}, nil).Times(1)
			},
			wantErr: true,
		},
		{
			name:       "Max usage reached",
			code:       "MAXED",
			totalPrice: 100,
			setup: func() {
				suite.mockRepo.On("GetByCode", mock.Anything, "MAXED").
					Return(&model.Coupon{ID: "c1", Code: "MAXED", MaxUsage: 5, UsedCount: 5}, nil).Times(1)
			},
			wantErr: true,
		},
		{
			name:       "Below min order amount",
			code:       "MIN50",
			totalPrice: 30,
			setup: func() {
				suite.mockRepo.On("GetByCode", mock.Anything, "MIN50").
					Return(&model.Coupon{ID: "c1", Code: "MIN50", DiscountType: model.DiscountTypeFixed, DiscountValue: 10, MinOrderAmount: 50}, nil).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			discount, coupon, err := suite.service.Apply(context.Background(), tc.code, tc.totalPrice)
			if tc.wantErr {
				suite.NotNil(err)
				suite.Nil(coupon)
				suite.Equal(float64(0), discount)
			} else {
				suite.Nil(err)
				suite.NotNil(coupon)
				suite.Equal(tc.wantDiscount, discount)
			}
		})
	}
}

func (suite *CouponServiceTestSuite) TestIncrUsedCount() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockRepo.On("IncrUsedCount", mock.Anything, "c1").Return(nil).Times(1)
			},
		},
		{
			name: "DB error",
			setup: func() {
				suite.mockRepo.On("IncrUsedCount", mock.Anything, "c1").Return(errors.New("db error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			err := suite.service.IncrUsedCount(context.Background(), "c1")
			if tc.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}
