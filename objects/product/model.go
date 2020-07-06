package product

import (
	"github.com/jinzhu/gorm"
)

type Product struct {
	gorm.Model
	UUID        string `json:"uuid,omitempty" bson:"uuid,omitempty"`
	Code        string `json:"code,omitempty" bson:"code,omitempty"`
	Name        string `json:"name,omitempty" bson:"name,omitempty"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	CategUUID   string `json:"categ_uuid,omitempty" bson:"code,omitempty"`
}

type ProductResponse struct {
	UUID        string `json:"uuid" bson:"uuid"`
	Code        string `json:"code" bson:"code"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	CategUUID   string `json:"categ_uuid" bson:"code"`
}

type ProductRequest struct {
	Name        string `json:"name,omitempty" bson:"name,omitempty"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	CategUUID   string `json:"categ_uuid,omitempty" bson:"categ_uuid,omitempty"`
}
