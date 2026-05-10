package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	orderModel "goshop/internal/order/model"
	orderSvcMocks "goshop/internal/order/service/mocks"
	"goshop/internal/payment/model"
	"goshop/internal/payment/repository"
	"goshop/pkg/payment"
)

type stubProvider struct {
	createFn func(ctx context.Context, p payment.CreateIntentParams) (*payment.Intent, error)
	verifyFn func(payload []byte, sig string) (*payment.Event, error)
}

func (s *stubProvider) CreateIntent(ctx context.Context, p payment.CreateIntentParams) (*payment.Intent, error) {
	return s.createFn(ctx, p)
}
func (s *stubProvider) VerifyWebhook(payload []byte, sig string) (*payment.Event, error) {
	return s.verifyFn(payload, sig)
}

type stubRepo struct {
	getFn      func(ctx context.Context, orderID string) (*model.Payment, error)
	createFn   func(ctx context.Context, p *model.Payment) error
	updateFn   func(ctx context.Context, p *model.Payment) error
	recordFn   func(ctx context.Context, provider, eventID string) error
	createCall int
	updateCall int
}

func (s *stubRepo) GetByOrderID(ctx context.Context, o string) (*model.Payment, error) {
	return s.getFn(ctx, o)
}
func (s *stubRepo) Create(ctx context.Context, p *model.Payment) error {
	s.createCall++
	return s.createFn(ctx, p)
}
func (s *stubRepo) Update(ctx context.Context, p *model.Payment) error {
	s.updateCall++
	return s.updateFn(ctx, p)
}
func (s *stubRepo) RecordProviderEvent(ctx context.Context, provider, eventID string) error {
	return s.recordFn(ctx, provider, eventID)
}

type stubOrderQuery struct {
	getFn func(ctx context.Context, id string) (*orderModel.Order, error)
}

func (s *stubOrderQuery) GetOrderByID(ctx context.Context, id string) (*orderModel.Order, error) {
	return s.getFn(ctx, id)
}

func newOrderSvcMock(t *testing.T) *orderSvcMocks.OrderService {
	return orderSvcMocks.NewOrderService(t)
}

func TestCreateIntent_OrderNotFound(t *testing.T) {
	osvc := newOrderSvcMock(t)
	q := &stubOrderQuery{getFn: func(_ context.Context, _ string) (*orderModel.Order, error) {
		return nil, errors.New("not found")
	}}
	svc := NewPaymentService(&stubProvider{}, &stubRepo{}, q, osvc)
	_, err := svc.CreateIntentForOrder(context.Background(), "o1")
	require.Error(t, err)
}

func TestCreateIntent_OrderNotPending(t *testing.T) {
	osvc := newOrderSvcMock(t)
	q := &stubOrderQuery{getFn: func(_ context.Context, _ string) (*orderModel.Order, error) {
		return &orderModel.Order{ID: "o1", Status: orderModel.OrderStatusPaid}, nil
	}}
	svc := NewPaymentService(&stubProvider{}, &stubRepo{}, q, osvc)
	_, err := svc.CreateIntentForOrder(context.Background(), "o1")
	require.Error(t, err)
}

func TestCreateIntent_ReusesExistingPending(t *testing.T) {
	osvc := newOrderSvcMock(t)
	q := &stubOrderQuery{getFn: func(_ context.Context, _ string) (*orderModel.Order, error) {
		return &orderModel.Order{ID: "o1", Status: orderModel.OrderStatusPendingPayment, FinalPrice: 10}, nil
	}}
	repo := &stubRepo{getFn: func(_ context.Context, _ string) (*model.Payment, error) {
		return &model.Payment{ProviderIntentID: "pi_1", Status: model.PaymentStatusPending, Amount: 1000, Currency: "usd"}, nil
	}}
	svc := NewPaymentService(&stubProvider{}, repo, q, osvc)
	intent, err := svc.CreateIntentForOrder(context.Background(), "o1")
	require.NoError(t, err)
	require.Equal(t, "pi_1", intent.ID)
}

func TestCreateIntent_GetPaymentErrorPropagates(t *testing.T) {
	osvc := newOrderSvcMock(t)
	q := &stubOrderQuery{getFn: func(_ context.Context, _ string) (*orderModel.Order, error) {
		return &orderModel.Order{ID: "o1", Status: orderModel.OrderStatusPendingPayment}, nil
	}}
	repo := &stubRepo{getFn: func(_ context.Context, _ string) (*model.Payment, error) {
		return nil, errors.New("db down")
	}}
	svc := NewPaymentService(&stubProvider{}, repo, q, osvc)
	_, err := svc.CreateIntentForOrder(context.Background(), "o1")
	require.Error(t, err)
}

func TestCreateIntent_CreatesNewWhenNoExisting(t *testing.T) {
	osvc := newOrderSvcMock(t)
	q := &stubOrderQuery{getFn: func(_ context.Context, _ string) (*orderModel.Order, error) {
		return &orderModel.Order{ID: "o1", Status: orderModel.OrderStatusPendingPayment, FinalPrice: 1.5}, nil
	}}
	repo := &stubRepo{
		getFn: func(_ context.Context, _ string) (*model.Payment, error) { return nil, gorm.ErrRecordNotFound },
		createFn: func(_ context.Context, p *model.Payment) error {
			require.Equal(t, "stripe", p.Provider)
			return nil
		},
	}
	prov := &stubProvider{createFn: func(_ context.Context, p payment.CreateIntentParams) (*payment.Intent, error) {
		require.Equal(t, "order_o1", p.IdempotencyKey)
		return &payment.Intent{ID: "pi_new", Amount: p.Amount, Currency: p.Currency}, nil
	}}
	svc := NewPaymentService(prov, repo, q, osvc)
	intent, err := svc.CreateIntentForOrder(context.Background(), "o1")
	require.NoError(t, err)
	require.Equal(t, "pi_new", intent.ID)
}

func TestCreateIntent_ProviderError(t *testing.T) {
	osvc := newOrderSvcMock(t)
	q := &stubOrderQuery{getFn: func(_ context.Context, _ string) (*orderModel.Order, error) {
		return &orderModel.Order{ID: "o1", Status: orderModel.OrderStatusPendingPayment}, nil
	}}
	repo := &stubRepo{getFn: func(_ context.Context, _ string) (*model.Payment, error) { return nil, gorm.ErrRecordNotFound }}
	prov := &stubProvider{createFn: func(_ context.Context, _ payment.CreateIntentParams) (*payment.Intent, error) {
		return nil, errors.New("stripe down")
	}}
	svc := NewPaymentService(prov, repo, q, osvc)
	_, err := svc.CreateIntentForOrder(context.Background(), "o1")
	require.Error(t, err)
}

func TestCreateIntent_RepoCreateError(t *testing.T) {
	osvc := newOrderSvcMock(t)
	q := &stubOrderQuery{getFn: func(_ context.Context, _ string) (*orderModel.Order, error) {
		return &orderModel.Order{ID: "o1", Status: orderModel.OrderStatusPendingPayment}, nil
	}}
	repo := &stubRepo{
		getFn:    func(_ context.Context, _ string) (*model.Payment, error) { return nil, gorm.ErrRecordNotFound },
		createFn: func(_ context.Context, _ *model.Payment) error { return errors.New("disk full") },
	}
	prov := &stubProvider{createFn: func(_ context.Context, _ payment.CreateIntentParams) (*payment.Intent, error) {
		return &payment.Intent{ID: "pi"}, nil
	}}
	svc := NewPaymentService(prov, repo, q, osvc)
	_, err := svc.CreateIntentForOrder(context.Background(), "o1")
	require.Error(t, err)
}

func TestHandleWebhook_VerifyError(t *testing.T) {
	prov := &stubProvider{verifyFn: func(_ []byte, _ string) (*payment.Event, error) {
		return nil, payment.ErrInvalidSignature
	}}
	svc := NewPaymentService(prov, &stubRepo{}, &stubOrderQuery{}, newOrderSvcMock(t))
	require.Error(t, svc.HandleWebhook(context.Background(), nil, "sig"))
}

func TestHandleWebhook_DuplicateIsNoop(t *testing.T) {
	prov := &stubProvider{verifyFn: func(_ []byte, _ string) (*payment.Event, error) {
		return &payment.Event{ID: "evt_1", OrderID: "o1", Type: payment.EventPaymentSucceeded}, nil
	}}
	repo := &stubRepo{recordFn: func(_ context.Context, _, _ string) error { return repository.ErrEventAlreadyProcessed }}
	svc := NewPaymentService(prov, repo, &stubOrderQuery{}, newOrderSvcMock(t))
	require.NoError(t, svc.HandleWebhook(context.Background(), nil, "sig"))
}

func TestHandleWebhook_RecordError(t *testing.T) {
	prov := &stubProvider{verifyFn: func(_ []byte, _ string) (*payment.Event, error) {
		return &payment.Event{ID: "evt_1", OrderID: "o1"}, nil
	}}
	repo := &stubRepo{recordFn: func(_ context.Context, _, _ string) error { return errors.New("db down") }}
	svc := NewPaymentService(prov, repo, &stubOrderQuery{}, newOrderSvcMock(t))
	require.Error(t, svc.HandleWebhook(context.Background(), nil, "sig"))
}

func TestHandleWebhook_MissingOrderID(t *testing.T) {
	prov := &stubProvider{verifyFn: func(_ []byte, _ string) (*payment.Event, error) {
		return &payment.Event{ID: "evt_1"}, nil
	}}
	repo := &stubRepo{recordFn: func(_ context.Context, _, _ string) error { return nil }}
	svc := NewPaymentService(prov, repo, &stubOrderQuery{}, newOrderSvcMock(t))
	require.Error(t, svc.HandleWebhook(context.Background(), nil, "sig"))
}

func TestHandleWebhook_GetByOrderError(t *testing.T) {
	prov := &stubProvider{verifyFn: func(_ []byte, _ string) (*payment.Event, error) {
		return &payment.Event{ID: "evt", OrderID: "o1", Type: payment.EventPaymentSucceeded}, nil
	}}
	repo := &stubRepo{
		recordFn: func(_ context.Context, _, _ string) error { return nil },
		getFn:    func(_ context.Context, _ string) (*model.Payment, error) { return nil, errors.New("gone") },
	}
	svc := NewPaymentService(prov, repo, &stubOrderQuery{}, newOrderSvcMock(t))
	require.Error(t, svc.HandleWebhook(context.Background(), nil, "sig"))
}

func newWebhookSvc(t *testing.T, evt payment.EventType) (PaymentService, *stubRepo, *orderSvcMocks.OrderService) {
	prov := &stubProvider{verifyFn: func(_ []byte, _ string) (*payment.Event, error) {
		return &payment.Event{ID: "evt", OrderID: "o1", Type: evt}, nil
	}}
	repo := &stubRepo{
		recordFn: func(_ context.Context, _, _ string) error { return nil },
		getFn:    func(_ context.Context, _ string) (*model.Payment, error) { return &model.Payment{ID: "p1"}, nil },
		updateFn: func(_ context.Context, _ *model.Payment) error { return nil },
	}
	osvc := newOrderSvcMock(t)
	return NewPaymentService(prov, repo, &stubOrderQuery{}, osvc), repo, osvc
}

func TestHandleWebhook_Succeeded(t *testing.T) {
	svc, repo, osvc := newWebhookSvc(t, payment.EventPaymentSucceeded)
	osvc.On("MarkOrderPaid", mock.Anything, "o1").Return(&orderModel.Order{}, nil).Once()
	require.NoError(t, svc.HandleWebhook(context.Background(), nil, "sig"))
	require.Equal(t, 1, repo.updateCall)
}

func TestHandleWebhook_Succeeded_UpdateError(t *testing.T) {
	svc, repo, _ := newWebhookSvc(t, payment.EventPaymentSucceeded)
	repo.updateFn = func(_ context.Context, _ *model.Payment) error { return errors.New("db") }
	require.Error(t, svc.HandleWebhook(context.Background(), nil, "sig"))
}

func TestHandleWebhook_Succeeded_MarkPaidError(t *testing.T) {
	svc, _, osvc := newWebhookSvc(t, payment.EventPaymentSucceeded)
	osvc.On("MarkOrderPaid", mock.Anything, "o1").Return(nil, errors.New("commit failed")).Once()
	require.Error(t, svc.HandleWebhook(context.Background(), nil, "sig"))
}

func TestHandleWebhook_Failed(t *testing.T) {
	svc, _, osvc := newWebhookSvc(t, payment.EventPaymentFailed)
	osvc.On("UpdateOrderStatus", mock.Anything, "o1", orderModel.OrderStatusPaymentFailed).Return(&orderModel.Order{}, nil).Once()
	require.NoError(t, svc.HandleWebhook(context.Background(), nil, "sig"))
}

func TestHandleWebhook_Failed_UpdateError(t *testing.T) {
	svc, repo, _ := newWebhookSvc(t, payment.EventPaymentFailed)
	repo.updateFn = func(_ context.Context, _ *model.Payment) error { return errors.New("db") }
	require.Error(t, svc.HandleWebhook(context.Background(), nil, "sig"))
}

func TestHandleWebhook_Failed_OrderUpdateError(t *testing.T) {
	svc, _, osvc := newWebhookSvc(t, payment.EventPaymentFailed)
	osvc.On("UpdateOrderStatus", mock.Anything, "o1", orderModel.OrderStatusPaymentFailed).Return(nil, errors.New("nope")).Once()
	require.Error(t, svc.HandleWebhook(context.Background(), nil, "sig"))
}

func TestHandleWebhook_Canceled(t *testing.T) {
	svc, _, osvc := newWebhookSvc(t, payment.EventPaymentCanceled)
	osvc.On("UpdateOrderStatus", mock.Anything, "o1", orderModel.OrderStatusCancelled).Return(&orderModel.Order{}, nil).Once()
	require.NoError(t, svc.HandleWebhook(context.Background(), nil, "sig"))
}

func TestHandleWebhook_Canceled_UpdateError(t *testing.T) {
	svc, repo, _ := newWebhookSvc(t, payment.EventPaymentCanceled)
	repo.updateFn = func(_ context.Context, _ *model.Payment) error { return errors.New("db") }
	require.Error(t, svc.HandleWebhook(context.Background(), nil, "sig"))
}

func TestHandleWebhook_Canceled_OrderUpdateError(t *testing.T) {
	svc, _, osvc := newWebhookSvc(t, payment.EventPaymentCanceled)
	osvc.On("UpdateOrderStatus", mock.Anything, "o1", orderModel.OrderStatusCancelled).Return(nil, errors.New("nope")).Once()
	require.Error(t, svc.HandleWebhook(context.Background(), nil, "sig"))
}

func TestHandleWebhook_Processing(t *testing.T) {
	svc, _, _ := newWebhookSvc(t, payment.EventPaymentProcessing)
	require.NoError(t, svc.HandleWebhook(context.Background(), nil, "sig"))
}

func TestHandleWebhook_Processing_UpdateError(t *testing.T) {
	svc, repo, _ := newWebhookSvc(t, payment.EventPaymentProcessing)
	repo.updateFn = func(_ context.Context, _ *model.Payment) error { return errors.New("db") }
	require.Error(t, svc.HandleWebhook(context.Background(), nil, "sig"))
}

func TestHandleWebhook_RequiresAction(t *testing.T) {
	svc, _, _ := newWebhookSvc(t, payment.EventPaymentRequiresAction)
	require.NoError(t, svc.HandleWebhook(context.Background(), nil, "sig"))
}

func TestHandleWebhook_RequiresAction_UpdateError(t *testing.T) {
	svc, repo, _ := newWebhookSvc(t, payment.EventPaymentRequiresAction)
	repo.updateFn = func(_ context.Context, _ *model.Payment) error { return errors.New("db") }
	require.Error(t, svc.HandleWebhook(context.Background(), nil, "sig"))
}

func TestHandleWebhook_UnknownEventTypeIgnored(t *testing.T) {
	svc, _, _ := newWebhookSvc(t, payment.EventType("unknown.weird"))
	require.NoError(t, svc.HandleWebhook(context.Background(), nil, "sig"))
}
