package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"

	"goshop/internal/order/domain"
	"goshop/internal/order/service"
	"goshop/pkg/apperror"
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
//	@Param		_	body		domain.CreateCouponReq	true	"Body"
//	@Success	200	{object}	domain.Coupon
//	@Router		/api/v1/coupons [post]
func (h *CouponHandler) CreateCoupon(c *gin.Context) {
	var req domain.CreateCouponReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to get body: ", err)
		apperror.Wrap(apperror.ErrBadRequest, err).HTTPError(c)
		return
	}

	coupon, err := h.service.Create(c, &req)
	if err != nil {
		logger.Error("Failed to create coupon: ", err)
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}

	var res domain.Coupon
	utils.Copy(&res, coupon)
	response.JSON(c, http.StatusOK, res)
}

// GetCouponByCode godoc
//
//	@Summary	get coupon by code
//	@Tags		coupons
//	@Produce	json
//	@Param		code	path		string	true	"Coupon code"
//	@Success	200		{object}	domain.Coupon
//	@Router		/api/v1/coupons/{code} [get]
func (h *CouponHandler) GetCouponByCode(c *gin.Context) {
	code := c.Param("code")
	coupon, err := h.service.GetByCode(c, code)
	if err != nil {
		logger.Errorf("Failed to get coupon %s: %s", code, err)
		apperror.Wrap(apperror.ErrNotFound, err).HTTPError(c)
		return
	}

	var res domain.Coupon
	utils.Copy(&res, coupon)
	response.JSON(c, http.StatusOK, res)
}
