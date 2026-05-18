package repository

import (
	"context"
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	dbsMocks "goshop/pkg/dbs/mocks"
)

func newProductSQLMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, m, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = sqlDB.Close() })
	g, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	require.NoError(t, err)
	return g, m
}

// TestProductReservation exercises the three single-statement UPDATEs that drive the
// inventory state machine. Each method has the same outcome matrix: zero rows affected
// is a sentinel error, a DB error propagates, otherwise success.
func TestProductReservation(t *testing.T) {
	const sqlPattern = `UPDATE "products"`

	type outcome struct {
		rowsAffected int64
		dbErr        error
		wantErrIs    error // assert errors.Is on this; nil means require.NoError
		wantAnyErr   bool  // when wantErrIs is nil but we still expect *some* error
	}
	type call func(repo ProductRepository) error

	reserve := func(repo ProductRepository) error {
		return repo.ReserveStock(context.Background(), "p1", 2)
	}
	commit := func(repo ProductRepository) error {
		return repo.CommitReservation(context.Background(), "p1", 1)
	}
	release := func(repo ProductRepository) error {
		return repo.ReleaseReservation(context.Background(), "p1", 1)
	}

	tests := []struct {
		name    string
		fn      call
		outcome outcome
	}{
		{"reserve_success", reserve, outcome{rowsAffected: 1}},
		{"reserve_insufficient_stock", reserve, outcome{rowsAffected: 0, wantErrIs: ErrInsufficientStock}},
		{"reserve_db_error", reserve, outcome{dbErr: errors.New("boom"), wantAnyErr: true}},

		{"commit_success", commit, outcome{rowsAffected: 1}},
		{"commit_no_rows", commit, outcome{rowsAffected: 0, wantAnyErr: true}},
		{"commit_db_error", commit, outcome{dbErr: errors.New("boom"), wantAnyErr: true}},

		{"release_success", release, outcome{rowsAffected: 1}},
		{"release_no_rows_sentinel", release, outcome{rowsAffected: 0, wantErrIs: ErrReservationAlreadyReleased}},
		{"release_db_error", release, outcome{dbErr: errors.New("boom"), wantAnyErr: true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, m := newProductSQLMockDB(t)
			dbm := dbsMocks.NewDatabase(t)
			dbm.On("GetDB").Return(g)
			if tt.outcome.dbErr != nil {
				m.ExpectExec(sqlPattern).WillReturnError(tt.outcome.dbErr)
			} else {
				m.ExpectExec(sqlPattern).WillReturnResult(sqlmock.NewResult(0, tt.outcome.rowsAffected))
			}

			err := tt.fn(NewProductRepository(dbm))
			switch {
			case tt.outcome.wantErrIs != nil:
				require.ErrorIs(t, err, tt.outcome.wantErrIs)
			case tt.outcome.wantAnyErr:
				require.Error(t, err)
			default:
				require.NoError(t, err)
			}
		})
	}
}
