package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrder_BeforeCreate_DefaultStatus(t *testing.T) {
	order := &Order{}
	err := order.BeforeCreate(nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, order.ID)
	assert.NotEmpty(t, order.Code)
	assert.Equal(t, OrderStatusNew, order.Status)
}

func TestOrder_BeforeCreate_ExistingStatus(t *testing.T) {
	order := &Order{Status: OrderStatusInProgress}
	err := order.BeforeCreate(nil)
	assert.NoError(t, err)
	assert.Equal(t, OrderStatusInProgress, order.Status)
}
