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

func TestReserveStock_Success(t *testing.T) {
	g, m := newProductSQLMockDB(t)
	dbm := dbsMocks.NewDatabase(t)
	dbm.On("GetDB").Return(g)
	m.ExpectExec(`UPDATE "products" SET "reserved_quantity"`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, NewProductRepository(dbm).ReserveStock(context.Background(), "p1", 2))
}

func TestReserveStock_InsufficientStock(t *testing.T) {
	g, m := newProductSQLMockDB(t)
	dbm := dbsMocks.NewDatabase(t)
	dbm.On("GetDB").Return(g)
	m.ExpectExec(`UPDATE "products" SET "reserved_quantity"`).
		WillReturnResult(sqlmock.NewResult(0, 0))
	err := NewProductRepository(dbm).ReserveStock(context.Background(), "p1", 2)
	require.ErrorIs(t, err, ErrInsufficientStock)
}

func TestReserveStock_DBError(t *testing.T) {
	g, m := newProductSQLMockDB(t)
	dbm := dbsMocks.NewDatabase(t)
	dbm.On("GetDB").Return(g)
	m.ExpectExec(`UPDATE "products"`).WillReturnError(errors.New("boom"))
	require.Error(t, NewProductRepository(dbm).ReserveStock(context.Background(), "p1", 2))
}

func TestCommitReservation_Success(t *testing.T) {
	g, m := newProductSQLMockDB(t)
	dbm := dbsMocks.NewDatabase(t)
	dbm.On("GetDB").Return(g)
	m.ExpectExec(`UPDATE "products" SET "reserved_quantity"=.+,"stock_quantity"=`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, NewProductRepository(dbm).CommitReservation(context.Background(), "p1", 1))
}

func TestCommitReservation_NoRowError(t *testing.T) {
	g, m := newProductSQLMockDB(t)
	dbm := dbsMocks.NewDatabase(t)
	dbm.On("GetDB").Return(g)
	m.ExpectExec(`UPDATE "products"`).WillReturnResult(sqlmock.NewResult(0, 0))
	require.Error(t, NewProductRepository(dbm).CommitReservation(context.Background(), "p1", 1))
}

func TestCommitReservation_DBError(t *testing.T) {
	g, m := newProductSQLMockDB(t)
	dbm := dbsMocks.NewDatabase(t)
	dbm.On("GetDB").Return(g)
	m.ExpectExec(`UPDATE "products"`).WillReturnError(errors.New("boom"))
	require.Error(t, NewProductRepository(dbm).CommitReservation(context.Background(), "p1", 1))
}

func TestReleaseReservation_Success(t *testing.T) {
	g, m := newProductSQLMockDB(t)
	dbm := dbsMocks.NewDatabase(t)
	dbm.On("GetDB").Return(g)
	m.ExpectExec(`UPDATE "products" SET "reserved_quantity"`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, NewProductRepository(dbm).ReleaseReservation(context.Background(), "p1", 1))
}

func TestReleaseReservation_NoRowReturnsSentinel(t *testing.T) {
	g, m := newProductSQLMockDB(t)
	dbm := dbsMocks.NewDatabase(t)
	dbm.On("GetDB").Return(g)
	m.ExpectExec(`UPDATE "products"`).WillReturnResult(sqlmock.NewResult(0, 0))
	err := NewProductRepository(dbm).ReleaseReservation(context.Background(), "p1", 1)
	require.ErrorIs(t, err, ErrReservationAlreadyReleased)
}

func TestReleaseReservation_DBError(t *testing.T) {
	g, m := newProductSQLMockDB(t)
	dbm := dbsMocks.NewDatabase(t)
	dbm.On("GetDB").Return(g)
	m.ExpectExec(`UPDATE "products"`).WillReturnError(errors.New("boom"))
	require.Error(t, NewProductRepository(dbm).ReleaseReservation(context.Background(), "p1", 1))
}
