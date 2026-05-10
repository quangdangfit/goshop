package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"goshop/internal/payment/model"
	dbsMocks "goshop/pkg/dbs/mocks"
)

func newSQLMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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

func TestPaymentRepo_GetByOrderID_Found(t *testing.T) {
	g, m := newSQLMockDB(t)
	dbm := dbsMocks.NewDatabase(t)
	dbm.On("GetDB").Return(g)

	rows := sqlmock.NewRows([]string{"id", "order_id"}).AddRow("p1", "o1")
	m.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "payments" WHERE order_id = $1`)).
		WithArgs("o1", 1).WillReturnRows(rows)

	p, err := NewPaymentRepository(dbm).GetByOrderID(context.Background(), "o1")
	require.NoError(t, err)
	require.Equal(t, "p1", p.ID)
}

func TestPaymentRepo_GetByOrderID_Error(t *testing.T) {
	g, m := newSQLMockDB(t)
	dbm := dbsMocks.NewDatabase(t)
	dbm.On("GetDB").Return(g)
	m.ExpectQuery(`SELECT \* FROM "payments"`).WillReturnError(errors.New("boom"))

	_, err := NewPaymentRepository(dbm).GetByOrderID(context.Background(), "o1")
	require.Error(t, err)
}

func TestPaymentRepo_Create(t *testing.T) {
	dbm := dbsMocks.NewDatabase(t)
	dbm.On("Create", mock.Anything, mock.Anything).Return(nil).Once()
	require.NoError(t, NewPaymentRepository(dbm).Create(context.Background(), &model.Payment{}))
}

func TestPaymentRepo_Update(t *testing.T) {
	dbm := dbsMocks.NewDatabase(t)
	dbm.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
	require.NoError(t, NewPaymentRepository(dbm).Update(context.Background(), &model.Payment{}))
}

func TestPaymentRepo_RecordProviderEvent_Success(t *testing.T) {
	g, m := newSQLMockDB(t)
	dbm := dbsMocks.NewDatabase(t)
	dbm.On("GetDB").Return(g)

	m.ExpectExec(`INSERT INTO "provider_events"`).WillReturnResult(sqlmock.NewResult(1, 1))

	err := NewPaymentRepository(dbm).RecordProviderEvent(context.Background(), "stripe", "evt_1")
	require.NoError(t, err)
}

func TestPaymentRepo_RecordProviderEvent_DuplicateBecomesSentinel(t *testing.T) {
	g, m := newSQLMockDB(t)
	dbm := dbsMocks.NewDatabase(t)
	dbm.On("GetDB").Return(g)

	m.ExpectExec(`INSERT INTO "provider_events"`).
		WillReturnError(errors.New(`pq: duplicate key value violates unique constraint`))

	err := NewPaymentRepository(dbm).RecordProviderEvent(context.Background(), "stripe", "evt_1")
	require.ErrorIs(t, err, ErrEventAlreadyProcessed)
}

func TestPaymentRepo_RecordProviderEvent_GormDuplicatedKeyBecomesSentinel(t *testing.T) {
	g, m := newSQLMockDB(t)
	dbm := dbsMocks.NewDatabase(t)
	dbm.On("GetDB").Return(g)
	m.ExpectExec(`INSERT INTO "provider_events"`).WillReturnError(gorm.ErrDuplicatedKey)

	err := NewPaymentRepository(dbm).RecordProviderEvent(context.Background(), "stripe", "evt_1")
	require.ErrorIs(t, err, ErrEventAlreadyProcessed)
}

func TestPaymentRepo_RecordProviderEvent_OtherErrorPropagates(t *testing.T) {
	g, m := newSQLMockDB(t)
	dbm := dbsMocks.NewDatabase(t)
	dbm.On("GetDB").Return(g)
	m.ExpectExec(`INSERT INTO "provider_events"`).WillReturnError(errors.New("connection lost"))

	err := NewPaymentRepository(dbm).RecordProviderEvent(context.Background(), "stripe", "evt_1")
	require.Error(t, err)
	require.NotErrorIs(t, err, ErrEventAlreadyProcessed)
}

func TestIsDuplicateKey_NilFalse(t *testing.T) {
	require.False(t, isDuplicateKey(nil))
}

func TestContainsAny_EmptySubReturnsFalse(t *testing.T) {
	require.False(t, containsAny("hello", ""))
}

func TestIndexOf_NotFound(t *testing.T) {
	require.Equal(t, -1, indexOf("abc", "xyz"))
}

func TestIndexOf_Found(t *testing.T) {
	require.Equal(t, 1, indexOf("abc", "bc"))
}
