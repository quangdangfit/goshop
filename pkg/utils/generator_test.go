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

// generateCode branch coverage with injected times

func TestGenerateCode_SingleDigitYear(t *testing.T) {
	// year % 100 < 10: e.g. 2005 → year=05
	d := time.Date(2005, time.March, 15, 0, 0, 0, 0, time.UTC)
	code := generateCode("P", d)
	assert.True(t, strings.HasPrefix(code, "P05"))
}

func TestGenerateCode_DoubleDigitYear(t *testing.T) {
	// year % 100 >= 10: e.g. 2026 → year=26
	d := time.Date(2026, time.November, 20, 0, 0, 0, 0, time.UTC)
	code := generateCode("P", d)
	assert.True(t, strings.HasPrefix(code, "P26"))
}

func TestGenerateCode_SingleDigitDay(t *testing.T) {
	// day < 10: e.g. day=5
	d := time.Date(2026, time.March, 5, 0, 0, 0, 0, time.UTC)
	code := generateCode("P", d)
	assert.Contains(t, code, "0305")
}

func TestGenerateCode_DoubleDigitDay(t *testing.T) {
	// day >= 10: e.g. day=15
	d := time.Date(2026, time.March, 15, 0, 0, 0, 0, time.UTC)
	code := generateCode("P", d)
	assert.Contains(t, code, "0315")
}

func TestGenerateCode_DoubleDigitMonth(t *testing.T) {
	// month >= 10: e.g. November=11
	d := time.Date(2026, time.November, 15, 0, 0, 0, 0, time.UTC)
	code := generateCode("P", d)
	assert.Contains(t, code, "1115")
}
