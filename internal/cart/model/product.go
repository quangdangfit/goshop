package model

type Product struct {
	ID          string  `json:"id" gorm:"unique;not null;index;primary_key"`
	Code        string  `json:"code" gorm:"uniqueIndex:idx_product_code,not null"`
	Name        string  `json:"name" gorm:"uniqueIndex:idx_product_name,not null"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}
