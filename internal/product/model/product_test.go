package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProduct_BeforeCreate(t *testing.T) {
	p := &Product{}
	err := p.BeforeCreate(nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, p.ID)
	assert.NotEmpty(t, p.Code)
	assert.True(t, p.Active)
}
