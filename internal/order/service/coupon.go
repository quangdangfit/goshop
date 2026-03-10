package service

import (
	"context"
	"errors"
	"time"

	"github.com/quangdangfit/gocommon/validation"

	"goshop/internal/order/dto"
	"goshop/internal/order/model"
	"goshop/internal/order/repository"
	"goshop/pkg/utils"
)

//go:generate mockery --name=CouponService
type CouponService interface {
	GetByCode(ctx context.Context, code string) (*model.Coupon, error)
	Create(ctx context.Context, req *dto.CreateCouponReq) (*model.Coupon, error)
	Apply(ctx context.Context, code string, totalPrice float64) (discountAmount float64, coupon *model.Coupon, err error)
	IncrUsedCount(ctx context.Context, id string) error
}

type couponSvc struct {
	validator validation.Validation
	repo      repository.CouponRepository
}

func NewCouponService(validator validation.Validation, repo repository.CouponRepository) CouponService {
	return &couponSvc{validator: validator, repo: repo}
}

func (s *couponSvc) GetByCode(ctx context.Context, code string) (*model.Coupon, error) {
	return s.repo.GetByCode(ctx, code)
}

func (s *couponSvc) Create(ctx context.Context, req *dto.CreateCouponReq) (*model.Coupon, error) {
	if err := s.validator.ValidateStruct(req); err != nil {
		return nil, err
	}
	var coupon model.Coupon
	utils.Copy(&coupon, req)
	if err := s.repo.Create(ctx, &coupon); err != nil {
		return nil, err
	}
	return &coupon, nil
}

func (s *couponSvc) Apply(ctx context.Context, code string, totalPrice float64) (float64, *model.Coupon, error) {
	coupon, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return 0, nil, errors.New("coupon not found")
	}

	now := time.Now()
	if coupon.ExpiresAt != nil && coupon.ExpiresAt.Before(now) {
		return 0, nil, errors.New("coupon has expired")
	}
	if coupon.MaxUsage > 0 && coupon.UsedCount >= coupon.MaxUsage {
		return 0, nil, errors.New("coupon has reached maximum usage")
	}
	if totalPrice < coupon.MinOrderAmount {
		return 0, nil, errors.New("order total is below minimum required for this coupon")
	}

	var discount float64
	switch coupon.DiscountType {
	case model.DiscountTypeFixed:
		discount = coupon.DiscountValue
		if discount > totalPrice {
			discount = totalPrice
		}
	case model.DiscountTypePercentage:
		discount = totalPrice * coupon.DiscountValue / 100
	}

	return discount, coupon, nil
}

func (s *couponSvc) IncrUsedCount(ctx context.Context, id string) error {
	return s.repo.IncrUsedCount(ctx, id)
}
