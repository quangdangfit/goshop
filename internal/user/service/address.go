package service

import (
	"context"

	"github.com/quangdangfit/gocommon/validation"

	"goshop/internal/user/dto"
	"goshop/internal/user/model"
	"goshop/internal/user/repository"
	"goshop/pkg/utils"
)

//go:generate mockery --name=AddressService
type AddressService interface {
	ListAddresses(ctx context.Context, userID string) ([]*model.Address, error)
	GetAddressByID(ctx context.Context, id, userID string) (*model.Address, error)
	Create(ctx context.Context, userID string, req *dto.CreateAddressReq) (*model.Address, error)
	Update(ctx context.Context, id, userID string, req *dto.UpdateAddressReq) (*model.Address, error)
	Delete(ctx context.Context, id, userID string) error
	SetDefault(ctx context.Context, id, userID string) error
}

type addressSvc struct {
	validator validation.Validation
	repo      repository.AddressRepository
}

func NewAddressService(validator validation.Validation, repo repository.AddressRepository) AddressService {
	return &addressSvc{validator: validator, repo: repo}
}

func (s *addressSvc) ListAddresses(ctx context.Context, userID string) ([]*model.Address, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *addressSvc) GetAddressByID(ctx context.Context, id, userID string) (*model.Address, error) {
	return s.repo.GetByID(ctx, id, userID)
}

func (s *addressSvc) Create(ctx context.Context, userID string, req *dto.CreateAddressReq) (*model.Address, error) {
	if err := s.validator.ValidateStruct(req); err != nil {
		return nil, err
	}
	var address model.Address
	utils.Copy(&address, req)
	address.UserID = userID
	if err := s.repo.Create(ctx, &address); err != nil {
		return nil, err
	}
	return &address, nil
}

func (s *addressSvc) Update(ctx context.Context, id, userID string, req *dto.UpdateAddressReq) (*model.Address, error) {
	address, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	utils.Copy(address, req)
	if err := s.repo.Update(ctx, address); err != nil {
		return nil, err
	}
	return address, nil
}

func (s *addressSvc) Delete(ctx context.Context, id, userID string) error {
	return s.repo.Delete(ctx, id, userID)
}

func (s *addressSvc) SetDefault(ctx context.Context, id, userID string) error {
	return s.repo.SetDefault(ctx, id, userID)
}
