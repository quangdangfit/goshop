package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/stretchr/testify/require"

	"goshop/internal/notification/model"
	"goshop/pkg/config"
)

type stubPrefSvc struct {
	listFn func(ctx context.Context, userID string) ([]*model.Preference, error)
	setFn  func(ctx context.Context, userID, eventType, channel string, enabled bool) (*model.Preference, error)
}

func (s *stubPrefSvc) List(ctx context.Context, u string) ([]*model.Preference, error) {
	return s.listFn(ctx, u)
}
func (s *stubPrefSvc) Set(ctx context.Context, u, e, c string, en bool) (*model.Preference, error) {
	return s.setFn(ctx, u, e, c, en)
}

func setupHandler(svc *stubPrefSvc) (*Handler, *gin.Engine) {
	logger.Initialize(config.ProductionEnv)
	gin.SetMode(gin.TestMode)
	h := NewHandler(svc)
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("userId", "u1") })
	r.GET("/me/notification-preferences", h.ListPreferences)
	r.PUT("/me/notification-preferences", h.SetPreference)
	return h, r
}

func TestListPreferences_OK(t *testing.T) {
	svc := &stubPrefSvc{listFn: func(_ context.Context, u string) ([]*model.Preference, error) {
		require.Equal(t, "u1", u)
		return []*model.Preference{{ID: "p1"}}, nil
	}}
	_, r := setupHandler(svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/me/notification-preferences", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestListPreferences_Error(t *testing.T) {
	svc := &stubPrefSvc{listFn: func(_ context.Context, _ string) ([]*model.Preference, error) {
		return nil, errors.New("db down")
	}}
	_, r := setupHandler(svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/me/notification-preferences", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestSetPreference_OK(t *testing.T) {
	called := false
	svc := &stubPrefSvc{setFn: func(_ context.Context, u, e, c string, en bool) (*model.Preference, error) {
		called = true
		require.Equal(t, "u1", u)
		require.Equal(t, "OrderPaid", e)
		require.Equal(t, "email", c)
		require.True(t, en)
		return &model.Preference{ID: "p1", Enabled: true}, nil
	}}
	_, r := setupHandler(svc)

	body, _ := json.Marshal(map[string]any{"event_type": "OrderPaid", "channel": "email", "enabled": true})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/me/notification-preferences", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	require.True(t, called)
}

func TestSetPreference_BadBody(t *testing.T) {
	svc := &stubPrefSvc{}
	_, r := setupHandler(svc)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/me/notification-preferences", bytes.NewReader([]byte("not-json")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSetPreference_ServiceError(t *testing.T) {
	svc := &stubPrefSvc{setFn: func(_ context.Context, _, _, _ string, _ bool) (*model.Preference, error) {
		return nil, errors.New("upsert failed")
	}}
	_, r := setupHandler(svc)
	body, _ := json.Marshal(map[string]any{"event_type": "OrderPaid", "channel": "email", "enabled": false})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/me/notification-preferences", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}
