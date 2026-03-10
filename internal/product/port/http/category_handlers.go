package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"

	"goshop/internal/product/dto"
	"goshop/internal/product/service"
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
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}
	var res []*dto.Category
	utils.Copy(&res, &categories)
	response.JSON(c, http.StatusOK, res)
}

func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, errors.New("missing id"), "Missing category ID")
		return
	}
	category, err := h.service.GetCategoryByID(c, id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err, "Not found")
		return
	}
	var res dto.Category
	utils.Copy(&res, category)
	response.JSON(c, http.StatusOK, res)
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}
	category, err := h.service.Create(c, &req)
	if err != nil {
		logger.Error("Failed to create category: ", err)
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}
	var res dto.Category
	utils.Copy(&res, category)
	response.JSON(c, http.StatusOK, res)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateCategoryReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}
	category, err := h.service.Update(c, id, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}
	var res dto.Category
	utils.Copy(&res, category)
	response.JSON(c, http.StatusOK, res)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(c, id); err != nil {
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}
	response.JSON(c, http.StatusOK, nil)
}
