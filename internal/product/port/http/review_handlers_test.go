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

	"goshop/internal/product/domain"
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

func (suite *ReviewHandlerTestSuite) TestListReviews() {
	tests := []struct {
		name      string
		path      string
		setup     func()
		expected  int
		checkBody func(writer *httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			path: "/api/v1/products/p1/reviews",
			setup: func() {
				suite.mockService.On("ListReviews", mock.Anything, "p1", int64(0), int64(0)).
					Return(
						[]*model.Review{
							{ID: "r1", ProductID: "p1", UserID: "u1", Rating: 5, Comment: "Great!"},
							{ID: "r2", ProductID: "p1", UserID: "u2", Rating: 4, Comment: "Good"},
						},
						&paging.Pagination{Total: 2, CurrentPage: 1, Limit: 20},
						nil,
					).Times(1)
			},
			expected: http.StatusOK,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				var reviews domain.ListReviewRes
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				_ = utils.Copy(&reviews, &res.Result)
				suite.Equal(2, len(reviews.Reviews))
				suite.Equal("r1", reviews.Reviews[0].ID)
				suite.Equal(5, reviews.Reviews[0].Rating)
			},
		},
		{
			name: "WithPageAndLimit",
			path: "/api/v1/products/p1/reviews?page=2&limit=10",
			setup: func() {
				suite.mockService.On("ListReviews", mock.Anything, "p1", int64(2), int64(10)).
					Return(
						[]*model.Review{},
						&paging.Pagination{Total: 0, CurrentPage: 2, Limit: 10},
						nil,
					).Times(1)
			},
			expected: http.StatusOK,
		},
		{
			name: "Fail",
			path: "/api/v1/products/p1/reviews",
			setup: func() {
				suite.mockService.On("ListReviews", mock.Anything, "p1", int64(0), int64(0)).
					Return(nil, nil, errors.New("db error")).Times(1)
			},
			expected: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContext("GET", tc.path, nil)
			ctx.Params = gin.Params{{Key: "id", Value: "p1"}}
			tc.setup()
			suite.handler.ListReviews(ctx)
			suite.Equal(tc.expected, writer.Code)
			if tc.checkBody != nil {
				tc.checkBody(writer)
			}
		})
	}
}

func (suite *ReviewHandlerTestSuite) TestCreateReview() {
	tests := []struct {
		name      string
		body      any
		userId    string
		setup     func()
		expected  int
		checkBody func(writer *httptest.ResponseRecorder)
	}{
		{
			name:   "Success",
			body:   &domain.CreateReviewReq{Rating: 5, Comment: "Excellent product!"},
			userId: "u1",
			setup: func() {
				suite.mockService.On("CreateReview", mock.Anything, "p1", "u1", &domain.CreateReviewReq{
					Rating: 5, Comment: "Excellent product!",
				}).Return(&model.Review{
					ID: "r1", ProductID: "p1", UserID: "u1", Rating: 5, Comment: "Excellent product!",
				}, nil).Times(1)
			},
			expected: http.StatusOK,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				var review domain.Review
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				_ = utils.Copy(&review, &res.Result)
				suite.Equal("r1", review.ID)
				suite.Equal(5, review.Rating)
				suite.Equal("Excellent product!", review.Comment)
			},
		},
		{
			name:     "Unauthorized",
			body:     &domain.CreateReviewReq{Rating: 5},
			userId:   "",
			setup:    func() {},
			expected: http.StatusUnauthorized,
		},
		{
			name:     "InvalidBody",
			body:     map[string]any{"rating": "five"},
			userId:   "u1",
			setup:    func() {},
			expected: http.StatusBadRequest,
		},
		{
			name:   "Fail",
			body:   &domain.CreateReviewReq{Rating: 5, Comment: "Great"},
			userId: "u1",
			setup: func() {
				suite.mockService.On("CreateReview", mock.Anything, "p1", "u1", &domain.CreateReviewReq{
					Rating: 5, Comment: "Great",
				}).Return(nil, errors.New("already reviewed")).Times(1)
			},
			expected: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContext("POST", "/api/v1/products/p1/reviews", tc.body)
			ctx.Params = gin.Params{{Key: "id", Value: "p1"}}
			if tc.userId != "" {
				ctx.Set("userId", tc.userId)
			}
			tc.setup()
			suite.handler.CreateReview(ctx)
			suite.Equal(tc.expected, writer.Code)
			if tc.checkBody != nil {
				tc.checkBody(writer)
			}
		})
	}
}

func (suite *ReviewHandlerTestSuite) TestUpdateReview() {
	tests := []struct {
		name      string
		body      any
		userId    string
		setup     func()
		expected  int
		checkBody func(writer *httptest.ResponseRecorder)
	}{
		{
			name:   "Success",
			body:   &domain.UpdateReviewReq{Rating: 4, Comment: "Updated comment"},
			userId: "u1",
			setup: func() {
				suite.mockService.On("UpdateReview", mock.Anything, "r1", "u1", &domain.UpdateReviewReq{
					Rating: 4, Comment: "Updated comment",
				}).Return(&model.Review{
					ID: "r1", ProductID: "p1", UserID: "u1", Rating: 4, Comment: "Updated comment",
				}, nil).Times(1)
			},
			expected: http.StatusOK,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				var review domain.Review
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				_ = utils.Copy(&review, &res.Result)
				suite.Equal("r1", review.ID)
				suite.Equal(4, review.Rating)
				suite.Equal("Updated comment", review.Comment)
			},
		},
		{
			name:     "Unauthorized",
			body:     &domain.UpdateReviewReq{Rating: 4},
			userId:   "",
			setup:    func() {},
			expected: http.StatusUnauthorized,
		},
		{
			name:     "InvalidBody",
			body:     map[string]any{"rating": "four"},
			userId:   "u1",
			setup:    func() {},
			expected: http.StatusBadRequest,
		},
		{
			name:   "Fail",
			body:   &domain.UpdateReviewReq{Comment: "Updated"},
			userId: "u1",
			setup: func() {
				suite.mockService.On("UpdateReview", mock.Anything, "r1", "u1", &domain.UpdateReviewReq{
					Comment: "Updated",
				}).Return(nil, errors.New("not found")).Times(1)
			},
			expected: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContext("PUT", "/api/v1/products/p1/reviews/r1", tc.body)
			ctx.Params = gin.Params{{Key: "id", Value: "p1"}, {Key: "reviewId", Value: "r1"}}
			if tc.userId != "" {
				ctx.Set("userId", tc.userId)
			}
			tc.setup()
			suite.handler.UpdateReview(ctx)
			suite.Equal(tc.expected, writer.Code)
			if tc.checkBody != nil {
				tc.checkBody(writer)
			}
		})
	}
}

func (suite *ReviewHandlerTestSuite) TestDeleteReview() {
	tests := []struct {
		name     string
		userId   string
		setup    func()
		expected int
	}{
		{
			name:   "Success",
			userId: "u1",
			setup: func() {
				suite.mockService.On("DeleteReview", mock.Anything, "r1", "u1").Return(nil).Times(1)
			},
			expected: http.StatusOK,
		},
		{
			name:     "Unauthorized",
			userId:   "",
			setup:    func() {},
			expected: http.StatusUnauthorized,
		},
		{
			name:   "Fail",
			userId: "u1",
			setup: func() {
				suite.mockService.On("DeleteReview", mock.Anything, "r1", "u1").Return(errors.New("not found")).Times(1)
			},
			expected: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContext("DELETE", "/api/v1/products/p1/reviews/r1", nil)
			ctx.Params = gin.Params{{Key: "id", Value: "p1"}, {Key: "reviewId", Value: "r1"}}
			if tc.userId != "" {
				ctx.Set("userId", tc.userId)
			}
			tc.setup()
			suite.handler.DeleteReview(ctx)
			suite.Equal(tc.expected, writer.Code)
		})
	}
}
