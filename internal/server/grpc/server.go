package grpc

import (
	"fmt"
	"net"

	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/redis"
	"github.com/quangdangfit/gocommon/validation"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	userGRPC "goshop/internal/user/port/grpc"
	"goshop/pkg/config"
	"goshop/pkg/middleware"
)

type Server struct {
	engine    *grpc.Server
	cfg       *config.Schema
	validator validation.Validation
	db        *gorm.DB
	cache     redis.IRedis
}

func NewServer(validator validation.Validation, db *gorm.DB, cache redis.IRedis) *Server {
	interceptor := middleware.NewAuthInterceptor([]string{})

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcRecovery.UnaryServerInterceptor(),
			interceptor.Unary(),
		),
	)

	return &Server{
		engine:    grpcServer,
		cfg:       config.GetConfig(),
		validator: validator,
		db:        db,
		cache:     cache,
	}
}

func (s Server) Run() error {
	userGRPC.RegisterHandlers(s.engine, s.db, s.validator)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.GrpcPort))
	logger.Info("GRPC server is listening on PORT: ", s.cfg.GrpcPort)
	if err != nil {
		logger.Error("Failed to listen: ", err)
		return err
	}

	// Start grpc server
	err = s.engine.Serve(lis)
	if err != nil {
		logger.Fatal("Failed to serve grpc: ", err)
		return err
	}

	return nil
}
