package apperror

import (
	"errors"

	"github.com/gin-gonic/gin"

	"goshop/pkg/config"
)

// HTTPError writes this AppError as an HTTP JSON error response.
// Usage: apperror.ErrUnauthorized.HTTPError(c)
func (e *AppError) HTTPError(c *gin.Context) {
	errorRes := map[string]interface{}{
		"code":    e.Code,
		"message": e.Message,
	}

	cfg := config.GetConfig()
	if cfg.Environment != config.ProductionEnv {
		errorRes["debug"] = e.Error()
	}

	c.JSON(e.HTTPStatus, gin.H{
		"result": nil,
		"error":  errorRes,
	})
}

// ToHTTPError writes any error as an HTTP JSON error response.
// If the error is an *AppError, its HTTPStatus, Code, and Message are used.
// Otherwise it falls back to the provided status and message.
func ToHTTPError(c *gin.Context, err error, fallbackStatus int, fallbackMessage string) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		appErr.HTTPError(c)
		return
	}

	errorRes := map[string]interface{}{
		"message": fallbackMessage,
	}

	cfg := config.GetConfig()
	if cfg.Environment != config.ProductionEnv {
		errorRes["debug"] = err.Error()
	}

	c.JSON(fallbackStatus, gin.H{
		"result": nil,
		"error":  errorRes,
	})
}
