package product

import (
	"goshop/dbs"
)

type Repository interface {
	GetProducts() (*[]Product, error)
	GetProductByID(uuid string) (*Product, error)
	//CreateProduct(payload map[string]interface{}) (*Product, error)
	//UpdateProduct(uuid string, payload map[string]interface{}) (*Product, error)
}

type repo struct {
}

func NewRepository() Repository {
	return &repo{}
}

func (r *repo) GetProducts() (*[]Product, error) {
	var products []Product
	if dbs.Database.Find(&products).RecordNotFound() {
		return nil, nil
	}

	return &products, nil

}

func (r *repo) GetProductByID(uuid string) (*Product, error) {
	var product Product
	if dbs.Database.Where("uuid = ?", uuid).Find(&product).RecordNotFound() {
		return nil, nil
	}

	return &product, nil
}
