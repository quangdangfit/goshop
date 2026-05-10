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

	"goshop/internal/notification/model"
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

func TestPreferenceRepo_ListByUser(t *testing.T) {
	g, m := newSQLMockDB(t)
	dbm := dbsMocks.NewDatabase(t)
	dbm.On("GetDB").Return(g)

	rows := sqlmock.NewRows([]string{"id", "user_id"}).AddRow("p1", "u1")
	m.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "preferences" WHERE user_id = $1`)).
		WithArgs("u1").WillReturnRows(rows)

	got, err := NewPreferenceRepository(dbm).ListByUser(context.Background(), "u1")
	require.NoError(t, err)
	require.Len(t, got, 1)
}

func TestPreferenceRepo_ListByUser_DBError(t *testing.T) {
	g, m := newSQLMockDB(t)
	dbm := dbsMocks.NewDatabase(t)
	dbm.On("GetDB").Return(g)

	m.ExpectQuery(`SELECT \* FROM "preferences"`).WillReturnError(errors.New("boom"))

	_, err := NewPreferenceRepository(dbm).ListByUser(context.Background(), "u1")
	require.Error(t, err)
}

func TestPreferenceRepo_GetFound(t *testing.T) {
	g, m := newSQLMockDB(t)
	dbm := dbsMocks.NewDatabase(t)
	dbm.On("GetDB").Return(g)

	rows := sqlmock.NewRows([]string{"id", "user_id", "event_type", "channel", "enabled"}).
		AddRow("p1", "u1", "OrderPaid", "email", true)
	m.ExpectQuery(`SELECT \* FROM "preferences" WHERE user_id = .+ AND event_type = .+ AND channel = .+`).
		WithArgs("u1", "OrderPaid", "email", 1).WillReturnRows(rows)

	got, err := NewPreferenceRepository(dbm).Get(context.Background(), "u1", "OrderPaid", "email")
	require.NoError(t, err)
	require.NotNil(t, got)
	require.True(t, got.Enabled)
}

func TestPreferenceRepo_GetNotFound(t *testing.T) {
	g, m := newSQLMockDB(t)
	dbm := dbsMocks.NewDatabase(t)
	dbm.On("GetDB").Return(g)

	m.ExpectQuery(`SELECT \* FROM "preferences"`).WillReturnError(gorm.ErrRecordNotFound)

	got, err := NewPreferenceRepository(dbm).Get(context.Background(), "u1", "OrderPaid", "email")
	require.NoError(t, err)
	require.Nil(t, got)
}

func TestPreferenceRepo_GetOtherError(t *testing.T) {
	g, m := newSQLMockDB(t)
	dbm := dbsMocks.NewDatabase(t)
	dbm.On("GetDB").Return(g)

	m.ExpectQuery(`SELECT \* FROM "preferences"`).WillReturnError(errors.New("db down"))

	_, err := NewPreferenceRepository(dbm).Get(context.Background(), "u1", "OrderPaid", "email")
	require.Error(t, err)
}

func TestPreferenceRepo_UpsertCreatesNewRow(t *testing.T) {
	g, m := newSQLMockDB(t)
	dbm := dbsMocks.NewDatabase(t)
	dbm.On("GetDB").Return(g)
	dbm.On("Create", mock.Anything, mock.Anything).Return(nil).Once()

	m.ExpectQuery(`SELECT \* FROM "preferences"`).WillReturnError(gorm.ErrRecordNotFound)

	err := NewPreferenceRepository(dbm).Upsert(context.Background(), &model.Preference{
		UserID: "u1", EventType: "OrderPaid", Channel: "email", Enabled: true,
	})
	require.NoError(t, err)
}

func TestPreferenceRepo_UpsertUpdatesExisting(t *testing.T) {
	g, m := newSQLMockDB(t)
	dbm := dbsMocks.NewDatabase(t)
	dbm.On("GetDB").Return(g)
	dbm.On("Update", mock.Anything, mock.Anything).Return(nil).Once()

	rows := sqlmock.NewRows([]string{"id", "user_id", "event_type", "channel", "enabled"}).
		AddRow("p1", "u1", "OrderPaid", "email", true)
	m.ExpectQuery(`SELECT \* FROM "preferences"`).WillReturnRows(rows)

	err := NewPreferenceRepository(dbm).Upsert(context.Background(), &model.Preference{
		UserID: "u1", EventType: "OrderPaid", Channel: "email", Enabled: false,
	})
	require.NoError(t, err)
}

func TestPreferenceRepo_UpsertGetError(t *testing.T) {
	g, m := newSQLMockDB(t)
	dbm := dbsMocks.NewDatabase(t)
	dbm.On("GetDB").Return(g)

	m.ExpectQuery(`SELECT \* FROM "preferences"`).WillReturnError(errors.New("db down"))

	err := NewPreferenceRepository(dbm).Upsert(context.Background(), &model.Preference{UserID: "u1"})
	require.Error(t, err)
}
