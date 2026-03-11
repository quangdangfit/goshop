package grpc

import (
	"github.com/quangdangfit/gocommon/validation"
	"google.golang.org/grpc"

	cartRepository "goshop/internal/cart/repository"
	"goshop/internal/order/repository"
	"goshop/internal/order/service"
	"goshop/pkg/dbs"
	"goshop/pkg/notification"
	pb "goshop/proto/gen/go/order"
)

func RegisterHandlers(svr *grpc.Server, db dbs.Database, validator validation.Validation) {
	oRepo := repository.NewOrderRepository(db)
	pRepo := repository.NewProductRepository(db)
	uRepo := repository.NewUserRepository(db)
	couponRepo := repository.NewCouponRepository(db)
	cartRepo := cartRepository.NewCartRepository(db)
	couponSvc := service.NewCouponService(validator, couponRepo)
	notifier := notification.NewLoggerNotifier()
	orderSvc := service.NewOrderService(validator, oRepo, pRepo, uRepo, cartRepo, couponSvc, notifier)
	orderHandler := NewOrderHandler(orderSvc)

	pb.RegisterOrderServiceServer(svr, orderHandler)
}
