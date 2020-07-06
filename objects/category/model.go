package category

import (
	"github.com/jinzhu/gorm"
)

type Category struct {
	gorm.Model
	UUID        string
	Code        string
	Name        string
	Description uint
}
