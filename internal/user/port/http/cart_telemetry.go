package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"

	"goshop/pkg/apperror"
	"goshop/pkg/response"
)

// cartSnapshotReq is the FE-supplied snapshot of the client-side cart. Telemetry only —
// the server NEVER reads it back to construct an order. Used for analytics on abandoned
// carts and cross-device hints.
type cartSnapshotReq struct {
	Items []cartSnapshotItem `json:"items" binding:"required"`
}

type cartSnapshotItem struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
}

// PutCartSnapshot accepts a logged-in user's cart as opaque telemetry. Capped item count
// keeps a noisy/bot client from flooding logs.
func PutCartSnapshot(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		apperror.ErrUnauthorized.HTTPError(c)
		return
	}

	var req cartSnapshotReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apperror.Wrap(apperror.ErrBadRequest, err).HTTPError(c)
		return
	}
	if len(req.Items) > 100 {
		apperror.ErrBadRequest.HTTPError(c)
		return
	}

	logger.Infof("cart_snapshot user=%s items=%d", userID, len(req.Items))
	response.JSON(c, http.StatusNoContent, gin.H{})
}
