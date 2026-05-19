package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	goOIDC "github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/stretchr/testify/suite"
	"golang.org/x/oauth2"

	"goshop/internal/user/service/mocks"
	"goshop/pkg/config"
)

// fakeOIDC is a stub oidcClient. Each method returns whatever was set on the
// struct, which keeps tests free of any real network or JWT crypto.
type fakeOIDC struct {
	authURL      string
	exchangeTok  *oauth2.Token
	exchangeErr  error
	verifyTok    *goOIDC.IDToken
	verifyErr    error
	exchangeCode string
	verifyRaw    string
}

func (f *fakeOIDC) AuthCodeURL(string) string { return f.authURL }
func (f *fakeOIDC) Exchange(_ context.Context, code string) (*oauth2.Token, error) {
	f.exchangeCode = code
	return f.exchangeTok, f.exchangeErr
}
func (f *fakeOIDC) Verify(_ context.Context, raw string) (*goOIDC.IDToken, error) {
	f.verifyRaw = raw
	return f.verifyTok, f.verifyErr
}

type OIDCHandlerTestSuite struct {
	suite.Suite
	svc *mocks.UserService
	cfg *config.Schema
}

func (s *OIDCHandlerTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)
	s.svc = mocks.NewUserService(s.T())
	s.cfg = &config.Schema{
		AuthSecret:      "test-secret",
		FrontendBaseURL: "https://shop.example.com",
	}
}

func TestOIDCHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(OIDCHandlerTestSuite))
}

// newRequest builds a Gin context wired to a recorder for handler invocation.
func (s *OIDCHandlerTestSuite) newGet(target string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, target, nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = r
	return c, w
}

// Login
// =================================================================================================

func (s *OIDCHandlerTestSuite) TestLogin_NoProviderHint_RedirectsToAuthCodeURL() {
	const authURL = "https://auth.example.com/authorize?client_id=x&state=goshop-state"
	h := newOIDCHandlerWith(s.svc, s.cfg, &fakeOIDC{authURL: authURL})

	c, w := s.newGet("/auth/login")
	h.Login(c)

	s.Equal(http.StatusFound, w.Code)
	s.Equal(authURL, w.Header().Get("Location"))
}

func (s *OIDCHandlerTestSuite) TestLogin_ValidProvider_AppendsIdpHint() {
	const authURL = "https://auth.example.com/authorize?client_id=x&state=goshop-state"
	h := newOIDCHandlerWith(s.svc, s.cfg, &fakeOIDC{authURL: authURL})

	for _, p := range []string{"google", "facebook", "apple"} {
		c, w := s.newGet("/auth/login?provider=" + p)
		h.Login(c)

		s.Equal(http.StatusFound, w.Code)
		loc, err := url.Parse(w.Header().Get("Location"))
		s.Require().NoError(err)
		s.Equal(p, loc.Query().Get("idp"), "idp hint must match provider %q", p)
		// Pre-existing query params survive.
		s.Equal("x", loc.Query().Get("client_id"))
		s.Equal("goshop-state", loc.Query().Get("state"))
	}
}

func (s *OIDCHandlerTestSuite) TestLogin_UnsupportedProvider_Returns400() {
	h := newOIDCHandlerWith(s.svc, s.cfg, &fakeOIDC{authURL: "https://auth.example.com/authorize"})

	c, w := s.newGet("/auth/login?provider=myspace")
	h.Login(c)

	s.Equal(http.StatusBadRequest, w.Code)
}

// Reaches the url.Parse error branch by feeding Login an AuthCodeURL with
// an unescaped '%' — net/url rejects that as an invalid escape. In
// production this branch is unreachable (oauth2 always produces a valid
// URL), but the defensive check still warrants coverage.
func (s *OIDCHandlerTestSuite) TestLogin_MalformedAuthCodeURL_Returns500() {
	h := newOIDCHandlerWith(s.svc, s.cfg, &fakeOIDC{authURL: "http://%"})

	c, w := s.newGet("/auth/login?provider=google")
	h.Login(c)

	s.Equal(http.StatusInternalServerError, w.Code, "body=%s", w.Body.String())
}

func (s *OIDCHandlerTestSuite) TestLogin_ProviderHint_IsCaseInsensitive() {
	h := newOIDCHandlerWith(s.svc, s.cfg, &fakeOIDC{authURL: "https://auth.example.com/a"})

	c, w := s.newGet("/auth/login?provider=GOOGLE")
	h.Login(c)

	s.Equal(http.StatusFound, w.Code)
	loc, _ := url.Parse(w.Header().Get("Location"))
	s.Equal("google", loc.Query().Get("idp"))
}

// NewOIDCHandler init-failure
// =================================================================================================

// When the OIDC provider can't be discovered (e.g. issuer URL unreachable),
// NewOIDCHandler delegates to onValidatorInitError. Production uses
// logger.Fatalf; here we swap in a capturing hook so the branch is covered
// without killing the test process.
func (s *OIDCHandlerTestSuite) TestNewOIDCHandler_ValidatorInitFailure_InvokesHook() {
	var captured error
	original := onValidatorInitError
	onValidatorInitError = func(err error) { captured = err }
	s.T().Cleanup(func() { onValidatorInitError = original })

	cfg := &config.Schema{
		// 127.0.0.1:1 is reliably refused; provider discovery fails fast.
		OIDCIssuer:       "http://127.0.0.1:1",
		OIDCClientID:     "x",
		OIDCClientSecret: "y",
		OIDCRedirectURL:  "http://localhost/cb",
	}
	h := NewOIDCHandler(s.svc, cfg)

	s.Nil(h, "handler must be nil when validator init fails")
	s.Require().Error(captured)
	s.Contains(captured.Error(), "discover OIDC provider")
}

// Callback
// =================================================================================================

func (s *OIDCHandlerTestSuite) TestCallback_MissingCode_Returns400() {
	h := newOIDCHandlerWith(s.svc, s.cfg, &fakeOIDC{})

	c, w := s.newGet("/auth/callback")
	h.Callback(c)

	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *OIDCHandlerTestSuite) TestCallback_ExchangeFailure_Returns401() {
	h := newOIDCHandlerWith(s.svc, s.cfg, &fakeOIDC{exchangeErr: errors.New("upstream down")})

	c, w := s.newGet("/auth/callback?code=abc")
	h.Callback(c)

	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *OIDCHandlerTestSuite) TestCallback_TokenWithoutIDToken_Returns401() {
	tok := &oauth2.Token{AccessToken: "a"} // no id_token extra
	h := newOIDCHandlerWith(s.svc, s.cfg, &fakeOIDC{exchangeTok: tok})

	c, w := s.newGet("/auth/callback?code=abc")
	h.Callback(c)

	s.Equal(http.StatusUnauthorized, w.Code)
}

// Note: a full success-path test would need a real *oidc.IDToken whose
// Claims() can be unmarshalled — go-oidc does not export a constructor for
// fakes. That branch is covered by the integration test (BE auth flow).
