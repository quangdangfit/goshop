package apperror

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCStatus returns this AppError as a gRPC status error.
func (e *AppError) GRPCStatus() error {
	return status.Error(e.GRPCCode, e.Message)
}

// ToGRPCStatus converts any error to a gRPC status error.
// If the error is an *AppError, its GRPCCode and Message are used.
// Otherwise it returns codes.Internal with the original error message.
func ToGRPCStatus(err error) error {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.GRPCStatus()
	}
	return status.Error(codes.Internal, err.Error())
}
