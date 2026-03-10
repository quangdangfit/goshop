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

// ListReviews
// =================================================================================================

func (suite *ReviewServiceTestSuite) TestListReviewsSuccess() {
	suite.mockRepo.On("ListByProduct", mock.Anything, "p1", int64(1), int64(10)).
		Return([]*model.Review{
			{ID: "r1", ProductID: "p1", UserID: "u1", Rating: 5},
			{ID: "r2", ProductID: "p1", UserID: "u2", Rating: 4},
		}, &paging.Pagination{Total: 2}, nil).Times(1)

	reviews, pagination, err := suite.service.ListReviews(context.Background(), "p1", 1, 10)
	suite.Nil(err)
	suite.Equal(2, len(reviews))
	suite.NotNil(pagination)
}

func (suite *ReviewServiceTestSuite) TestListReviewsFail() {
	suite.mockRepo.On("ListByProduct", mock.Anything, "p1", int64(1), int64(10)).
		Return(nil, nil, errors.New("db error")).Times(1)

	reviews, pagination, err := suite.service.ListReviews(context.Background(), "p1", 1, 10)
	suite.NotNil(err)
	suite.Nil(reviews)
	suite.Nil(pagination)
}

// CreateReview
// =================================================================================================

func (suite *ReviewServiceTestSuite) TestCreateReviewSuccess() {
	req := &dto.CreateReviewReq{Rating: 5, Comment: "Great!"}

	suite.mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Times(1)
	suite.mockRepo.On("GetAggregates", mock.Anything, "p1").
		Return(float64(5), 1, nil).Maybe()
	suite.mockProductRepo.On("UpdateRating", mock.Anything, "p1", float64(5), 1).
		Return(nil).Maybe()

	review, err := suite.service.CreateReview(context.Background(), "p1", "u1", req)
	suite.Nil(err)
	suite.NotNil(review)
	suite.Equal("p1", review.ProductID)
	suite.Equal("u1", review.UserID)
	suite.Equal(5, review.Rating)
}

func (suite *ReviewServiceTestSuite) TestCreateReviewValidationFail() {
	req := &dto.CreateReviewReq{} // missing required Rating

	review, err := suite.service.CreateReview(context.Background(), "p1", "u1", req)
	suite.NotNil(err)
	suite.Nil(review)
}

func (suite *ReviewServiceTestSuite) TestCreateReviewDBFail() {
	req := &dto.CreateReviewReq{Rating: 5, Comment: "Great!"}

	suite.mockRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error")).Times(1)

	review, err := suite.service.CreateReview(context.Background(), "p1", "u1", req)
	suite.NotNil(err)
	suite.Nil(review)
}

func (suite *ReviewServiceTestSuite) TestCreateReviewGetAggregatesFail() {
	req := &dto.CreateReviewReq{Rating: 5, Comment: "Great!"}

	suite.mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Times(1)
	suite.mockRepo.On("GetAggregates", mock.Anything, "p1").
		Return(float64(0), 0, errors.New("db error")).Maybe()

	review, err := suite.service.CreateReview(context.Background(), "p1", "u1", req)
	suite.Nil(err)
	suite.NotNil(review)
}

// UpdateReview
// =================================================================================================

func (suite *ReviewServiceTestSuite) TestUpdateReviewSuccess() {
	req := &dto.UpdateReviewReq{Rating: 4, Comment: "Updated"}

	suite.mockRepo.On("GetByID", mock.Anything, "r1").
		Return(&model.Review{ID: "r1", ProductID: "p1", UserID: "u1", Rating: 5}, nil).Times(1)
	suite.mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Times(1)
	suite.mockRepo.On("GetAggregates", mock.Anything, "p1").
		Return(float64(4), 1, nil).Maybe()
	suite.mockProductRepo.On("UpdateRating", mock.Anything, "p1", float64(4), 1).
		Return(nil).Maybe()

	review, err := suite.service.UpdateReview(context.Background(), "r1", "u1", req)
	suite.Nil(err)
	suite.NotNil(review)
	suite.Equal(4, review.Rating)
}

func (suite *ReviewServiceTestSuite) TestUpdateReviewGetByIDFail() {
	req := &dto.UpdateReviewReq{Rating: 4}

	suite.mockRepo.On("GetByID", mock.Anything, "notfound").
		Return(nil, errors.New("not found")).Times(1)

	review, err := suite.service.UpdateReview(context.Background(), "notfound", "u1", req)
	suite.NotNil(err)
	suite.Nil(review)
}

func (suite *ReviewServiceTestSuite) TestUpdateReviewPermissionDenied() {
	req := &dto.UpdateReviewReq{Rating: 4}

	suite.mockRepo.On("GetByID", mock.Anything, "r1").
		Return(&model.Review{ID: "r1", ProductID: "p1", UserID: "u2"}, nil).Times(1)

	review, err := suite.service.UpdateReview(context.Background(), "r1", "u1", req)
	suite.NotNil(err)
	suite.Nil(review)
	suite.Equal("permission denied", err.Error())
}

func (suite *ReviewServiceTestSuite) TestUpdateReviewDBFail() {
	req := &dto.UpdateReviewReq{Rating: 4}

	suite.mockRepo.On("GetByID", mock.Anything, "r1").
		Return(&model.Review{ID: "r1", ProductID: "p1", UserID: "u1", Rating: 5}, nil).Times(1)
	suite.mockRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error")).Times(1)

	review, err := suite.service.UpdateReview(context.Background(), "r1", "u1", req)
	suite.NotNil(err)
	suite.Nil(review)
}

// DeleteReview
// =================================================================================================

func (suite *ReviewServiceTestSuite) TestDeleteReviewSuccess() {
	suite.mockRepo.On("GetByID", mock.Anything, "r1").
		Return(&model.Review{ID: "r1", ProductID: "p1", UserID: "u1"}, nil).Times(1)
	suite.mockRepo.On("Delete", mock.Anything, "r1", "u1").Return(nil).Times(1)
	suite.mockRepo.On("GetAggregates", mock.Anything, "p1").
		Return(float64(0), 0, nil).Maybe()
	suite.mockProductRepo.On("UpdateRating", mock.Anything, "p1", float64(0), 0).
		Return(nil).Maybe()

	err := suite.service.DeleteReview(context.Background(), "r1", "u1")
	suite.Nil(err)
}

func (suite *ReviewServiceTestSuite) TestDeleteReviewGetByIDFail() {
	suite.mockRepo.On("GetByID", mock.Anything, "notfound").
		Return(nil, errors.New("not found")).Times(1)

	err := suite.service.DeleteReview(context.Background(), "notfound", "u1")
	suite.NotNil(err)
}

func (suite *ReviewServiceTestSuite) TestDeleteReviewPermissionDenied() {
	suite.mockRepo.On("GetByID", mock.Anything, "r1").
		Return(&model.Review{ID: "r1", ProductID: "p1", UserID: "u2"}, nil).Times(1)

	err := suite.service.DeleteReview(context.Background(), "r1", "u1")
	suite.NotNil(err)
	suite.Equal("permission denied", err.Error())
}

func (suite *ReviewServiceTestSuite) TestDeleteReviewDBFail() {
	suite.mockRepo.On("GetByID", mock.Anything, "r1").
		Return(&model.Review{ID: "r1", ProductID: "p1", UserID: "u1"}, nil).Times(1)
	suite.mockRepo.On("Delete", mock.Anything, "r1", "u1").Return(errors.New("db error")).Times(1)

	err := suite.service.DeleteReview(context.Background(), "r1", "u1")
	suite.NotNil(err)
}
