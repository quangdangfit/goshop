package response

import (
	"github.com/gin-gonic/gin"

	"goshop/pkg/apperror"
)

// Error sends an error response. If err is an *apperror.AppError, the HTTP status
// and user-facing message are derived automatically via err.HTTPError(c).
// Otherwise it falls back to the provided status and message.
func Error(c *gin.Context, status int, err error, message string) {
	apperror.ToHTTPError(c, err, status, message)
}
