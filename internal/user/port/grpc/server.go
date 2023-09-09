package grpc

import (
	"github.com/quangdangfit/gocommon/validation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"goshop/internal/user/repository"
	"goshop/internal/user/service"
	"goshop/pkg/dbs"
	pb "goshop/proto/gen/go/user"
)

func RegisterHandlers(svr *grpc.Server, db dbs.IDatabase, validator validation.Validation) {
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(validator, userRepo)
	userHandler := NewUserHandler(userSvc)

	pb.RegisterUserServiceServer(svr, userHandler)
	reflection.Register(svr)
}
