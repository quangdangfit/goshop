package service

import (
	"context"
	"time"

	"github.com/quangdangfit/gocommon/validation"

	"goshop/internal/order/domain"
	"goshop/internal/order/model"
	"goshop/internal/order/repository"
	"goshop/pkg/apperror"
	"goshop/pkg/utils"
)

//go:generate mockery --name=CouponService
type CouponService interface {
	GetByCode(ctx context.Context, code string) (*model.Coupon, error)
	Create(ctx context.Context, req *domain.CreateCouponReq) (*model.Coupon, error)
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

func (s *couponSvc) Create(ctx context.Context, req *domain.CreateCouponReq) (*model.Coupon, error) {
	if err := s.validator.ValidateStruct(req); err != nil {
		return nil, err
	}
	var coupon model.Coupon
	if err := utils.Copy(&coupon, req); err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, &coupon); err != nil {
		return nil, err
	}
	return &coupon, nil
}

func (s *couponSvc) Apply(ctx context.Context, code string, totalPrice float64) (float64, *model.Coupon, error) {
	coupon, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return 0, nil, apperror.Wrap(apperror.ErrNotFound, err)
	}

	now := time.Now()
	if coupon.ExpiresAt != nil && coupon.ExpiresAt.Before(now) {
		return 0, nil, apperror.ErrCouponExpired
	}
	if coupon.MaxUsage > 0 && coupon.UsedCount >= coupon.MaxUsage {
		return 0, nil, apperror.ErrCouponMaxUsage
	}
	if totalPrice < coupon.MinOrderAmount {
		return 0, nil, apperror.ErrCouponMinOrder
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
