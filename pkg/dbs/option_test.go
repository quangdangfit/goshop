package dbs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithOptions(t *testing.T) {
	tests := []struct {
		name     string
		opts     []FindOption
		assertFn func(t *testing.T, opt Option)
	}{
		{
			name: "WithQuery",
			opts: []FindOption{WithQuery(NewQuery("id = ?", "123"))},
			assertFn: func(t *testing.T, opt Option) {
				assert.Equal(t, []Query{NewQuery("id = ?", "123")}, opt.query)
			},
		},
		{
			name: "WithOffset",
			opts: []FindOption{WithOffset(10)},
			assertFn: func(t *testing.T, opt Option) {
				assert.Equal(t, 10, opt.offset)
			},
		},
		{
			name: "WithLimit",
			opts: []FindOption{WithLimit(50)},
			assertFn: func(t *testing.T, opt Option) {
				assert.Equal(t, 50, opt.limit)
			},
		},
		{
			name: "WithOrder",
			opts: []FindOption{WithOrder("created_at DESC")},
			assertFn: func(t *testing.T, opt Option) {
				assert.Equal(t, "created_at DESC", opt.order)
			},
		},
		{
			name: "WithPreload",
			opts: []FindOption{WithPreload([]string{"Lines", "Lines.Product"})},
			assertFn: func(t *testing.T, opt Option) {
				assert.Equal(t, []string{"Lines", "Lines.Product"}, opt.preloads)
			},
		},
		{
			name: "defaults when no options",
			opts: nil,
			assertFn: func(t *testing.T, opt Option) {
				assert.Equal(t, 0, opt.offset)
				assert.Equal(t, 1000, opt.limit)
				assert.Equal(t, "id", opt.order)
				assert.Empty(t, opt.preloads)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			opt := getOption(tc.opts...)
			tc.assertFn(t, opt)
		})
	}
}

func TestNewQuery(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		args     []any
		wantArgs []any
	}{
		{
			name:     "with args",
			query:    "name = ?",
			args:     []any{"test"},
			wantArgs: []any{"test"},
		},
		{
			name:     "no args",
			query:    "deleted_at IS NULL",
			args:     nil,
			wantArgs: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			q := NewQuery(tc.query, tc.args...)
			assert.Equal(t, tc.query, q.Query)
			assert.Equal(t, tc.wantArgs, q.Args)
		})
	}
}
