package service

import (
	"context"
	"errors"
	"testing"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"goshop/internal/product/dto"
	"goshop/internal/product/model"
	"goshop/internal/product/repository/mocks"
	"goshop/pkg/config"
)

type CategoryServiceTestSuite struct {
	suite.Suite
	mockRepo *mocks.CategoryRepository
	service  CategoryService
}

func (suite *CategoryServiceTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)
	validator := validation.New()
	suite.mockRepo = mocks.NewCategoryRepository(suite.T())
	suite.service = NewCategoryService(validator, suite.mockRepo)
}

func TestCategoryServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CategoryServiceTestSuite))
}

func (suite *CategoryServiceTestSuite) TestListCategories() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
		wantLen int
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockRepo.On("List", mock.Anything).
					Return([]*model.Category{{ID: "cat1", Name: "Electronics"}, {ID: "cat2", Name: "Clothing"}}, nil).Times(1)
			},
			wantLen: 2,
		},
		{
			name: "DB error",
			setup: func() {
				suite.mockRepo.On("List", mock.Anything).
					Return(nil, errors.New("db error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			categories, err := suite.service.ListCategories(context.Background())
			if tc.wantErr {
				suite.NotNil(err)
				suite.Nil(categories)
			} else {
				suite.Nil(err)
				suite.Equal(tc.wantLen, len(categories))
			}
		})
	}
}

func (suite *CategoryServiceTestSuite) TestGetCategoryByID() {
	tests := []struct {
		name    string
		id      string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			id:   "cat1",
			setup: func() {
				suite.mockRepo.On("GetByID", mock.Anything, "cat1").
					Return(&model.Category{ID: "cat1", Name: "Electronics"}, nil).Times(1)
			},
		},
		{
			name: "Not found",
			id:   "notfound",
			setup: func() {
				suite.mockRepo.On("GetByID", mock.Anything, "notfound").
					Return(nil, errors.New("not found")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			category, err := suite.service.GetCategoryByID(context.Background(), tc.id)
			if tc.wantErr {
				suite.NotNil(err)
				suite.Nil(category)
			} else {
				suite.Nil(err)
				suite.Equal(tc.id, category.ID)
			}
		})
	}
}

func (suite *CategoryServiceTestSuite) TestCreate() {
	tests := []struct {
		name    string
		req     *dto.CreateCategoryReq
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			req:  &dto.CreateCategoryReq{Name: "Electronics", Slug: "electronics"},
			setup: func() {
				suite.mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name:    "Validation fail",
			req:     &dto.CreateCategoryReq{},
			setup:   func() {},
			wantErr: true,
		},
		{
			name: "DB fail",
			req:  &dto.CreateCategoryReq{Name: "Electronics", Slug: "electronics"},
			setup: func() {
				suite.mockRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			category, err := suite.service.Create(context.Background(), tc.req)
			if tc.wantErr {
				suite.NotNil(err)
				suite.Nil(category)
			} else {
				suite.Nil(err)
				suite.NotNil(category)
				suite.Equal("Electronics", category.Name)
			}
		})
	}
}

func (suite *CategoryServiceTestSuite) TestUpdate() {
	tests := []struct {
		name    string
		id      string
		req     *dto.UpdateCategoryReq
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			id:   "cat1",
			req:  &dto.UpdateCategoryReq{Name: "Updated Electronics"},
			setup: func() {
				suite.mockRepo.On("GetByID", mock.Anything, "cat1").
					Return(&model.Category{ID: "cat1", Name: "Electronics"}, nil).Times(1)
				suite.mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name: "GetByID fail",
			id:   "notfound",
			req:  &dto.UpdateCategoryReq{Name: "Updated"},
			setup: func() {
				suite.mockRepo.On("GetByID", mock.Anything, "notfound").
					Return(nil, errors.New("not found")).Times(1)
			},
			wantErr: true,
		},
		{
			name: "DB fail",
			id:   "cat1",
			req:  &dto.UpdateCategoryReq{Name: "Updated Electronics"},
			setup: func() {
				suite.mockRepo.On("GetByID", mock.Anything, "cat1").
					Return(&model.Category{ID: "cat1", Name: "Electronics"}, nil).Times(1)
				suite.mockRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			category, err := suite.service.Update(context.Background(), tc.id, tc.req)
			if tc.wantErr {
				suite.NotNil(err)
				suite.Nil(category)
			} else {
				suite.Nil(err)
				suite.NotNil(category)
				suite.Equal("Updated Electronics", category.Name)
			}
		})
	}
}

func (suite *CategoryServiceTestSuite) TestDelete() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockRepo.On("Delete", mock.Anything, "cat1").Return(nil).Times(1)
			},
		},
		{
			name: "Not found",
			setup: func() {
				suite.mockRepo.On("Delete", mock.Anything, "cat1").Return(errors.New("not found")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			err := suite.service.Delete(context.Background(), "cat1")
			if tc.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}
