package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	type Src struct {
		Name  string
		Value int
	}
	type Dst struct {
		Name  string
		Value int
	}
	type MapDst struct {
		Name string `json:"name"`
	}

	tests := []struct {
		name     string
		src      interface{}
		newDst   func() interface{}
		assertFn func(t *testing.T, dst interface{})
	}{
		{
			name: "struct to struct",
			src:  &Src{Name: "test", Value: 42},
			newDst: func() interface{} {
				return &Dst{}
			},
			assertFn: func(t *testing.T, dst interface{}) {
				d := dst.(*Dst)
				assert.Equal(t, "test", d.Name)
				assert.Equal(t, 42, d.Value)
			},
		},
		{
			name: "nil source",
			src:  nil,
			newDst: func() interface{} {
				var dst map[string]interface{}
				return &dst
			},
			assertFn: func(t *testing.T, dst interface{}) {
				d := dst.(*map[string]interface{})
				assert.Nil(t, *d)
			},
		},
		{
			name: "map to struct",
			src:  map[string]interface{}{"name": "hello"},
			newDst: func() interface{} {
				return &MapDst{}
			},
			assertFn: func(t *testing.T, dst interface{}) {
				d := dst.(*MapDst)
				assert.Equal(t, "hello", d.Name)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dst := tc.newDst()
			Copy(dst, tc.src)
			tc.assertFn(t, dst)
		})
	}
}
