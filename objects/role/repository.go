package role

import (
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"goshop/dbs"
)

type Repository interface {
	CreateRole(req *RoleRequest) (*Role, error)
}

type repo struct {
	db *gorm.DB
}

func NewRepository() Repository {
	return &repo{db: dbs.Database}
}

func (r *repo) CreateRole(req *RoleRequest) (*Role, error) {
	var role Role
	copier.Copy(&role, &req)

	if err := r.db.Create(&role).Error; err != nil {
		return nil, err
	}

	return &role, nil
}
