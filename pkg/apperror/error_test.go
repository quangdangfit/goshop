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
	tests := []struct {
		name    string
		err     *AppError
		wantMsg string
	}{
		{
			name:    "simple error message",
			err:     ErrNotFound,
			wantMsg: "Resource not found",
		},
		{
			name:    "error with wrapped inner error",
			err:     Wrap(ErrNotFound, fmt.Errorf("record not found")),
			wantMsg: "Resource not found: record not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.wantMsg, tc.err.Error())
		})
	}
}

func TestAppError_Unwrap(t *testing.T) {
	tests := []struct {
		name      string
		err       *AppError
		wantInner error
	}{
		{
			name:      "with inner error",
			err:       Wrap(ErrInternal, fmt.Errorf("db error")),
			wantInner: fmt.Errorf("db error"),
		},
		{
			name:      "without inner error",
			err:       ErrNotFound,
			wantInner: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			inner := errors.Unwrap(tc.err)
			if tc.wantInner == nil {
				assert.Nil(t, inner)
			} else {
				assert.Equal(t, tc.wantInner.Error(), inner.Error())
			}
		})
	}
}

func TestErrorsAs(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode string
		wantHTTP int
		wantGRPC codes.Code
	}{
		{
			name:     "wrapped AppError",
			err:      Wrap(ErrNotFound, fmt.Errorf("db error")),
			wantCode: "NOT_FOUND",
			wantHTTP: http.StatusNotFound,
			wantGRPC: codes.NotFound,
		},
		{
			name:     "AppError wrapped in fmt.Errorf",
			err:      fmt.Errorf("operation failed: %w", Wrap(ErrForbidden, fmt.Errorf("not owner"))),
			wantCode: "FORBIDDEN",
			wantHTTP: http.StatusForbidden,
			wantGRPC: codes.PermissionDenied,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var appErr *AppError
			assert.True(t, errors.As(tc.err, &appErr))
			assert.Equal(t, tc.wantCode, appErr.Code)
			assert.Equal(t, tc.wantHTTP, appErr.HTTPStatus)
			assert.Equal(t, tc.wantGRPC, appErr.GRPCCode)
		})
	}
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
