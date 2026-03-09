package dbs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithQuery(t *testing.T) {
	q := NewQuery("id = ?", "123")
	opt := getOption(WithQuery(q))
	assert.Equal(t, []Query{q}, opt.query)
}

func TestWithOffset(t *testing.T) {
	opt := getOption(WithOffset(10))
	assert.Equal(t, 10, opt.offset)
}

func TestWithLimit(t *testing.T) {
	opt := getOption(WithLimit(50))
	assert.Equal(t, 50, opt.limit)
}

func TestWithOrder(t *testing.T) {
	opt := getOption(WithOrder("created_at DESC"))
	assert.Equal(t, "created_at DESC", opt.order)
}

func TestWithPreload(t *testing.T) {
	preloads := []string{"Lines", "Lines.Product"}
	opt := getOption(WithPreload(preloads))
	assert.Equal(t, preloads, opt.preloads)
}

func TestGetOption_Defaults(t *testing.T) {
	opt := getOption()
	assert.Equal(t, 0, opt.offset)
	assert.Equal(t, 1000, opt.limit)
	assert.Equal(t, "id", opt.order)
	assert.Empty(t, opt.preloads)
}

func TestNewQuery(t *testing.T) {
	q := NewQuery("name = ?", "test")
	assert.Equal(t, "name = ?", q.Query)
	assert.Equal(t, []any{"test"}, q.Args)
}

func TestNewQuery_NoArgs(t *testing.T) {
	q := NewQuery("deleted_at IS NULL")
	assert.Equal(t, "deleted_at IS NULL", q.Query)
	assert.Nil(t, q.Args)
}
