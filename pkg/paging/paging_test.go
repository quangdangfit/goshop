package paging

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew_DefaultPageSize(t *testing.T) {
	p := New(1, 0, 100)
	assert.Equal(t, int64(20), p.Limit)
	assert.Equal(t, int64(1), p.CurrentPage)
	assert.Equal(t, int64(100), p.Total)
	assert.Equal(t, int64(5), p.TotalPage)
	assert.Equal(t, int64(0), p.Skip)
}

func TestNew_CustomPageSize(t *testing.T) {
	p := New(2, 10, 50)
	assert.Equal(t, int64(10), p.Limit)
	assert.Equal(t, int64(2), p.CurrentPage)
	assert.Equal(t, int64(5), p.TotalPage)
	assert.Equal(t, int64(10), p.Skip)
}

func TestNew_PageSizeExceedsMax(t *testing.T) {
	p := New(1, 100, 200)
	assert.Equal(t, int64(DefaultPageSize), p.Limit)
}

func TestNew_PageZeroResetsToOne(t *testing.T) {
	p := New(0, 10, 50)
	assert.Equal(t, int64(1), p.CurrentPage)
}

func TestNew_PageNegativeResetsToOne(t *testing.T) {
	p := New(-5, 10, 50)
	assert.Equal(t, int64(1), p.CurrentPage)
}

func TestNew_ZeroTotal(t *testing.T) {
	p := New(1, 10, 0)
	assert.Equal(t, int64(0), p.TotalPage)
	assert.Equal(t, int64(1), p.CurrentPage)
}

func TestNew_LastPage(t *testing.T) {
	p := New(3, 10, 25)
	assert.Equal(t, int64(3), p.TotalPage)
	assert.Equal(t, int64(3), p.CurrentPage)
	assert.Equal(t, int64(20), p.Skip)
}
