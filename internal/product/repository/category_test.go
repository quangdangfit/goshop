package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"goshop/internal/product/model"
	"goshop/pkg/config"
	dbsMocks "goshop/pkg/dbs/mocks"
)

type CategoryRepositoryTestSuite struct {
	suite.Suite
	mockDB *dbsMocks.Database
	repo   CategoryRepository
}

func (suite *CategoryRepositoryTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)
	suite.mockDB = dbsMocks.NewDatabase(suite.T())
	suite.repo = NewCategoryRepository(suite.mockDB)
}

func TestCategoryRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CategoryRepositoryTestSuite))
}

func (suite *CategoryRepositoryTestSuite) TestList() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockDB.On("Find", mock.Anything, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name: "DB error",
			setup: func() {
				suite.mockDB.On("Find", mock.Anything, mock.Anything).Return(errors.New("db error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			categories, err := suite.repo.List(context.Background())
			if tc.wantErr {
				suite.NotNil(err)
				suite.Nil(categories)
			} else {
				suite.Nil(err)
				suite.Equal(0, len(categories))
			}
		})
	}
}

func (suite *CategoryRepositoryTestSuite) TestGetByID() {
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
				suite.mockDB.On("FindById", mock.Anything, "cat1", &model.Category{}).Return(nil).Times(1)
			},
		},
		{
			name: "Not found",
			id:   "notfound",
			setup: func() {
				suite.mockDB.On("FindById", mock.Anything, "notfound", &model.Category{}).Return(errors.New("not found")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			category, err := suite.repo.GetByID(context.Background(), tc.id)
			if tc.wantErr {
				suite.NotNil(err)
				suite.Nil(category)
			} else {
				suite.Nil(err)
				suite.NotNil(category)
			}
		})
	}
}

func (suite *CategoryRepositoryTestSuite) TestCreate() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockDB.On("Create", mock.Anything, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name: "DB error",
			setup: func() {
				suite.mockDB.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			category := &model.Category{Name: "Electronics", Slug: "electronics"}
			err := suite.repo.Create(context.Background(), category)
			if tc.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}

func (suite *CategoryRepositoryTestSuite) TestUpdate() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockDB.On("Update", mock.Anything, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name: "DB error",
			setup: func() {
				suite.mockDB.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			category := &model.Category{ID: "cat1", Name: "Updated Electronics"}
			err := suite.repo.Update(context.Background(), category)
			if tc.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}

func (suite *CategoryRepositoryTestSuite) TestDelete() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockDB.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name: "DB error",
			setup: func() {
				suite.mockDB.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("db error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			err := suite.repo.Delete(context.Background(), "cat1")
			if tc.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}
