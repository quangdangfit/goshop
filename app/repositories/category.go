package repositories

import (
	"errors"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"

	"goshop/app/models"
	"goshop/app/schema"
	"goshop/dbs"
	"goshop/pkg/utils"
)

type ICategoryRepository interface {
	GetCategories(query *schema.CategoryQueryParam) (*[]models.Category, error)
	GetCategoryByID(uuid string) (*models.Category, error)
	CreateCategory(item *schema.Category) (*models.Category, error)
	UpdateCategory(uuid string, item *schema.CategoryBodyParam) (*models.Category, error)
}

type categRepo struct {
	db *gorm.DB
}

func NewCategoryRepository() ICategoryRepository {
	return &categRepo{db: dbs.Database}
}

func (r *categRepo) GetCategories(query *schema.CategoryQueryParam) (*[]models.Category, error) {
	var categories []models.Category
	var mapQuery map[string]interface{}
	utils.Copy(&mapQuery, query)
	if r.db.Where(mapQuery).Find(&categories).RecordNotFound() {
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

func (r *categRepo) CreateCategory(item *schema.Category) (*models.Category, error) {
	var category models.Category
	copier.Copy(&category, &item)

	if err := r.db.Create(&category).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *categRepo) UpdateCategory(uuid string, item *schema.CategoryBodyParam) (*models.Category, error) {
	var category models.Category
	var value map[string]interface{}
	if r.db.Where("uuid = ? ", uuid).First(&category).RecordNotFound() {
		return nil, errors.New("not found category")
	}

	utils.Copy(&value, item)
	if err := r.db.Model(&category).Updates(value).Error; err != nil {
		return nil, err
	}

	return &category, nil
}
