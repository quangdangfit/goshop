package grpc

import (
	"testing"

	"github.com/quangdangfit/gocommon/validation"
	goGRPC "google.golang.org/grpc"

	"goshop/pkg/dbs/mocks"
)

func TestRegisterHandlers(t *testing.T) {
	mockDB := mocks.NewIDatabase(t)
	RegisterHandlers(goGRPC.NewServer(), mockDB, validation.New())
}
