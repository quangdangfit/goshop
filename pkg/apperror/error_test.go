package apperror

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func TestAppError_Error(t *testing.T) {
	err := ErrNotFound
	assert.Equal(t, "Resource not found", err.Error())
}

func TestAppError_ErrorWithWrapped(t *testing.T) {
	inner := fmt.Errorf("record not found")
	err := Wrap(ErrNotFound, inner)
	assert.Equal(t, "Resource not found: record not found", err.Error())
}

func TestAppError_Unwrap(t *testing.T) {
	inner := fmt.Errorf("db error")
	err := Wrap(ErrInternal, inner)
	assert.Equal(t, inner, errors.Unwrap(err))
}

func TestAppError_UnwrapNil(t *testing.T) {
	err := ErrNotFound
	assert.Nil(t, errors.Unwrap(err))
}

func TestErrorsAs(t *testing.T) {
	inner := fmt.Errorf("db error")
	err := Wrap(ErrNotFound, inner)

	var appErr *AppError
	assert.True(t, errors.As(err, &appErr))
	assert.Equal(t, "NOT_FOUND", appErr.Code)
	assert.Equal(t, http.StatusNotFound, appErr.HTTPStatus)
	assert.Equal(t, codes.NotFound, appErr.GRPCCode)
}

func TestErrorsAs_WrappedInFmtErrorf(t *testing.T) {
	appErr := Wrap(ErrForbidden, fmt.Errorf("not owner"))
	wrapped := fmt.Errorf("operation failed: %w", appErr)

	var target *AppError
	assert.True(t, errors.As(wrapped, &target))
	assert.Equal(t, "FORBIDDEN", target.Code)
}

func TestNew(t *testing.T) {
	err := New("CUSTOM", "custom message", http.StatusTeapot, codes.DataLoss)
	assert.Equal(t, "CUSTOM", err.Code)
	assert.Equal(t, "custom message", err.Message)
	assert.Equal(t, http.StatusTeapot, err.HTTPStatus)
	assert.Equal(t, codes.DataLoss, err.GRPCCode)
}

func TestWrap_PreservesCodeAndStatus(t *testing.T) {
	inner := fmt.Errorf("underlying")
	err := Wrap(ErrUnauthorized, inner)

	assert.Equal(t, "UNAUTHORIZED", err.Code)
	assert.Equal(t, http.StatusUnauthorized, err.HTTPStatus)
	assert.Equal(t, codes.Unauthenticated, err.GRPCCode)
	assert.Equal(t, "Unauthorized", err.Message)
	assert.Equal(t, inner, err.Err)
}

func TestWrapMessage(t *testing.T) {
	err := WrapMessage(ErrBadRequest, nil, "missing field 'name'")
	assert.Equal(t, "BAD_REQUEST", err.Code)
	assert.Equal(t, "missing field 'name'", err.Message)
	assert.Equal(t, http.StatusBadRequest, err.HTTPStatus)
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		err        *AppError
		code       string
		httpStatus int
		grpcCode   codes.Code
	}{
		{ErrBadRequest, "BAD_REQUEST", 400, codes.InvalidArgument},
		{ErrUnauthorized, "UNAUTHORIZED", 401, codes.Unauthenticated},
		{ErrForbidden, "FORBIDDEN", 403, codes.PermissionDenied},
		{ErrNotFound, "NOT_FOUND", 404, codes.NotFound},
		{ErrConflict, "CONFLICT", 409, codes.AlreadyExists},
		{ErrInternal, "INTERNAL", 500, codes.Internal},
		{ErrInvalidCredentials, "INVALID_CREDENTIALS", 401, codes.Unauthenticated},
		{ErrValidation, "VALIDATION_FAILED", 400, codes.InvalidArgument},
		{ErrInvalidStatus, "INVALID_STATUS", 422, codes.FailedPrecondition},
		{ErrCouponExpired, "COUPON_EXPIRED", 400, codes.FailedPrecondition},
		{ErrCouponMaxUsage, "COUPON_MAX_USAGE", 400, codes.FailedPrecondition},
		{ErrCouponMinOrder, "COUPON_MIN_ORDER", 400, codes.FailedPrecondition},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			assert.Equal(t, tt.code, tt.err.Code)
			assert.Equal(t, tt.httpStatus, tt.err.HTTPStatus)
			assert.Equal(t, tt.grpcCode, tt.err.GRPCCode)
		})
	}
}
