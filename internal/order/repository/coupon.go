package repository

import (
	"context"

	"gorm.io/gorm"

	"goshop/internal/order/model"
	"goshop/pkg/dbs"
)

//go:generate mockery --name=CouponRepository
type CouponRepository interface {
	GetByCode(ctx context.Context, code string) (*model.Coupon, error)
	Create(ctx context.Context, coupon *model.Coupon) error
	IncrUsedCount(ctx context.Context, id string) error
}

type couponRepo struct {
	db dbs.Database
}

func NewCouponRepository(db dbs.Database) CouponRepository {
	return &couponRepo{db: db}
}

func (r *couponRepo) GetByCode(ctx context.Context, code string) (*model.Coupon, error) {
	var coupon model.Coupon
	if err := r.db.FindOne(ctx, &coupon, dbs.WithQuery(dbs.NewQuery("code = ?", code))); err != nil {
		return nil, err
	}
	return &coupon, nil
}

func (r *couponRepo) Create(ctx context.Context, coupon *model.Coupon) error {
	return r.db.Create(ctx, coupon)
}

func (r *couponRepo) IncrUsedCount(ctx context.Context, id string) error {
	return r.db.GetDB().WithContext(ctx).
		Model(&model.Coupon{}).
		Where("id = ?", id).
		UpdateColumn("used_count", gorm.Expr("used_count + 1")).Error
}
