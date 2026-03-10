package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCategory_BeforeCreate(t *testing.T) {
	c := &Category{}
	err := c.BeforeCreate(nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, c.ID)
}
