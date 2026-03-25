package apperror

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGRPCStatus_Method(t *testing.T) {
	err := ErrNotFound.GRPCStatus()
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Equal(t, "Resource not found", st.Message())
}

func TestGRPCStatus_MethodOnWrapped(t *testing.T) {
	inner := fmt.Errorf("record not found")
	appErr := Wrap(ErrNotFound, inner)
	err := appErr.GRPCStatus()

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Equal(t, "Resource not found", st.Message())
}

func TestToGRPCStatus_WithAppError(t *testing.T) {
	err := ToGRPCStatus(ErrNotFound)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Equal(t, "Resource not found", st.Message())
}

func TestToGRPCStatus_WithPlainError(t *testing.T) {
	err := ToGRPCStatus(fmt.Errorf("something broke"))
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, "something broke", st.Message())
}

func TestGRPCStatus_WithCustomMessage(t *testing.T) {
	err := WrapMessage(ErrBadRequest, nil, "ID is required").GRPCStatus()
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Equal(t, "ID is required", st.Message())
}
