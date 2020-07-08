package models

import (
	"github.com/jinzhu/gorm"
)

type OrderLine struct {
	gorm.Model
	UUID      string
	ProductID string
	OrderID   string
	Quantity  uint
	Price     uint
}
