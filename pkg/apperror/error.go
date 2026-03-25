package apperror

import (
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
)

// AppError is a structured application error that carries an error code,
// a user-facing message, the original error, and transport-level status mappings.
type AppError struct {
	Code       string     `json:"code"`
	Message    string     `json:"message"`
	HTTPStatus int        `json:"-"`
	GRPCCode   codes.Code `json:"-"`
	Err        error      `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Err.Error())
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError with the given code, message, and status mappings.
func New(code string, message string, httpStatus int, grpcCode codes.Code) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		GRPCCode:   grpcCode,
	}
}

// Wrap wraps an existing error with an AppError.
func Wrap(appErr *AppError, err error) *AppError {
	return &AppError{
		Code:       appErr.Code,
		Message:    appErr.Message,
		HTTPStatus: appErr.HTTPStatus,
		GRPCCode:   appErr.GRPCCode,
		Err:        err,
	}
}

// WrapMessage wraps an existing error with an AppError and overrides the message.
func WrapMessage(appErr *AppError, err error, message string) *AppError {
	return &AppError{
		Code:       appErr.Code,
		Message:    message,
		HTTPStatus: appErr.HTTPStatus,
		GRPCCode:   appErr.GRPCCode,
		Err:        err,
	}
}

// Pre-defined application errors.
var (
	ErrBadRequest = New(
		"BAD_REQUEST",
		"Invalid request parameters",
		http.StatusBadRequest,
		codes.InvalidArgument,
	)

	ErrUnauthorized = New(
		"UNAUTHORIZED",
		"Unauthorized",
		http.StatusUnauthorized,
		codes.Unauthenticated,
	)

	ErrForbidden = New(
		"FORBIDDEN",
		"Permission denied",
		http.StatusForbidden,
		codes.PermissionDenied,
	)

	ErrNotFound = New(
		"NOT_FOUND",
		"Resource not found",
		http.StatusNotFound,
		codes.NotFound,
	)

	ErrConflict = New(
		"CONFLICT",
		"Resource already exists",
		http.StatusConflict,
		codes.AlreadyExists,
	)

	ErrInternal = New(
		"INTERNAL",
		"Something went wrong",
		http.StatusInternalServerError,
		codes.Internal,
	)

	ErrInvalidCredentials = New(
		"INVALID_CREDENTIALS",
		"Invalid email or password",
		http.StatusUnauthorized,
		codes.Unauthenticated,
	)

	ErrValidation = New(
		"VALIDATION_FAILED",
		"Validation failed",
		http.StatusBadRequest,
		codes.InvalidArgument,
	)

	ErrInvalidStatus = New(
		"INVALID_STATUS",
		"Invalid status transition",
		http.StatusUnprocessableEntity,
		codes.FailedPrecondition,
	)

	ErrCouponExpired = New(
		"COUPON_EXPIRED",
		"Coupon has expired",
		http.StatusBadRequest,
		codes.FailedPrecondition,
	)

	ErrCouponMaxUsage = New(
		"COUPON_MAX_USAGE",
		"Coupon has reached maximum usage",
		http.StatusBadRequest,
		codes.FailedPrecondition,
	)

	ErrCouponMinOrder = New(
		"COUPON_MIN_ORDER",
		"Order total is below minimum required for this coupon",
		http.StatusBadRequest,
		codes.FailedPrecondition,
	)
)
