package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCart_BeforeCreate(t *testing.T) {
	cart := &Cart{}
	err := cart.BeforeCreate(nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, cart.ID)
}
