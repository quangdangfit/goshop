package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"goshop/internal/order/model"
	dbsMocks "goshop/pkg/dbs/mocks"
)

func newReservationSQLMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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

func TestReservationRepo_CreateMany_Empty(t *testing.T) {
	dbm := dbsMocks.NewDatabase(t)
	require.NoError(t, NewReservationRepository(dbm).CreateMany(context.Background(), nil))
}

func TestReservationRepo_CreateMany(t *testing.T) {
	dbm := dbsMocks.NewDatabase(t)
	dbm.On("CreateInBatches", mock.Anything, mock.Anything, 2).Return(nil).Once()
	err := NewReservationRepository(dbm).CreateMany(context.Background(), []*model.StockReservation{{}, {}})
	require.NoError(t, err)
}

func TestReservationRepo_FindActiveByOrderID(t *testing.T) {
	g, m := newReservationSQLMockDB(t)
	dbm := dbsMocks.NewDatabase(t)
	dbm.On("GetDB").Return(g)

	rows := sqlmock.NewRows([]string{"id"}).AddRow("r1")
	m.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "stock_reservations" WHERE order_id = $1 AND status = $2`)).
		WithArgs("o1", string(model.ReservationStatusActive)).WillReturnRows(rows)

	got, err := NewReservationRepository(dbm).FindActiveByOrderID(context.Background(), "o1")
	require.NoError(t, err)
	require.Len(t, got, 1)
}

func TestReservationRepo_FindActiveByOrderID_Error(t *testing.T) {
	g, m := newReservationSQLMockDB(t)
	dbm := dbsMocks.NewDatabase(t)
	dbm.On("GetDB").Return(g)
	m.ExpectQuery(`SELECT \* FROM "stock_reservations"`).WillReturnError(errors.New("boom"))
	_, err := NewReservationRepository(dbm).FindActiveByOrderID(context.Background(), "o1")
	require.Error(t, err)
}

func TestReservationRepo_FindExpired(t *testing.T) {
	g, m := newReservationSQLMockDB(t)
	dbm := dbsMocks.NewDatabase(t)
	dbm.On("GetDB").Return(g)

	rows := sqlmock.NewRows([]string{"id"}).AddRow("r1")
	m.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "stock_reservations" WHERE status = $1 AND expires_at < $2 LIMIT $3`)).
		WithArgs(string(model.ReservationStatusActive), sqlmock.AnyArg(), 5).
		WillReturnRows(rows)

	got, err := NewReservationRepository(dbm).FindExpired(context.Background(), time.Now(), 5)
	require.NoError(t, err)
	require.Len(t, got, 1)
}

func TestReservationRepo_FindExpired_Error(t *testing.T) {
	g, m := newReservationSQLMockDB(t)
	dbm := dbsMocks.NewDatabase(t)
	dbm.On("GetDB").Return(g)
	m.ExpectQuery(`SELECT \* FROM "stock_reservations"`).WillReturnError(errors.New("boom"))
	_, err := NewReservationRepository(dbm).FindExpired(context.Background(), time.Now(), 5)
	require.Error(t, err)
}

func TestReservationRepo_UpdateStatus_Empty(t *testing.T) {
	dbm := dbsMocks.NewDatabase(t)
	require.NoError(t, NewReservationRepository(dbm).UpdateStatus(context.Background(), nil, model.ReservationStatusReleased))
}

func TestReservationRepo_UpdateStatus(t *testing.T) {
	g, m := newReservationSQLMockDB(t)
	dbm := dbsMocks.NewDatabase(t)
	dbm.On("GetDB").Return(g)

	m.ExpectExec(`UPDATE "stock_reservations" SET`).
		WillReturnResult(sqlmock.NewResult(0, 2))

	err := NewReservationRepository(dbm).UpdateStatus(context.Background(), []string{"r1", "r2"}, model.ReservationStatusCommitted)
	require.NoError(t, err)
}
