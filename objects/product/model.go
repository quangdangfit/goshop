package product

import (
	"github.com/jinzhu/gorm"
)

type Product struct {
	gorm.Model
	UUID        string
	Code        string
	Name        string
	Description uint
	CategUUID   string
}
