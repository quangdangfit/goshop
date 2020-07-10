package repositories

import (
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"

	"goshop/dbs"
	"goshop/models"
)

type RoleRepository interface {
	CreateRole(req *models.RoleRequest) (*models.Role, error)
}

type roleRepo struct {
	db *gorm.DB
}

func NewRoleRepository() RoleRepository {
	return &roleRepo{db: dbs.Database}
}

func (r *roleRepo) CreateRole(req *models.RoleRequest) (*models.Role, error) {
	var role models.Role
	copier.Copy(&role, &req)

	if err := r.db.Create(&role).Error; err != nil {
		return nil, err
	}

	return &role, nil
}
