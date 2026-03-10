package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWishlist_BeforeCreate(t *testing.T) {
	w := &Wishlist{}
	err := w.BeforeCreate(nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, w.ID)
}
