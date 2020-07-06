package order

import (
	"github.com/jinzhu/gorm"
)

type Order struct {
	gorm.Model
	UUID       string
	Code       string
	TotalPrice uint
}
