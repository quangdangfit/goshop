package http

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/quangdangfit/gocommon/validation"

	"goshop/internal/user/model"
	"goshop/internal/user/service/mocks"
	"goshop/pkg/config"
	dbsMocks "goshop/pkg/dbs/mocks"
)

// fakeOIDCServer is a minimal stand-in for an Authentik instance. It serves
// just enough of the OIDC surface (discovery + JWKS + token endpoint) for
// go-oidc to validate id_tokens we sign with the test RSA key. The auth
// endpoint is never hit since we call /token directly via the OIDC code
// exchange.
type fakeOIDCServer struct {
	*httptest.Server
	key      *rsa.PrivateKey
	kid      string
	clientID string
	// claims is the payload baked into the id_token returned by /token.
	claims map[string]any
}

func newFakeOIDCServer(t *testing.T, clientID string, claims map[string]any) *fakeOIDCServer {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	f := &fakeOIDCServer{key: key, kid: "test-kid", clientID: clientID, claims: claims}

	mux := http.NewServeMux()
	srv := httptest.NewServer(mux)
	f.Server = srv

	mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"issuer":                                srv.URL,
			"authorization_endpoint":                srv.URL + "/authorize",
			"token_endpoint":                        srv.URL + "/token",
			"jwks_uri":                              srv.URL + "/jwks",
			"id_token_signing_alg_values_supported": []string{"RS256"},
		})
	})

	mux.HandleFunc("/jwks", func(w http.ResponseWriter, _ *http.Request) {
		n := base64.RawURLEncoding.EncodeToString(key.N.Bytes())
		e := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(key.E)).Bytes())
		_ = json.NewEncoder(w).Encode(map[string]any{
			"keys": []map[string]any{{"kty": "RSA", "kid": f.kid, "alg": "RS256", "use": "sig", "n": n, "e": e}},
		})
	})

	mux.HandleFunc("/token", func(w http.ResponseWriter, _ *http.Request) {
		idToken := f.signIDToken(t)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token": "fake-access",
			"token_type":   "Bearer",
			"expires_in":   3600,
			"id_token":     idToken,
		})
	})

	t.Cleanup(srv.Close)
	return f
}

// signIDToken builds and signs an RS256 id_token with the configured claims.
func (f *fakeOIDCServer) signIDToken(t *testing.T) string {
	t.Helper()
	claims := jwt.MapClaims{
		"iss": f.URL,
		"aud": f.clientID,
		"exp": time.Now().Add(5 * time.Minute).Unix(),
		"iat": time.Now().Unix(),
	}
	for k, v := range f.claims {
		claims[k] = v
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tok.Header["kid"] = f.kid
	signed, err := tok.SignedString(f.key)
	require.NoError(t, err)
	return signed
}

// OIDCCallbackIntegrationTestSuite runs the full Callback path against a
// real *oidc.Validator pointed at the fake provider — exercising both the
// production constructor (NewOIDCHandler) and the verify/claims/upsert/jwt
// branches that the unit suite skips.
type OIDCCallbackIntegrationTestSuite struct {
	suite.Suite
	svc *mocks.UserService
}

func (s *OIDCCallbackIntegrationTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)
	s.svc = mocks.NewUserService(s.T())
}

func TestOIDCCallbackIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(OIDCCallbackIntegrationTestSuite))
}

func (s *OIDCCallbackIntegrationTestSuite) buildHandler(claims map[string]any) (*OIDCHandler, *fakeOIDCServer) {
	const clientID = "goshop-test-client"
	fake := newFakeOIDCServer(s.T(), clientID, claims)

	cfg := &config.Schema{
		AuthSecret:       "test-secret",
		OIDCIssuer:       fake.URL,
		OIDCClientID:     clientID,
		OIDCClientSecret: "secret",
		OIDCRedirectURL:  "http://localhost/callback",
		FrontendBaseURL:  "https://shop.example.com",
	}

	// NewOIDCHandler is the production constructor; this also covers it.
	h := NewOIDCHandler(s.svc, cfg)
	return h, fake
}

func (s *OIDCCallbackIntegrationTestSuite) doCallback(h *OIDCHandler) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/auth/callback?code=fake-code", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = r
	h.Callback(c)
	return w
}

func (s *OIDCCallbackIntegrationTestSuite) TestCallback_Success_RedirectsToFEWithTokens() {
	sub, email, name := "sub-abc", "user@example.com", "Test User"
	h, _ := s.buildHandler(map[string]any{"sub": sub, "email": email, "name": name})

	s.svc.On("UpsertUserFromOIDC", mock.Anything, sub, email, name).
		Return(&model.User{ID: sub, Email: email, Role: model.UserRoleCustomer}, nil).Once()

	w := s.doCallback(h)

	s.Equal(http.StatusFound, w.Code, "body=%s", w.Body.String())
	loc, err := url.Parse(w.Header().Get("Location"))
	s.Require().NoError(err)
	s.Equal("shop.example.com", loc.Host)
	s.Equal("/auth/callback", loc.Path)
	s.NotEmpty(loc.Query().Get("access_token"))
	s.NotEmpty(loc.Query().Get("refresh_token"))
}

func (s *OIDCCallbackIntegrationTestSuite) TestCallback_EmptyEmailClaim_Returns400() {
	h, _ := s.buildHandler(map[string]any{"sub": "sub-1", "email": ""})

	w := s.doCallback(h)

	s.Equal(http.StatusBadRequest, w.Code, "body=%s", w.Body.String())
}

func (s *OIDCCallbackIntegrationTestSuite) TestCallback_UpsertFailure_Returns500() {
	h, _ := s.buildHandler(map[string]any{"sub": "sub-2", "email": "x@example.com"})
	s.svc.On("UpsertUserFromOIDC", mock.Anything, "sub-2", "x@example.com", "").
		Return(nil, fmt.Errorf("db down")).Once()

	w := s.doCallback(h)

	s.Equal(http.StatusInternalServerError, w.Code, "body=%s", w.Body.String())
}

// Verify branch: the id_token is signed with the wrong key so go-oidc's
// signature check fails inside Callback.
func (s *OIDCCallbackIntegrationTestSuite) TestCallback_VerifyFailure_Returns401() {
	h, fake := s.buildHandler(map[string]any{"sub": "sub-v", "email": "v@example.com"})

	// Re-key the fake server after the handler discovered its JWKS so the
	// next signed token can't be verified.
	wrongKey, err := rsa.GenerateKey(rand.Reader, 2048)
	s.Require().NoError(err)
	fake.key = wrongKey

	w := s.doCallback(h)

	s.Equal(http.StatusUnauthorized, w.Code, "body=%s", w.Body.String())
}

// Claims-unmarshal branch: go-oidc's Verify only parses the standard claims
// (sub/aud/iss/exp/iat), so a non-string `email` slips past Verify and only
// blows up when our handler decodes the payload into its own struct.
func (s *OIDCCallbackIntegrationTestSuite) TestCallback_ClaimsUnmarshalFailure_Returns401() {
	h, _ := s.buildHandler(map[string]any{"sub": "sub-c", "email": 123, "name": "n"})

	w := s.doCallback(h)

	s.Equal(http.StatusUnauthorized, w.Code, "body=%s", w.Body.String())
}

func (s *OIDCCallbackIntegrationTestSuite) TestCallback_InvalidFrontendURL_Returns500() {
	const clientID = "goshop-test-client"
	fake := newFakeOIDCServer(s.T(), clientID, map[string]any{"sub": "sub-3", "email": "y@example.com"})
	cfg := &config.Schema{
		AuthSecret:       "test-secret",
		OIDCIssuer:       fake.URL,
		OIDCClientID:     clientID,
		OIDCClientSecret: "secret",
		OIDCRedirectURL:  "http://localhost/callback",
		// Empty host trips the url.Parse host check at the bottom of Callback.
		FrontendBaseURL: "not-a-url",
	}
	h := NewOIDCHandler(s.svc, cfg)
	s.svc.On("UpsertUserFromOIDC", mock.Anything, "sub-3", "y@example.com", "").
		Return(&model.User{ID: "sub-3", Email: "y@example.com"}, nil).Once()

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/auth/callback?code=c", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = r
	h.Callback(c)

	s.Equal(http.StatusInternalServerError, w.Code, "body=%s", w.Body.String())
	// And the redirect should not have been written.
	s.True(strings.Contains(w.Body.String(), "frontend_base_url") || w.Header().Get("Location") == "")
}

// TestRoutes_OIDCMode_MountsGETLoginAndCallback covers the routes.go OIDC
// branch end-to-end by pointing the global config at a fake provider and
// wiring up the full route group. We assert the OIDC endpoints actually
// exist by hitting them through the router.
func TestRoutes_OIDCMode_MountsGETLoginAndCallback(t *testing.T) {
	const clientID = "goshop-routes-client"
	fake := newFakeOIDCServer(t, clientID, map[string]any{"sub": "irrelevant"})

	// Load config from env so config.GetConfig() reflects OIDC mode for the
	// duration of this test. t.Setenv restores values automatically.
	t.Setenv("auth_mode", "oidc")
	t.Setenv("oidc_issuer", fake.URL)
	t.Setenv("oidc_client_id", clientID)
	t.Setenv("oidc_client_secret", "secret")
	t.Setenv("oidc_redirect_url", "http://localhost/cb")
	t.Setenv("frontend_base_url", "https://shop.example.com")
	t.Setenv("auth_secret", "test-secret")
	_ = config.LoadConfig()
	// Restore the pkg-level config after this test so other tests in the
	// package don't inherit auth_mode=oidc + a dead fake-server URL. We must
	// unset the OIDC env vars first because t.Setenv's own restoration runs
	// after this cleanup (Cleanups fire LIFO and t.Setenv registered first).
	t.Cleanup(func() {
		for _, k := range []string{"auth_mode", "oidc_issuer", "oidc_client_id", "oidc_client_secret", "oidc_redirect_url", "frontend_base_url"} {
			_ = os.Unsetenv(k)
		}
		_ = config.LoadConfig()
	})

	db := dbsMocks.NewDatabase(t)
	engine := gin.New()
	Routes(engine.Group("/api/v1"), db, validation.New())

	// GET /auth/login should now exist (OIDC redirect to fake provider).
	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/login", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	require.Equal(t, http.StatusFound, w.Code, "GET /auth/login should redirect in OIDC mode: %s", w.Body.String())

	loc, err := url.Parse(w.Header().Get("Location"))
	require.NoError(t, err)
	require.Contains(t, loc.Host, strings.TrimPrefix(fake.URL, "http://"), "should redirect to the fake Authentik provider")

	// And GET /auth/callback should be wired (missing code → 400).
	req = httptest.NewRequest(http.MethodGet, "/api/v1/auth/callback", nil)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}
