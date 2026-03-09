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

	src := Src{Name: "test", Value: 42}
	var dst Dst
	Copy(&dst, &src)

	assert.Equal(t, "test", dst.Name)
	assert.Equal(t, 42, dst.Value)
}

func TestCopy_NilSrc(t *testing.T) {
	var dst map[string]interface{}
	Copy(&dst, nil)
	assert.Nil(t, dst)
}

func TestCopy_MapToStruct(t *testing.T) {
	src := map[string]interface{}{
		"name": "hello",
	}
	type Dst struct {
		Name string `json:"name"`
	}
	var dst Dst
	Copy(&dst, src)
	assert.Equal(t, "hello", dst.Name)
}
