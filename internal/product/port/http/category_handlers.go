package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"

	"goshop/internal/product/domain"
	"goshop/internal/product/service"
	"goshop/pkg/apperror"
	"goshop/pkg/response"
	"goshop/pkg/utils"
)

type CategoryHandler struct {
	service service.CategoryService
}

func NewCategoryHandler(svc service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: svc}
}

func (h *CategoryHandler) ListCategories(c *gin.Context) {
	categories, err := h.service.ListCategories(c)
	if err != nil {
		logger.Error("Failed to list categories: ", err)
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}
	var res []*domain.Category
	if err := utils.Copy(&res, &categories); err != nil {
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}
	response.JSON(c, http.StatusOK, res)
}

func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		apperror.WrapMessage(apperror.ErrBadRequest, nil, "Missing category ID").HTTPError(c)
		return
	}
	category, err := h.service.GetCategoryByID(c, id)
	if err != nil {
		apperror.Wrap(apperror.ErrNotFound, err).HTTPError(c)
		return
	}
	var res domain.Category
	if err := utils.Copy(&res, category); err != nil {
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}
	response.JSON(c, http.StatusOK, res)
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req domain.CreateCategoryReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		apperror.Wrap(apperror.ErrBadRequest, err).HTTPError(c)
		return
	}
	category, err := h.service.Create(c, &req)
	if err != nil {
		logger.Error("Failed to create category: ", err)
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}
	var res domain.Category
	if err := utils.Copy(&res, category); err != nil {
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}
	response.JSON(c, http.StatusOK, res)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var req domain.UpdateCategoryReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		apperror.Wrap(apperror.ErrBadRequest, err).HTTPError(c)
		return
	}
	category, err := h.service.Update(c, id, &req)
	if err != nil {
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}
	var res domain.Category
	if err := utils.Copy(&res, category); err != nil {
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}
	response.JSON(c, http.StatusOK, res)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(c, id); err != nil {
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}
	response.JSON(c, http.StatusOK, nil)
}
