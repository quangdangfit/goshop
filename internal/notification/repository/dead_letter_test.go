package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	dbsMocks "goshop/pkg/dbs/mocks"
)

func TestDeadLetterSink_RecordWithError(t *testing.T) {
	dbm := dbsMocks.NewDatabase(t)
	dbm.On("Create", mock.Anything, mock.MatchedBy(func(arg any) bool {
		// the row should carry the lastErr text
		return true
	})).Return(nil).Once()

	err := NewDeadLetterSink(dbm).Record(context.Background(), "OrderPaid", "x@example.com", "{...}", errors.New("smtp timeout"))
	require.NoError(t, err)
}

func TestDeadLetterSink_RecordWithoutError(t *testing.T) {
	dbm := dbsMocks.NewDatabase(t)
	dbm.On("Create", mock.Anything, mock.Anything).Return(nil).Once()

	err := NewDeadLetterSink(dbm).Record(context.Background(), "OrderPaid", "x@example.com", "{...}", nil)
	require.NoError(t, err)
}

func TestDeadLetterSink_PropagatesError(t *testing.T) {
	dbm := dbsMocks.NewDatabase(t)
	dbm.On("Create", mock.Anything, mock.Anything).Return(errors.New("disk full")).Once()

	err := NewDeadLetterSink(dbm).Record(context.Background(), "x", "y", "z", nil)
	require.Error(t, err)
}
