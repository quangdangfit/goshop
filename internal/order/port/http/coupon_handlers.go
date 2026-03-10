package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"

	"goshop/internal/order/dto"
	"goshop/internal/order/service"
	"goshop/pkg/response"
	"goshop/pkg/utils"
)

type CouponHandler struct {
	service service.CouponService
}

func NewCouponHandler(svc service.CouponService) *CouponHandler {
	return &CouponHandler{service: svc}
}

// CreateCoupon godoc
//
//	@Summary	create coupon
//	@Tags		coupons
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		_	body		dto.CreateCouponReq	true	"Body"
//	@Success	200	{object}	dto.Coupon
//	@Router		/api/v1/coupons [post]
func (h *CouponHandler) CreateCoupon(c *gin.Context) {
	var req dto.CreateCouponReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to get body: ", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	coupon, err := h.service.Create(c, &req)
	if err != nil {
		logger.Error("Failed to create coupon: ", err)
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res dto.Coupon
	utils.Copy(&res, coupon)
	response.JSON(c, http.StatusOK, res)
}

// GetCouponByCode godoc
//
//	@Summary	get coupon by code
//	@Tags		coupons
//	@Produce	json
//	@Param		code	path		string	true	"Coupon code"
//	@Success	200		{object}	dto.Coupon
//	@Router		/api/v1/coupons/{code} [get]
func (h *CouponHandler) GetCouponByCode(c *gin.Context) {
	code := c.Param("code")
	coupon, err := h.service.GetByCode(c, code)
	if err != nil {
		logger.Errorf("Failed to get coupon %s: %s", code, err)
		response.Error(c, http.StatusNotFound, err, "Not found")
		return
	}

	var res dto.Coupon
	utils.Copy(&res, coupon)
	response.JSON(c, http.StatusOK, res)
}
