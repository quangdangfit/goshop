package repositories

import (
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"

	"goshop/app/models"
	"goshop/app/serializers"
	"goshop/dbs"
)

type IRoleRepository interface {
	CreateRole(req *serializers.RoleBodyParam) (*models.Role, error)
}

type RoleRepo struct {
	db *gorm.DB
}

func NewRoleRepository() *RoleRepo {
	return &RoleRepo{db: dbs.Database}
}

func (r *RoleRepo) CreateRole(req *serializers.RoleBodyParam) (*models.Role, error) {
	var role models.Role
	copier.Copy(&role, &req)

	if err := r.db.Create(&role).Error; err != nil {
		return nil, err
	}

	return &role, nil
}
