package paging

import (
	"math"
)

const (
	DefaultPageSize int64 = 20
)

// Pagination struct
type Pagination struct {
	CurrentPage int64 `json:"current_page"`
	Total       int64 `json:"total"`
	TotalPage   int64 `json:"total_page"`
	Limit       int64 `json:"limit"`
	Skip        int64 `json:"skip"`
}

// New paging object
func New(page int64, pageSize int64, total int64) *Pagination {
	var pageInfo Pagination
	limit := DefaultPageSize
	if pageSize > 0 && pageSize <= limit {
		pageInfo.Limit = pageSize
	} else {
		pageInfo.Limit = limit
	}

	totalPage := int64(math.Ceil(float64(total) / float64(pageInfo.Limit)))
	pageInfo.Total = total
	pageInfo.TotalPage = totalPage
	if page < 1 || totalPage == 0 {
		page = 1
	}

	pageInfo.CurrentPage = page
	pageInfo.Skip = (page - 1) * pageInfo.Limit
	return &pageInfo
}
