package grpc

import (
	"github.com/quangdangfit/gocommon/validation"
	"google.golang.org/grpc"

	notificationRepo "goshop/internal/notification/repository"
	notificationSvc "goshop/internal/notification/service"
	"goshop/internal/order/repository"
	"goshop/internal/order/service"
	userRepo "goshop/internal/user/repository"
	"goshop/pkg/config"
	"goshop/pkg/dbs"
	"goshop/pkg/notification"
	pb "goshop/proto/gen/go/order"
)

func RegisterHandlers(svr *grpc.Server, db dbs.Database, validator validation.Validation) {
	oRepo := repository.NewOrderRepository(db)
	pRepo := repository.NewProductRepository(db)
	uRepo := repository.NewUserRepository(db)
	couponRepo := repository.NewCouponRepository(db)
	reservationRepo := repository.NewReservationRepository(db)
	cfg := config.GetConfig()
	couponSvc := service.NewCouponService(validator, couponRepo)
	prefChecker := notificationSvc.NewDBPreferenceChecker(
		notificationSvc.NewUserRepoLookup(userRepo.NewUserRepository(db)),
		notificationRepo.NewPreferenceRepository(db),
	)
	notifier := notification.BuildDefault(notification.Settings{
		SMTPHost:     cfg.SMTPHost,
		SMTPPort:     cfg.SMTPPort,
		SMTPUser:     cfg.SMTPUser,
		SMTPPassword: cfg.SMTPPassword,
		EmailFrom:    cfg.EmailFrom,
		Prefs:        prefChecker,
	})
	orderSvc := service.NewOrderService(validator, db, oRepo, pRepo, uRepo, reservationRepo, couponSvc, notifier)
	orderHandler := NewOrderHandler(orderSvc)

	pb.RegisterOrderServiceServer(svr, orderHandler)
}
