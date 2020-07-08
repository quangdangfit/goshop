package repositories

import (
	"encoding/json"
	"errors"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"goshop/dbs"
	"goshop/models"
)

type CategoryRepository interface {
	GetCategories(active bool) (*[]models.Category, error)
	GetCategoryByID(uuid string) (*models.Category, error)
	CreateCategory(req *models.CategoryRequest) (*models.Category, error)
	UpdateCategory(uuid string, req *models.CategoryRequest) (*models.Category, error)
}

type categRepo struct {
	Repository
	db *gorm.DB
}

func NewCategoryRepository() CategoryRepository {
	return &categRepo{db: dbs.Database}
}

func (r *categRepo) GetCategories(active bool) (*[]models.Category, error) {
	var categories []models.Category
	if r.db.Where("active = ?", active).Find(&categories).RecordNotFound() {
		return nil, nil
	}

	return &categories, nil
}

func (r *categRepo) GetCategoryByID(uuid string) (*models.Category, error) {
	var category models.Category
	value := r.GetCache(uuid)
	if value != nil {
		json.Unmarshal(value, &category)
		return &category, nil
	}

	if r.db.Where("uuid = ?", uuid).Find(&category).RecordNotFound() {
		return nil, errors.New("not found category")
	}

	go r.SetCache(uuid, category)
	return &category, nil
}

func (r *categRepo) CreateCategory(req *models.CategoryRequest) (*models.Category, error) {
	var category models.Category
	copier.Copy(&category, &req)

	if err := r.db.Create(&category).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *categRepo) UpdateCategory(uuid string, req *models.CategoryRequest) (*models.Category, error) {
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
