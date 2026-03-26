package utils

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateCode(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
		check  func(t *testing.T, code string)
	}{
		{
			name:   "product prefix",
			prefix: "P",
			check: func(t *testing.T, code string) {
				assert.True(t, strings.HasPrefix(code, "P"))
				assert.Equal(t, strings.ToUpper(code), code)
			},
		},
		{
			name:   "order prefix",
			prefix: "SO",
			check: func(t *testing.T, code string) {
				assert.True(t, strings.HasPrefix(code, "SO"))
			},
		},
		{
			name:   "long prefix",
			prefix: "TEST",
			check: func(t *testing.T, code string) {
				assert.NotEmpty(t, code)
				assert.Equal(t, strings.ToUpper(code), code)
			},
		},
		{
			name:   "code is long enough",
			prefix: "X",
			check: func(t *testing.T, code string) {
				assert.True(t, len(code) > 5)
			},
		},
		{
			name:   "two codes are non-empty",
			prefix: "P",
			check: func(t *testing.T, code string) {
				code2 := GenerateCode("P")
				assert.NotEmpty(t, code)
				assert.NotEmpty(t, code2)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			code := GenerateCode(tc.prefix)
			tc.check(t, code)
		})
	}
}

func TestGenerateCode_WithInjectedTime(t *testing.T) {
	tests := []struct {
		name     string
		prefix   string
		time     time.Time
		contains string
		prefix2  string // expected prefix in output
	}{
		{
			name:    "single digit year",
			prefix:  "P",
			time:    time.Date(2005, time.March, 15, 0, 0, 0, 0, time.UTC),
			prefix2: "P05",
		},
		{
			name:    "double digit year",
			prefix:  "P",
			time:    time.Date(2026, time.November, 20, 0, 0, 0, 0, time.UTC),
			prefix2: "P26",
		},
		{
			name:     "single digit day",
			prefix:   "P",
			time:     time.Date(2026, time.March, 5, 0, 0, 0, 0, time.UTC),
			contains: "0305",
		},
		{
			name:     "double digit day",
			prefix:   "P",
			time:     time.Date(2026, time.March, 15, 0, 0, 0, 0, time.UTC),
			contains: "0315",
		},
		{
			name:     "double digit month",
			prefix:   "P",
			time:     time.Date(2026, time.November, 15, 0, 0, 0, 0, time.UTC),
			contains: "1115",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			code := generateCode(tc.prefix, tc.time)
			if tc.prefix2 != "" {
				assert.True(t, strings.HasPrefix(code, tc.prefix2))
			}
			if tc.contains != "" {
				assert.Contains(t, code, tc.contains)
			}
		})
	}
}
