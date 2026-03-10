package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddress_BeforeCreate(t *testing.T) {
	a := &Address{}
	err := a.BeforeCreate(nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, a.ID)
}
