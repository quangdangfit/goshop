package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"goshop/internal/product/dto"
	"goshop/internal/product/model"
	srvMocks "goshop/internal/product/service/mocks"
	"goshop/pkg/config"
	"goshop/pkg/paging"
	"goshop/pkg/response"
	"goshop/pkg/utils"
)

type ReviewHandlerTestSuite struct {
	suite.Suite
	mockService *srvMocks.ReviewService
	handler     *ReviewHandler
}

func (suite *ReviewHandlerTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)
	suite.mockService = srvMocks.NewReviewService(suite.T())
	suite.handler = NewReviewHandler(suite.mockService)
}

func TestReviewHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ReviewHandlerTestSuite))
}

func (suite *ReviewHandlerTestSuite) prepareContext(method, path string, body any) (*gin.Context, *httptest.ResponseRecorder) {
	requestBody, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, path, bytes.NewBuffer(requestBody))
	c, _ := gin.CreateTestContext(w)
	c.Request = r
	return c, w
}

// ListReviews
// =================================================================================================

func (suite *ReviewHandlerTestSuite) TestListReviewsSuccess() {
	ctx, writer := suite.prepareContext("GET", "/api/v1/products/p1/reviews", nil)
	ctx.Params = gin.Params{{Key: "id", Value: "p1"}}

	suite.mockService.On("ListReviews", mock.Anything, "p1", int64(0), int64(0)).
		Return(
			[]*model.Review{
				{ID: "r1", ProductID: "p1", UserID: "u1", Rating: 5, Comment: "Great!"},
				{ID: "r2", ProductID: "p1", UserID: "u2", Rating: 4, Comment: "Good"},
			},
			&paging.Pagination{Total: 2, CurrentPage: 1, Limit: 20},
			nil,
		).Times(1)

	suite.handler.ListReviews(ctx)

	var res response.Response
	var reviews dto.ListReviewRes
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	utils.Copy(&reviews, &res.Result)
	suite.Equal(http.StatusOK, writer.Code)
	suite.Equal(2, len(reviews.Reviews))
	suite.Equal("r1", reviews.Reviews[0].ID)
	suite.Equal(5, reviews.Reviews[0].Rating)
}

func (suite *ReviewHandlerTestSuite) TestListReviewsWithPageAndLimit() {
	ctx, writer := suite.prepareContext("GET", "/api/v1/products/p1/reviews?page=2&limit=10", nil)
	ctx.Params = gin.Params{{Key: "id", Value: "p1"}}

	suite.mockService.On("ListReviews", mock.Anything, "p1", int64(2), int64(10)).
		Return(
			[]*model.Review{},
			&paging.Pagination{Total: 0, CurrentPage: 2, Limit: 10},
			nil,
		).Times(1)

	suite.handler.ListReviews(ctx)

	suite.Equal(http.StatusOK, writer.Code)
}

func (suite *ReviewHandlerTestSuite) TestListReviewsFail() {
	ctx, writer := suite.prepareContext("GET", "/api/v1/products/p1/reviews", nil)
	ctx.Params = gin.Params{{Key: "id", Value: "p1"}}

	suite.mockService.On("ListReviews", mock.Anything, "p1", int64(0), int64(0)).
		Return(nil, nil, errors.New("db error")).Times(1)

	suite.handler.ListReviews(ctx)

	suite.Equal(http.StatusInternalServerError, writer.Code)
}

// CreateReview
// =================================================================================================

func (suite *ReviewHandlerTestSuite) TestCreateReviewSuccess() {
	req := &dto.CreateReviewReq{
		Rating:  5,
		Comment: "Excellent product!",
	}
	ctx, writer := suite.prepareContext("POST", "/api/v1/products/p1/reviews", req)
	ctx.Params = gin.Params{{Key: "id", Value: "p1"}}
	ctx.Set("userId", "u1")

	suite.mockService.On("CreateReview", mock.Anything, "p1", "u1", req).
		Return(&model.Review{
			ID:        "r1",
			ProductID: "p1",
			UserID:    "u1",
			Rating:    5,
			Comment:   "Excellent product!",
		}, nil).Times(1)

	suite.handler.CreateReview(ctx)

	var res response.Response
	var review dto.Review
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	utils.Copy(&review, &res.Result)
	suite.Equal(http.StatusOK, writer.Code)
	suite.Equal("r1", review.ID)
	suite.Equal(5, review.Rating)
	suite.Equal("Excellent product!", review.Comment)
}

func (suite *ReviewHandlerTestSuite) TestCreateReviewUnauthorized() {
	req := &dto.CreateReviewReq{Rating: 5}
	ctx, writer := suite.prepareContext("POST", "/api/v1/products/p1/reviews", req)
	ctx.Params = gin.Params{{Key: "id", Value: "p1"}}
	// userId not set

	suite.handler.CreateReview(ctx)

	suite.Equal(http.StatusUnauthorized, writer.Code)
}

func (suite *ReviewHandlerTestSuite) TestCreateReviewInvalidBody() {
	req := map[string]any{"rating": "five"}
	ctx, writer := suite.prepareContext("POST", "/api/v1/products/p1/reviews", req)
	ctx.Params = gin.Params{{Key: "id", Value: "p1"}}
	ctx.Set("userId", "u1")

	suite.handler.CreateReview(ctx)

	suite.Equal(http.StatusBadRequest, writer.Code)
}

func (suite *ReviewHandlerTestSuite) TestCreateReviewFail() {
	req := &dto.CreateReviewReq{Rating: 5, Comment: "Great"}
	ctx, writer := suite.prepareContext("POST", "/api/v1/products/p1/reviews", req)
	ctx.Params = gin.Params{{Key: "id", Value: "p1"}}
	ctx.Set("userId", "u1")

	suite.mockService.On("CreateReview", mock.Anything, "p1", "u1", req).
		Return(nil, errors.New("already reviewed")).Times(1)

	suite.handler.CreateReview(ctx)

	suite.Equal(http.StatusInternalServerError, writer.Code)
}

// UpdateReview
// =================================================================================================

func (suite *ReviewHandlerTestSuite) TestUpdateReviewSuccess() {
	req := &dto.UpdateReviewReq{
		Rating:  4,
		Comment: "Updated comment",
	}
	ctx, writer := suite.prepareContext("PUT", "/api/v1/products/p1/reviews/r1", req)
	ctx.Params = gin.Params{{Key: "id", Value: "p1"}, {Key: "reviewId", Value: "r1"}}
	ctx.Set("userId", "u1")

	suite.mockService.On("UpdateReview", mock.Anything, "r1", "u1", req).
		Return(&model.Review{
			ID:        "r1",
			ProductID: "p1",
			UserID:    "u1",
			Rating:    4,
			Comment:   "Updated comment",
		}, nil).Times(1)

	suite.handler.UpdateReview(ctx)

	var res response.Response
	var review dto.Review
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	utils.Copy(&review, &res.Result)
	suite.Equal(http.StatusOK, writer.Code)
	suite.Equal("r1", review.ID)
	suite.Equal(4, review.Rating)
	suite.Equal("Updated comment", review.Comment)
}

func (suite *ReviewHandlerTestSuite) TestUpdateReviewUnauthorized() {
	req := &dto.UpdateReviewReq{Rating: 4}
	ctx, writer := suite.prepareContext("PUT", "/api/v1/products/p1/reviews/r1", req)
	ctx.Params = gin.Params{{Key: "id", Value: "p1"}, {Key: "reviewId", Value: "r1"}}
	// userId not set

	suite.handler.UpdateReview(ctx)

	suite.Equal(http.StatusUnauthorized, writer.Code)
}

func (suite *ReviewHandlerTestSuite) TestUpdateReviewInvalidBody() {
	req := map[string]any{"rating": "four"}
	ctx, writer := suite.prepareContext("PUT", "/api/v1/products/p1/reviews/r1", req)
	ctx.Params = gin.Params{{Key: "id", Value: "p1"}, {Key: "reviewId", Value: "r1"}}
	ctx.Set("userId", "u1")

	suite.handler.UpdateReview(ctx)

	suite.Equal(http.StatusBadRequest, writer.Code)
}

func (suite *ReviewHandlerTestSuite) TestUpdateReviewFail() {
	req := &dto.UpdateReviewReq{Comment: "Updated"}
	ctx, writer := suite.prepareContext("PUT", "/api/v1/products/p1/reviews/r1", req)
	ctx.Params = gin.Params{{Key: "id", Value: "p1"}, {Key: "reviewId", Value: "r1"}}
	ctx.Set("userId", "u1")

	suite.mockService.On("UpdateReview", mock.Anything, "r1", "u1", req).
		Return(nil, errors.New("not found")).Times(1)

	suite.handler.UpdateReview(ctx)

	suite.Equal(http.StatusInternalServerError, writer.Code)
}

// DeleteReview
// =================================================================================================

func (suite *ReviewHandlerTestSuite) TestDeleteReviewSuccess() {
	ctx, writer := suite.prepareContext("DELETE", "/api/v1/products/p1/reviews/r1", nil)
	ctx.Params = gin.Params{{Key: "id", Value: "p1"}, {Key: "reviewId", Value: "r1"}}
	ctx.Set("userId", "u1")

	suite.mockService.On("DeleteReview", mock.Anything, "r1", "u1").Return(nil).Times(1)

	suite.handler.DeleteReview(ctx)

	suite.Equal(http.StatusOK, writer.Code)
}

func (suite *ReviewHandlerTestSuite) TestDeleteReviewUnauthorized() {
	ctx, writer := suite.prepareContext("DELETE", "/api/v1/products/p1/reviews/r1", nil)
	ctx.Params = gin.Params{{Key: "id", Value: "p1"}, {Key: "reviewId", Value: "r1"}}
	// userId not set

	suite.handler.DeleteReview(ctx)

	suite.Equal(http.StatusUnauthorized, writer.Code)
}

func (suite *ReviewHandlerTestSuite) TestDeleteReviewFail() {
	ctx, writer := suite.prepareContext("DELETE", "/api/v1/products/p1/reviews/r1", nil)
	ctx.Params = gin.Params{{Key: "id", Value: "p1"}, {Key: "reviewId", Value: "r1"}}
	ctx.Set("userId", "u1")

	suite.mockService.On("DeleteReview", mock.Anything, "r1", "u1").Return(errors.New("not found")).Times(1)

	suite.handler.DeleteReview(ctx)

	suite.Equal(http.StatusInternalServerError, writer.Code)
}
