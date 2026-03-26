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
	"goshop/pkg/paging"
)

type ReviewServiceTestSuite struct {
	suite.Suite
	mockRepo        *mocks.ReviewRepository
	mockProductRepo *mocks.ProductRepository
	service         ReviewService
}

func (suite *ReviewServiceTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)
	validator := validation.New()
	suite.mockRepo = mocks.NewReviewRepository(suite.T())
	suite.mockProductRepo = mocks.NewProductRepository(suite.T())
	suite.service = NewReviewService(validator, suite.mockRepo, suite.mockProductRepo)
}

func TestReviewServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ReviewServiceTestSuite))
}

func (suite *ReviewServiceTestSuite) TestListReviews() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
		wantLen int
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockRepo.On("ListByProduct", mock.Anything, "p1", int64(1), int64(10)).
					Return([]*model.Review{
						{ID: "r1", ProductID: "p1", UserID: "u1", Rating: 5},
						{ID: "r2", ProductID: "p1", UserID: "u2", Rating: 4},
					}, &paging.Pagination{Total: 2}, nil).Times(1)
			},
			wantLen: 2,
		},
		{
			name: "DB error",
			setup: func() {
				suite.mockRepo.On("ListByProduct", mock.Anything, "p1", int64(1), int64(10)).
					Return(nil, nil, errors.New("db error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			reviews, pagination, err := suite.service.ListReviews(context.Background(), "p1", 1, 10)
			if tc.wantErr {
				suite.NotNil(err)
				suite.Nil(reviews)
				suite.Nil(pagination)
			} else {
				suite.Nil(err)
				suite.Equal(tc.wantLen, len(reviews))
				suite.NotNil(pagination)
			}
		})
	}
}

func (suite *ReviewServiceTestSuite) TestCreateReview() {
	tests := []struct {
		name    string
		req     *dto.CreateReviewReq
		setup   func()
		wantErr bool
		wantNil bool
	}{
		{
			name: "Success",
			req:  &dto.CreateReviewReq{Rating: 5, Comment: "Great!"},
			setup: func() {
				suite.mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Times(1)
				suite.mockRepo.On("GetAggregates", mock.Anything, "p1").
					Return(float64(5), 1, nil).Maybe()
				suite.mockProductRepo.On("UpdateRating", mock.Anything, "p1", float64(5), 1).
					Return(nil).Maybe()
			},
		},
		{
			name:    "Validation fail",
			req:     &dto.CreateReviewReq{},
			setup:   func() {},
			wantErr: true,
			wantNil: true,
		},
		{
			name: "DB fail",
			req:  &dto.CreateReviewReq{Rating: 5, Comment: "Great!"},
			setup: func() {
				suite.mockRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error")).Times(1)
			},
			wantErr: true,
			wantNil: true,
		},
		{
			name: "GetAggregates fail",
			req:  &dto.CreateReviewReq{Rating: 5, Comment: "Great!"},
			setup: func() {
				suite.mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Times(1)
				suite.mockRepo.On("GetAggregates", mock.Anything, "p1").
					Return(float64(0), 0, errors.New("db error")).Maybe()
			},
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			review, err := suite.service.CreateReview(context.Background(), "p1", "u1", tc.req)
			if tc.wantErr {
				suite.NotNil(err)
				suite.Nil(review)
			} else {
				suite.Nil(err)
				suite.NotNil(review)
				suite.Equal("p1", review.ProductID)
				suite.Equal("u1", review.UserID)
			}
		})
	}
}

func (suite *ReviewServiceTestSuite) TestUpdateReview() {
	tests := []struct {
		name    string
		id      string
		userID  string
		req     *dto.UpdateReviewReq
		setup   func()
		wantErr bool
		errMsg  string
	}{
		{
			name:   "Success",
			id:     "r1",
			userID: "u1",
			req:    &dto.UpdateReviewReq{Rating: 4, Comment: "Updated"},
			setup: func() {
				suite.mockRepo.On("GetByID", mock.Anything, "r1").
					Return(&model.Review{ID: "r1", ProductID: "p1", UserID: "u1", Rating: 5}, nil).Times(1)
				suite.mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Times(1)
				suite.mockRepo.On("GetAggregates", mock.Anything, "p1").
					Return(float64(4), 1, nil).Maybe()
				suite.mockProductRepo.On("UpdateRating", mock.Anything, "p1", float64(4), 1).
					Return(nil).Maybe()
			},
		},
		{
			name:   "GetByID fail",
			id:     "notfound",
			userID: "u1",
			req:    &dto.UpdateReviewReq{Rating: 4},
			setup: func() {
				suite.mockRepo.On("GetByID", mock.Anything, "notfound").
					Return(nil, errors.New("not found")).Times(1)
			},
			wantErr: true,
		},
		{
			name:   "Permission denied",
			id:     "r1",
			userID: "u1",
			req:    &dto.UpdateReviewReq{Rating: 4},
			setup: func() {
				suite.mockRepo.On("GetByID", mock.Anything, "r1").
					Return(&model.Review{ID: "r1", ProductID: "p1", UserID: "u2"}, nil).Times(1)
			},
			wantErr: true,
			errMsg:  "Permission denied",
		},
		{
			name:   "DB fail",
			id:     "r1",
			userID: "u1",
			req:    &dto.UpdateReviewReq{Rating: 4},
			setup: func() {
				suite.mockRepo.On("GetByID", mock.Anything, "r1").
					Return(&model.Review{ID: "r1", ProductID: "p1", UserID: "u1", Rating: 5}, nil).Times(1)
				suite.mockRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			review, err := suite.service.UpdateReview(context.Background(), tc.id, tc.userID, tc.req)
			if tc.wantErr {
				suite.NotNil(err)
				suite.Nil(review)
				if tc.errMsg != "" {
					suite.Equal(tc.errMsg, err.Error())
				}
			} else {
				suite.Nil(err)
				suite.NotNil(review)
				suite.Equal(tc.req.Rating, review.Rating)
			}
		})
	}
}

func (suite *ReviewServiceTestSuite) TestDeleteReview() {
	tests := []struct {
		name    string
		id      string
		userID  string
		setup   func()
		wantErr bool
		errMsg  string
	}{
		{
			name:   "Success",
			id:     "r1",
			userID: "u1",
			setup: func() {
				suite.mockRepo.On("GetByID", mock.Anything, "r1").
					Return(&model.Review{ID: "r1", ProductID: "p1", UserID: "u1"}, nil).Times(1)
				suite.mockRepo.On("Delete", mock.Anything, "r1", "u1").Return(nil).Times(1)
				suite.mockRepo.On("GetAggregates", mock.Anything, "p1").
					Return(float64(0), 0, nil).Maybe()
				suite.mockProductRepo.On("UpdateRating", mock.Anything, "p1", float64(0), 0).
					Return(nil).Maybe()
			},
		},
		{
			name:   "GetByID fail",
			id:     "notfound",
			userID: "u1",
			setup: func() {
				suite.mockRepo.On("GetByID", mock.Anything, "notfound").
					Return(nil, errors.New("not found")).Times(1)
			},
			wantErr: true,
		},
		{
			name:   "Permission denied",
			id:     "r1",
			userID: "u1",
			setup: func() {
				suite.mockRepo.On("GetByID", mock.Anything, "r1").
					Return(&model.Review{ID: "r1", ProductID: "p1", UserID: "u2"}, nil).Times(1)
			},
			wantErr: true,
			errMsg:  "Permission denied",
		},
		{
			name:   "DB fail",
			id:     "r1",
			userID: "u1",
			setup: func() {
				suite.mockRepo.On("GetByID", mock.Anything, "r1").
					Return(&model.Review{ID: "r1", ProductID: "p1", UserID: "u1"}, nil).Times(1)
				suite.mockRepo.On("Delete", mock.Anything, "r1", "u1").Return(errors.New("db error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			err := suite.service.DeleteReview(context.Background(), tc.id, tc.userID)
			if tc.wantErr {
				suite.NotNil(err)
				if tc.errMsg != "" {
					suite.Equal(tc.errMsg, err.Error())
				}
			} else {
				suite.Nil(err)
			}
		})
	}
}
