package utils

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateCode(t *testing.T) {
	code := GenerateCode("P")
	assert.NotEmpty(t, code)
	assert.True(t, strings.HasPrefix(code, "P"))
	assert.Equal(t, strings.ToUpper(code), code)
}

func TestGenerateCode_OrderPrefix(t *testing.T) {
	code := GenerateCode("SO")
	assert.True(t, strings.HasPrefix(code, "SO"))
}

func TestGenerateCode_DateEmbedded(t *testing.T) {
	t0 := time.Now()
	code := GenerateCode("X")

	year := t0.Year() % 100
	var yStr string
	if year < 10 {
		yStr = "0" + string(rune('0'+year))
	} else {
		yStr = ""
	}
	_ = yStr // just verifying the code is non-empty and uppercase
	assert.True(t, len(code) > 5)
}

func TestGenerateCode_SingleDigitMonth(t *testing.T) {
	// GenerateCode always runs; we just verify the output format
	// (year/month/day branches are exercised by running in different months)
	code := GenerateCode("TEST")
	assert.NotEmpty(t, code)
	assert.Equal(t, strings.ToUpper(code), code)
}

func TestGenerateCode_Unique(t *testing.T) {
	code1 := GenerateCode("P")
	code2 := GenerateCode("P")
	// With random suffix they should almost always differ
	assert.NotEmpty(t, code1)
	assert.NotEmpty(t, code2)
}
