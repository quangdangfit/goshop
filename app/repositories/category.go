package repositories

import (
	"errors"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"

	"goshop/app/models"
	"goshop/dbs"
)

type CategoryRepository interface {
	GetCategories(query models.CategoryQueryRequest) (*[]models.Category, error)
	GetCategoryByID(uuid string) (*models.Category, error)
	CreateCategory(req *models.CategoryBodyRequest) (*models.Category, error)
	UpdateCategory(uuid string, req *models.CategoryBodyRequest) (*models.Category, error)
}

type categRepo struct {
	db *gorm.DB
}

func NewCategoryRepository() CategoryRepository {
	return &categRepo{db: dbs.Database}
}

func (r *categRepo) GetCategories(query models.CategoryQueryRequest) (*[]models.Category, error) {
	var categories []models.Category
	if r.db.Where(query).Find(&categories).RecordNotFound() {
		return nil, nil
	}

	return &categories, nil
}

func (r *categRepo) GetCategoryByID(uuid string) (*models.Category, error) {
	var category models.Category
	if r.db.Where("uuid = ?", uuid).Find(&category).RecordNotFound() {
		return nil, errors.New("not found category")
	}

	return &category, nil
}

func (r *categRepo) CreateCategory(req *models.CategoryBodyRequest) (*models.Category, error) {
	var category models.Category
	copier.Copy(&category, &req)

	if err := r.db.Create(&category).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *categRepo) UpdateCategory(uuid string, req *models.CategoryBodyRequest) (*models.Category, error) {
	var category models.Category
	if r.db.Where("uuid = ? ", uuid).First(&category).RecordNotFound() {
		return nil, errors.New("not found category")
	}

	copier.Copy(&category, &req)
	if err := r.db.Save(&category).Error; err != nil {
		return nil, err
	}

	return &category, nil
}
