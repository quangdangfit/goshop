package paging

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		page      int64
		limit     int64
		total     int64
		wantLimit int64
		wantPage  int64
		wantTotal int64
		wantPages int64
		wantSkip  int64
	}{
		{
			name:      "default page size when limit is 0",
			page:      1,
			limit:     0,
			total:     100,
			wantLimit: 20,
			wantPage:  1,
			wantTotal: 100,
			wantPages: 5,
			wantSkip:  0,
		},
		{
			name:      "custom page size",
			page:      2,
			limit:     10,
			total:     50,
			wantLimit: 10,
			wantPage:  2,
			wantTotal: 50,
			wantPages: 5,
			wantSkip:  10,
		},
		{
			name:      "page size exceeds max",
			page:      1,
			limit:     100,
			total:     200,
			wantLimit: DefaultPageSize,
			wantPage:  1,
			wantTotal: 200,
			wantPages: 200 / DefaultPageSize,
			wantSkip:  0,
		},
		{
			name:      "page zero resets to one",
			page:      0,
			limit:     10,
			total:     50,
			wantLimit: 10,
			wantPage:  1,
			wantTotal: 50,
			wantPages: 5,
			wantSkip:  0,
		},
		{
			name:      "negative page resets to one",
			page:      -5,
			limit:     10,
			total:     50,
			wantLimit: 10,
			wantPage:  1,
			wantTotal: 50,
			wantPages: 5,
			wantSkip:  0,
		},
		{
			name:      "zero total",
			page:      1,
			limit:     10,
			total:     0,
			wantLimit: 10,
			wantPage:  1,
			wantTotal: 0,
			wantPages: 0,
			wantSkip:  0,
		},
		{
			name:      "last page",
			page:      3,
			limit:     10,
			total:     25,
			wantLimit: 10,
			wantPage:  3,
			wantTotal: 25,
			wantPages: 3,
			wantSkip:  20,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.page, tc.limit, tc.total)
			assert.Equal(t, tc.wantLimit, p.Limit)
			assert.Equal(t, tc.wantPage, p.CurrentPage)
			assert.Equal(t, tc.wantTotal, p.Total)
			assert.Equal(t, tc.wantPages, p.TotalPage)
			assert.Equal(t, tc.wantSkip, p.Skip)
		})
	}
}
