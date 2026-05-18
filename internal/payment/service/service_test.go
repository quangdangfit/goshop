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

func TestCreateIntent_ExistingRow_ReplaysProviderForFreshClientSecret(t *testing.T) {
	osvc := newOrderSvcMock(t)
	q := &stubOrderQuery{getFn: func(_ context.Context, _ string) (*orderModel.Order, error) {
		return &orderModel.Order{ID: "o1", Status: orderModel.OrderStatusPendingPayment, FinalPrice: 10}, nil
	}}
	repo := &stubRepo{getFn: func(_ context.Context, _ string) (*model.Payment, error) {
		return &model.Payment{ProviderIntentID: "pi_1", Status: model.PaymentStatusPending, Amount: 1000, Currency: "usd"}, nil
	}}
	prov := &stubProvider{createFn: func(_ context.Context, p payment.CreateIntentParams) (*payment.Intent, error) {
		require.Equal(t, "order_o1", p.IdempotencyKey)
		// Stripe's idempotency replay returns the same intent with a fresh client_secret.
		return &payment.Intent{ID: "pi_1", ClientSecret: "pi_1_secret_replay", Amount: 1000, Currency: "usd"}, nil
	}}
	svc := NewPaymentService(prov, repo, q, osvc)
	intent, err := svc.CreateIntentForOrder(context.Background(), "o1")
	require.NoError(t, err)
	require.Equal(t, "pi_1", intent.ID)
	require.Equal(t, "pi_1_secret_replay", intent.ClientSecret, "must return a non-empty client_secret on repeat calls")
	require.Equal(t, 0, repo.createCall, "must not write a second payment row on replay")
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

// TestHandleWebhook_PerEventType covers the per-event-type dispatch matrix in
// HandleWebhook: for each provider event we exercise the happy path, the
// repo.Update failure, and (where applicable) the downstream orderService
// failure. The flow inside the handler is always repo.Update → orderService,
// so an Update failure short-circuits before any orderService call.
func TestHandleWebhook_PerEventType(t *testing.T) {
	const (
		osvcMarkPaid = "MarkOrderPaid"
		osvcUpdate   = "UpdateOrderStatus"
	)
	type osvcCall struct {
		method    string
		newStatus orderModel.OrderStatus
		returnErr bool
	}
	tests := []struct {
		name      string
		event     payment.EventType
		updateErr bool
		osvc      *osvcCall
		wantErr   bool
	}{
		{"succeeded_happy", payment.EventPaymentSucceeded, false, &osvcCall{method: osvcMarkPaid}, false},
		{"succeeded_update_error", payment.EventPaymentSucceeded, true, nil, true},
		{"succeeded_mark_paid_error", payment.EventPaymentSucceeded, false, &osvcCall{method: osvcMarkPaid, returnErr: true}, true},

		{"failed_happy", payment.EventPaymentFailed, false, &osvcCall{method: osvcUpdate, newStatus: orderModel.OrderStatusPaymentFailed}, false},
		{"failed_update_error", payment.EventPaymentFailed, true, nil, true},
		{"failed_order_update_error", payment.EventPaymentFailed, false, &osvcCall{method: osvcUpdate, newStatus: orderModel.OrderStatusPaymentFailed, returnErr: true}, true},

		{"canceled_happy", payment.EventPaymentCanceled, false, &osvcCall{method: osvcUpdate, newStatus: orderModel.OrderStatusCancelled}, false},
		{"canceled_update_error", payment.EventPaymentCanceled, true, nil, true},
		{"canceled_order_update_error", payment.EventPaymentCanceled, false, &osvcCall{method: osvcUpdate, newStatus: orderModel.OrderStatusCancelled, returnErr: true}, true},

		{"processing_happy", payment.EventPaymentProcessing, false, nil, false},
		{"processing_update_error", payment.EventPaymentProcessing, true, nil, true},

		{"requires_action_happy", payment.EventPaymentRequiresAction, false, nil, false},
		{"requires_action_update_error", payment.EventPaymentRequiresAction, true, nil, true},

		{"unknown_event_ignored", payment.EventType("unknown.weird"), false, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, osvc := newWebhookSvc(t, tt.event)
			if tt.updateErr {
				repo.updateFn = func(_ context.Context, _ *model.Payment) error { return errors.New("db") }
			}
			if tt.osvc != nil {
				var (
					retOrder *orderModel.Order
					retErr   error
				)
				if tt.osvc.returnErr {
					retErr = errors.New("downstream order svc failed")
				} else {
					retOrder = &orderModel.Order{}
				}
				switch tt.osvc.method {
				case osvcMarkPaid:
					osvc.On("MarkOrderPaid", mock.Anything, "o1").Return(retOrder, retErr).Once()
				case osvcUpdate:
					osvc.On("UpdateOrderStatus", mock.Anything, "o1", tt.osvc.newStatus).Return(retOrder, retErr).Once()
				}
			}

			err := svc.HandleWebhook(context.Background(), nil, "sig")
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.event == payment.EventPaymentSucceeded {
					require.Equal(t, 1, repo.updateCall)
				}
			}
		})
	}
}
