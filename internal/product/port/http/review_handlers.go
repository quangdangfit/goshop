package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"

	"goshop/internal/product/dto"
	"goshop/internal/product/service"
	"goshop/pkg/apperror"
	"goshop/pkg/response"
	"goshop/pkg/utils"
)

type ReviewHandler struct {
	service service.ReviewService
}

func NewReviewHandler(svc service.ReviewService) *ReviewHandler {
	return &ReviewHandler{service: svc}
}

func (h *ReviewHandler) ListReviews(c *gin.Context) {
	productID := c.Param("id")
	page, _ := strconv.ParseInt(c.Query("page"), 10, 64)
	limit, _ := strconv.ParseInt(c.Query("limit"), 10, 64)

	reviews, pagination, err := h.service.ListReviews(c, productID, page, limit)
	if err != nil {
		logger.Error("Failed to list reviews: ", err)
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}
	var res dto.ListReviewRes
	utils.Copy(&res.Reviews, &reviews)
	res.Pagination = pagination
	response.JSON(c, http.StatusOK, res)
}

func (h *ReviewHandler) CreateReview(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		apperror.ErrUnauthorized.HTTPError(c)
		return
	}
	productID := c.Param("id")
	var req dto.CreateReviewReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		apperror.Wrap(apperror.ErrBadRequest, err).HTTPError(c)
		return
	}
	review, err := h.service.CreateReview(c, productID, userID, &req)
	if err != nil {
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}
	var res dto.Review
	utils.Copy(&res, review)
	response.JSON(c, http.StatusOK, res)
}

func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		apperror.ErrUnauthorized.HTTPError(c)
		return
	}
	reviewID := c.Param("reviewId")
	var req dto.UpdateReviewReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		apperror.Wrap(apperror.ErrBadRequest, err).HTTPError(c)
		return
	}
	review, err := h.service.UpdateReview(c, reviewID, userID, &req)
	if err != nil {
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}
	var res dto.Review
	utils.Copy(&res, review)
	response.JSON(c, http.StatusOK, res)
}

func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		apperror.ErrUnauthorized.HTTPError(c)
		return
	}
	reviewID := c.Param("reviewId")
	if err := h.service.DeleteReview(c, reviewID, userID); err != nil {
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}
	response.JSON(c, http.StatusOK, nil)
}
