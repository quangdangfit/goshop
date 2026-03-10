package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoupon_BeforeCreate(t *testing.T) {
	c := &Coupon{}
	err := c.BeforeCreate(nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, c.ID)
}
