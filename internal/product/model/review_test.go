package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReview_BeforeCreate(t *testing.T) {
	r := &Review{}
	err := r.BeforeCreate(nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, r.ID)
}
