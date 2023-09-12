package grpc

import (
	"github.com/quangdangfit/gocommon/validation"
	"google.golang.org/grpc"

	"goshop/internal/cart/repository"
	"goshop/internal/cart/service"
	"goshop/pkg/dbs"
	pb "goshop/proto/gen/go/cart"
)

func RegisterHandlers(svr *grpc.Server, db dbs.IDatabase, validator validation.Validation) {
	cartRepo := repository.NewCartRepository(db)
	cartSvc := service.NewCartService(validator, cartRepo)
	cartHandler := NewCartHandler(cartSvc)

	pb.RegisterCartServiceServer(svr, cartHandler)
}
