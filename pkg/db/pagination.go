package db

import (
	"math"

	"gorm.io/gorm"
)

type SortType string

const (
	SORT_ASC  SortType = "asc"
	SORT_DESC SortType = "desc"
)

func ToSortType(sortTypeStr string) SortType {
	if sortTypeStr == "asc" {
		return SORT_ASC
	}
	return SORT_DESC
}

// Pagination is the model for pagination information.
type Pagination struct {
	Limit      int    `json:"limit,omitempty" example:"10"`
	Page       int    `json:"page,omitempty" example:"1"`
	Sort       string `json:"sort,omitempty" example:"Id desc"`
	TotalRows  int64  `json:"total_rows" example:"100"`
	TotalPages int    `json:"total_pages" example:"10"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetSort() string {
	if p.Sort == "" {
		p.Sort = "Id desc"
	}
	return p.Sort
}

func Paginate(value interface{}, pagination *Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(value).Count(&totalRows)

	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}
