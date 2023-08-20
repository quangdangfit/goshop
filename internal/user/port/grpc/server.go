package grpc

import (
	"github.com/quangdangfit/gocommon/validation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"

	"goshop/internal/user/repository"
	"goshop/internal/user/service"
	pb "goshop/proto/gen/go/user"
)

func RegisterHandlers(svr *grpc.Server, db *gorm.DB, validator validation.Validation) {
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(validator, userRepo)
	userHandler := NewUserHandler(userSvc)

	pb.RegisterUserServiceServer(svr, userHandler)
	reflection.Register(svr)
}
