package grpc

import (
	"testing"

	"github.com/quangdangfit/gocommon/validation"
	"github.com/stretchr/testify/mock"
	goGRPC "google.golang.org/grpc"

	"goshop/pkg/dbs/mocks"
)

func TestRegisterHandlers(t *testing.T) {
	mockDB := mocks.NewIDatabase(t)
	mockDB.On("AutoMigrate", mock.Anything).Return(nil).Times(1)
	RegisterHandlers(goGRPC.NewServer(), mockDB, validation.New())
}
