package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"goshop/internal/order/model"
)

func TestCouponFromModel(t *testing.T) {
	assert.Nil(t, CouponFromModel(nil))
	exp := time.Now()
	c := CouponFromModel(&model.Coupon{
		ID:             "c1",
		Code:           "X",
		DiscountType:   model.DiscountTypeFixed,
		DiscountValue:  10,
		MinOrderAmount: 5,
		MaxUsage:       2,
		UsedCount:      1,
		ExpiresAt:      &exp,
	})
	assert.Equal(t, "c1", c.ID)
	assert.Equal(t, "fixed", c.DiscountType)
	assert.Equal(t, &exp, c.ExpiresAt)
}

func TestOrderFromModel(t *testing.T) {
	assert.Nil(t, OrderFromModel(nil))
	o := OrderFromModel(&model.Order{
		ID:         "o1",
		Code:       "C",
		TotalPrice: 100,
		Status:     model.OrderStatusNew,
		Lines: []*model.OrderLine{
			{Quantity: 2, Price: 5, Product: &model.Product{ID: "p1", Code: "PC", Name: "n", Price: 5}},
			nil,
		},
	})
	assert.Equal(t, "o1", o.ID)
	assert.Equal(t, "new", o.Status)
	assert.Len(t, o.Lines, 2)
	assert.Equal(t, "p1", o.Lines[0].Product.ID)
	assert.Nil(t, o.Lines[1])
}

func TestOrdersFromModel(t *testing.T) {
	out := OrdersFromModel([]*model.Order{{ID: "a"}, nil})
	assert.Len(t, out, 2)
	assert.Equal(t, "a", out[0].ID)
	assert.Nil(t, out[1])
}

func TestOrderLineFromModel_NoProduct(t *testing.T) {
	assert.Nil(t, OrderLineFromModel(nil))
	out := OrderLineFromModel(&model.OrderLine{Quantity: 3, Price: 9})
	assert.Equal(t, uint(3), out.Quantity)
	assert.Equal(t, "", out.Product.ID)
}
