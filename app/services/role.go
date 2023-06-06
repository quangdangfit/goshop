package services

import (
	"context"

	"goshop/app/models"
	"goshop/app/repositories"
	"goshop/app/serializers"
)

type IRoleService interface {
	CreateRole(ctx context.Context, item *serializers.RoleBodyParam) (*models.Role, error)
}

type role struct {
	repo repositories.IRoleRepository
}

func NewRoleService(repo repositories.IRoleRepository) IRoleService {
	return &role{repo: repo}
}

func (r *role) CreateRole(ctx context.Context, item *serializers.RoleBodyParam) (*models.Role, error) {
	role, err := r.repo.CreateRole(item)
	if err != nil {
		return nil, err
	}

	return role, nil
}
