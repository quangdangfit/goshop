package category

import (
	"errors"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"goshop/dbs"
)

type Repository interface {
	GetCategories(active bool) (*[]Category, error)
	GetCategoryByID(uuid string) (*Category, error)
	CreateCategory(req *CategoryRequest) (*Category, error)
	UpdateCategory(uuid string, req *CategoryRequest) (*Category, error)
}

type repo struct {
	db *gorm.DB
}

func NewRepository() Repository {
	return &repo{db: dbs.Database}
}

func (r *repo) GetCategories(active bool) (*[]Category, error) {
	var categories []Category
	if r.db.Where("active = ?", active).Find(&categories).RecordNotFound() {
		return nil, nil
	}

	return &categories, nil
}

func (r *repo) GetCategoryByID(uuid string) (*Category, error) {
	var category Category
	if r.db.Where("uuid = ?", uuid).Find(&category).RecordNotFound() {
		return nil, errors.New("not found category")
	}

	return &category, nil
}

func (r *repo) CreateCategory(req *CategoryRequest) (*Category, error) {
	var category Category
	copier.Copy(&category, &req)

	r.db.Create(&category)

	return &category, nil
}

func (r *repo) UpdateCategory(uuid string, req *CategoryRequest) (*Category, error) {
	var category Category
	if r.db.Where("uuid = ? ", uuid).First(&category).RecordNotFound() {
		return nil, errors.New("not found category")
	}

	copier.Copy(&category, &req)
	r.db.Save(&category)

	return &category, nil
}
