package grpc

import (
	"fmt"
	"net"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	cartGRPC "goshop/internal/cart/port/grpc"
	orderGRPC "goshop/internal/order/port/grpc"
	productGRPC "goshop/internal/product/port/grpc"
	userGRPC "goshop/internal/user/port/grpc"
	"goshop/pkg/config"
	"goshop/pkg/dbs"
	"goshop/pkg/middleware"
	"goshop/pkg/redis"
)

type Server struct {
	engine    *grpc.Server
	cfg       *config.Schema
	validator validation.Validation
	db        dbs.Database
	cache     redis.Redis
}

func NewServer(validator validation.Validation, db dbs.Database, cache redis.Redis) *Server {
	interceptor := middleware.NewAuthInterceptor(config.AuthIgnoreMethods)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
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
	cartGRPC.RegisterHandlers(s.engine, s.db, s.validator)
	productGRPC.RegisterHandlers(s.engine, s.db, s.validator)
	orderGRPC.RegisterHandlers(s.engine, s.db, s.validator)

	reflection.Register(s.engine)

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
