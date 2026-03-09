package grpc

import (
	"github.com/quangdangfit/gocommon/validation"
	"google.golang.org/grpc"

	"goshop/internal/product/repository"
	"goshop/internal/product/service"
	"goshop/pkg/dbs"
	pb "goshop/proto/gen/go/product"
)

func RegisterHandlers(svr *grpc.Server, db dbs.Database, validator validation.Validation) {
	productRepo := repository.NewProductRepository(db)
	productSvc := service.NewProductService(validator, productRepo)
	productHandler := NewProductHandler(productSvc)

	pb.RegisterProductServiceServer(svr, productHandler)
}
