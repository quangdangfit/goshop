//go:build integration

// authflow_test.go covers end-to-end authentication + authorization paths
// across the HTTP layer: the stitched register→login→me→change-password
// loop, the admin-only role guard, and the /config/public bootstrap payload
// the FE reads to decide whether to render SSO buttons.

package tests_user

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	userDomain "goshop/internal/user/domain"
	userModel "goshop/internal/user/model"
	"goshop/pkg/jtoken"
)

// Register → login → me → change-password → re-login (new password).
// The whole sequence must succeed using only tokens minted along the way —
// catches regressions where a step silently breaks downstream consumers.
func TestAuthFlow_RegisterLoginMeChangePassword(t *testing.T) {
	email := "flow+" + randSuffix() + "@test.com"

	// Register (server returns only the user record — login mints tokens).
	reg := userDomain.RegisterReq{Email: email, Password: "initial-pw1"}
	w := makeRequest("POST", "/api/v1/auth/register", reg, "")
	require.Equal(t, http.StatusOK, w.Code, "register: %s", w.Body.String())

	// Login to obtain a token pair.
	w = makeRequest("POST", "/api/v1/auth/login", userDomain.LoginReq{Email: email, Password: "initial-pw1"}, "")
	require.Equal(t, http.StatusOK, w.Code, "login: %s", w.Body.String())

	var loginResp struct {
		Result struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		} `json:"result"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &loginResp))
	require.NotEmpty(t, loginResp.Result.AccessToken)
	require.NotEmpty(t, loginResp.Result.RefreshToken)

	// /auth/me with the access token.
	w = makeRequest("GET", "/api/v1/auth/me", nil, loginResp.Result.AccessToken)
	require.Equal(t, http.StatusOK, w.Code, "me: %s", w.Body.String())

	// Change password.
	chg := map[string]string{"password": "initial-pw1", "new_password": "rotated-pw2"}
	w = makeRequest("PUT", "/api/v1/auth/change-password", chg, loginResp.Result.AccessToken)
	require.Equal(t, http.StatusOK, w.Code, "change-password: %s", w.Body.String())

	// Old password no longer works.
	w = makeRequest("POST", "/api/v1/auth/login", userDomain.LoginReq{Email: email, Password: "initial-pw1"}, "")
	assert.NotEqual(t, http.StatusOK, w.Code, "old password should fail after rotation")

	// New password works.
	w = makeRequest("POST", "/api/v1/auth/login", userDomain.LoginReq{Email: email, Password: "rotated-pw2"}, "")
	require.Equal(t, http.StatusOK, w.Code, "login with new password: %s", w.Body.String())
}

// Refresh token issued by login must yield a fresh access token; an access
// token must not be accepted as a refresh token (and vice versa).
func TestAuthFlow_RefreshTokenRoundTrip(t *testing.T) {
	tokens := loginTokens()
	require.NotEmpty(t, tokens["refresh_token"])

	w := makeRequest("POST", "/api/v1/auth/refresh", nil, tokens["refresh_token"])
	require.Equal(t, http.StatusOK, w.Code, "refresh: %s", w.Body.String())

	var resp struct {
		Result struct {
			AccessToken string `json:"access_token"`
		} `json:"result"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.NotEmpty(t, resp.Result.AccessToken)
	// The freshly minted access token must work against a protected route.
	w = makeRequest("GET", "/api/v1/auth/me", nil, resp.Result.AccessToken)
	assert.Equal(t, http.StatusOK, w.Code, "refreshed access token must authenticate: %s", w.Body.String())

	// Cross-use must be rejected: access token sent to /refresh.
	w = makeRequest("POST", "/api/v1/auth/refresh", nil, tokens["access_token"])
	assert.Equal(t, http.StatusUnauthorized, w.Code, "access token must not be usable as refresh token")
}

// AdminOnly: customers are forbidden, admins pass. We exercise the guard via
// POST /api/v1/coupons which is mounted as authMiddleware + adminMiddleware.
func TestAuthZ_AdminOnly_CouponsEndpoint(t *testing.T) {
	ctx := context.Background()

	// Seed an admin and a customer directly (BeforeCreate hashes the password
	// and generates IDs; we then mint tokens for both via jtoken).
	admin := &userModel.User{Email: "admin+" + randSuffix() + "@test.com", Password: "adminpass1", Role: userModel.UserRoleAdmin}
	require.NoError(t, dbTest.Create(ctx, admin))
	customer := &userModel.User{Email: "cust+" + randSuffix() + "@test.com", Password: "custpass1", Role: userModel.UserRoleCustomer}
	require.NoError(t, dbTest.Create(ctx, customer))

	adminTok := jtoken.GenerateAccessToken(map[string]interface{}{
		"id":    admin.ID,
		"email": admin.Email,
		"role":  admin.Role,
	})
	custTok := jtoken.GenerateAccessToken(map[string]interface{}{
		"id":    customer.ID,
		"email": customer.Email,
		"role":  customer.Role,
	})

	coupon := map[string]interface{}{
		"code":             "AUTHZ-" + randSuffix(),
		"discount_percent": 10,
		"max_usage":        100,
	}

	// Customer: 403.
	w := makeRequest("POST", "/api/v1/coupons", coupon, custTok)
	assert.Equal(t, http.StatusForbidden, w.Code, "customer must not create coupons: %s", w.Body.String())

	// Admin: not 401/403 (validation may push it to 400, but auth must let it through).
	w = makeRequest("POST", "/api/v1/coupons", coupon, adminTok)
	assert.NotEqual(t, http.StatusUnauthorized, w.Code, "admin auth rejected: %s", w.Body.String())
	assert.NotEqual(t, http.StatusForbidden, w.Code, "admin role rejected: %s", w.Body.String())
}

// /api/v1/auth/me without an Authorization header must be unauthorized.
func TestAuthN_ProtectedRoute_MissingToken(t *testing.T) {
	w := makeRequest("GET", "/api/v1/auth/me", nil, "")
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// /config/public exposes the auth_mode flag the FE keys off to decide
// whether to render SSO buttons. Default config has auth_mode=jwt.
func TestConfigPublic_ContainsAuthMode(t *testing.T) {
	w := makeRequest("GET", "/api/v1/config/public", nil, "")
	require.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Result map[string]interface{} `json:"result"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	// The key must exist even if the test env didn't run LoadConfig (in which
	// case the value is the zero string). Either an empty string or one of
	// the two valid modes is acceptable here — the contract being tested is
	// that the FE will always see the field.
	mode, ok := resp.Result["auth_mode"].(string)
	require.True(t, ok, "auth_mode missing from /config/public: %v", resp.Result)
	assert.Contains(t, []string{"", "jwt", "oidc"}, mode)
}
