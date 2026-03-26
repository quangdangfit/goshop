package apperror

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGRPCStatus(t *testing.T) {
	tests := []struct {
		name     string
		err      *AppError
		wantCode codes.Code
		wantMsg  string
	}{
		{
			name:     "simple AppError",
			err:      ErrNotFound,
			wantCode: codes.NotFound,
			wantMsg:  "Resource not found",
		},
		{
			name:     "wrapped AppError",
			err:      Wrap(ErrNotFound, fmt.Errorf("record not found")),
			wantCode: codes.NotFound,
			wantMsg:  "Resource not found",
		},
		{
			name:     "custom message",
			err:      WrapMessage(ErrBadRequest, nil, "ID is required"),
			wantCode: codes.InvalidArgument,
			wantMsg:  "ID is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.err.GRPCStatus()
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tc.wantCode, st.Code())
			assert.Equal(t, tc.wantMsg, st.Message())
		})
	}
}

func TestToGRPCStatus(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode codes.Code
		wantMsg  string
	}{
		{
			name:     "with AppError",
			err:      ErrNotFound,
			wantCode: codes.NotFound,
			wantMsg:  "Resource not found",
		},
		{
			name:     "with plain error",
			err:      fmt.Errorf("something broke"),
			wantCode: codes.Internal,
			wantMsg:  "something broke",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ToGRPCStatus(tc.err)
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tc.wantCode, st.Code())
			assert.Equal(t, tc.wantMsg, st.Message())
		})
	}
}
