# Authentik Integration Plan — GoShop

## 1. Tổng quan

Mục tiêu: Tích hợp **Authentik** (open-source identity provider) làm OIDC/OAuth2 provider cho flow user identity của GoShop, thay thế dần custom JWT + bcrypt hiện tại.

---

## 2. Phân tích trạng thái hiện tại

### 2.1 Auth flow hiện tại

| Layer | File | Mô tả |
|-------|------|-------|
| Domain | [`internal/user/domain/user.go`](internal/user/domain/user.go) | `LoginReq/RegisterReq/LoginRes/RefreshTokenReq/ChangePasswordReq` |
| Model | [`internal/user/model/user.go`](internal/user/model/user.go) | `User{ID, Email, Password, Role}` — bcrypt hash trong `BeforeCreate` |
| Service | [`internal/user/service/user.go`](internal/user/service/user.go) | `Login/Register/GetUserByID/RefreshToken/ChangePassword` |
| Repository | [`internal/user/repository/user.go`](internal/user/repository/user.go) | `Create/Update/GetUserByID/GetUserByEmail` |
| HTTP Handler | [`internal/user/port/http/handlers.go`](internal/user/port/http/handlers.go) | `Login/Register/GetMe/RefreshToken/ChangePassword` |
| HTTP Routes | [`internal/user/port/http/routes.go`](internal/user/port/http/routes.go) | `/auth/register`, `/auth/login`, `/auth/refresh`, `/auth/me`, `/auth/change-password` |
| gRPC Handler | [`internal/user/port/grpc/handlers.go`](internal/user/port/grpc/handlers.go) | `Login/Register/GetMe/RefreshToken/ChangePassword` |
| JWT | [`pkg/jtoken/jwt.go`](pkg/jtoken/jwt.go) | HS256, access 5h, refresh 30d |
| HTTP Middleware | [`pkg/middleware/auth.go`](pkg/middleware/auth.go) | `JWTAuth()` / `JWTRefresh()` — đọc `Authorization: Bearer <token>` |
| gRPC Interceptor | [`pkg/middleware/auth_interceptor.go`](pkg/middleware/auth_interceptor.go) | Parse JWT từ gRPC metadata `token` |
| Config | [`pkg/config/config.go`](pkg/config/config.go) | `auth_secret` env var |
| Admin guard | [`pkg/middleware/admin.go`](pkg/middleware/admin.go) | Kiểm tra `role == "admin"` |

### 2.2 Điểm cần thay đổi

- Xóa dependency vào `bcrypt` và `auth_secret` (sau khi migrate xong)
- Thay thế JWT middleware bằng OIDC token introspection / JWT validation với JWKS
- Thêm flow OIDC: `login → Authentik → callback → tạo/update user local → issue session`
- Giữ lại `role` trong DB để `AdminOnly` middleware hoạt động

---

## 3. Kiến trúc Authentik

### 3.1 Flow OIDC Authorization Code + PKCE (khuyến nghị cho SPA/Mobile)

```
┌──────────┐     ┌──────────────┐     ┌──────────────┐     ┌──────────┐
│  Browser │────▶│  GoShop API  │────▶│   Authentik  │────▶│  User   │
│  (FE)    │     │  /auth/login  │     │  /application/│     │         │
└──────────┘     └──────────────┘     │  oidc/authorize│     └──────────┘
       ▲                  │            └──────────────┘
       │                  │                   │
       │                  │    callback với code
       │                  │◀──────────────────┘
       │                  │
       │     /auth/callback│
       │◀─────────────────┘
       │
       │  Set-Cookie: session (httpOnly)
       │
```

### 3.2 Flow OIDC Client Credentials (cho service-to-service / gRPC)

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│  GoShop gRPC │────▶│  GoShop API  │────▶│   Authentik  │
│  (internal)  │     │  /oidc/token │     │  /oauth2/token│
└──────────────┘     └──────────────┘     └──────────────┘
```

### 3.3 Vai trò Authentik trong kiến trúc

| Vai trò | Mô tả |
|---------|-------|
| **Identity Provider** | Xác thực user (email/password, SSO, MFA, TOTP, WebAuthn) |
| **OAuth2 / OIDC** | Cấp `id_token` + `access_token` + `refresh_token` |
| **User provisioning** | Đồng bộ user vào DB local qua API hoặc SCIM |
| **Session management** | Quản lý logout, session revocation |

---

## 4. Plan chi tiết (phía GoShop)

> **Lưu ý:** Authentik đã được deploy lên K8s bởi repo khác. Phần này chỉ liệt kê những gì cần làm trong repo GoShop.

### Prerequisites từ Authentik (đã có sẵn)

Sau khi Authentik đã chạy, cần lấy các giá trị này từ Admin UI (Applications → goshop-api → Show advanced settings):

| Config key | Mô tả | Ví dụ |
|-----------|-------|-------|
| `OIDC_ISSUER` | Issuer URL của OIDC app | `https://auth.cunghoclaptrinh.online/application/o/goshop-api/` |
| `OIDC_CLIENT_ID` | Client ID | `goshop-api-client` |
| `OIDC_CLIENT_SECRET` | Client secret | (từ Authentik UI) |
| `OIDC_REDIRECT_URL` | Redirect URL sau khi login | `https://goshop.cunghoclaptrinh.online/api/v1/auth/callback` |
| `OIDC_JWKS_URL` | JWKS endpoint để verify token | `https://auth.cunghoclaptrinh.online/application/o/goshop-api/jwks/` |
| `OIDC_SCOPES` | Scopes yêu cầu | `openid,email,profile` |

### Phase 1 — Config & dependencies

#### 1.1 Thêm config OIDC

File: [`pkg/config/config.go`](pkg/config/config.go)

```go
// AuthMode xác định cơ chế xác thực được sử dụng
type AuthMode string

const (
    AuthModeJWT  AuthMode = "jwt"   // JWT truyền thống (bcrypt + HS256)
    AuthModeOIDC AuthMode = "oidc"  // OIDC qua Authentik
)

// Thêm vào Schema struct
AuthMode              AuthMode `env:"auth_mode" envDefault:"jwt"`
OIDCIssuer            string   `env:"oidc_issuer"`
OIDCClientID          string   `env:"oidc_client_id"`
OIDCClientSecret      string   `env:"oidc_client_secret"`
OIDCRedirectURL       string   `env:"oidc_redirect_url"`
OIDCJWKSURL           string   `env:"oidc_jwks_url"`
OIDCScopes            string   `env:"oidc_scopes" envDefault:"openid,email,profile"`
```

File: [`config.sample.yaml`](config.sample.yaml)

```yaml
# Auth mode: "jwt" (mặc định) hoặc "oidc"
auth_mode: jwt

# OIDC / Authentik — chỉ cần khi auth_mode = oidc
oidc_issuer: https://auth.cunghoclaptrinh.online/application/o/goshop-api/
oidc_client_id: goshop-api-client
oidc_client_secret: ######
oidc_redirect_url: https://goshop.cunghoclaptrinh.online/api/v1/auth/callback
oidc_jwks_url: https://auth.cunghoclaptrinh.online/application/o/goshop-api/jwks/
oidc_scopes: openid,email,profile
```

#### 1.2 Thêm dependencies

```bash
go get github.com/coreos/go-oidc/v3
```

---

### Phase 1 — Cấu hình & thư viện

#### 1.1 Thêm config OIDC

File: [`pkg/config/config.go`](pkg/config/config.go)

```go
// Thêm vào Schema struct
OIDCIssuer            string `env:"oidc_issuer"`
OIDCClientID          string `env:"oidc_client_id"`
OIDCClientSecret      string `env:"oidc_client_secret"`
OIDCRedirectURL       string `env:"oidc_redirect_url"`
OIDCJWKSURL           string `env:"oidc_jwks_url"`
OIDCScopes            string `env:"oidc_scopes" envDefault:"openid,email,profile"`
```

File: [`config.sample.yaml`](config.sample.yaml)

```yaml
# OIDC / Authentik
oidc_issuer: http://localhost:9000/application/o/goshop-api/
oidc_client_id: goshop-api-client
oidc_client_secret: ######
oidc_redirect_url: http://localhost:8888/api/v1/auth/callback
oidc_jwks_url: http://localhost:9000/application/o/goshop-api/jwks/
oidc_scopes: openid,email,profile
```

#### 1.2 Thêm dependencies

```bash
go get github.com/coreos/go-oidc/v3
go get github.com/golang-jwt/jwt/v5   # đã có, giữ lại cho JWKS
```

---

### Phase 2 — OIDC Middleware

#### 2.1 JWKS Token Validator

File mới: [`pkg/oidc/validator.go`](pkg/oidc/validator.go)

```go
package oidc

import (
    "context"
    "crypto/rsa"
    "encoding/json"
    "fmt"
    "net/http"
    "sync"
    "time"

    "github.com/coreos/go-oidc/v3/oidc"
    "github.com/golang-jwt/jwt/v5"
    "github.com/quangdangfit/gocommon/logger"
    "golang.org/x/oauth2"
)

type Validator struct {
    provider *oidc.Provider
    verifier *oidc.IDTokenVerifier
    oauth2Config *oauth2.Config
    keyFunc jwt.Keyfunc
    mu sync.RWMutex
}

func NewValidator(issuer, clientID, clientSecret, redirectURL, jwksURL string) (*Validator, error) {
    ctx := context.Background()
    provider, err := oidc.NewProvider(ctx, issuer)
    if err != nil {
        return nil, fmt.Errorf("failed to get provider: %w", err)
    }

    verifier := provider.Verifier(&oidc.Config{
        ClientID: clientID,
    })

    oauth2Config := &oauth2.Config{
        ClientID:     clientID,
        ClientSecret: clientSecret,
        RedirectURL:  redirectURL,
        Endpoint:     provider.Endpoint(),
        Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
    }

    // JWKS key func cho JWT validation
    keyFunc := func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return getPublicKey(ctx, jwksURL, token)
    }

    return &Validator{
        provider:    provider,
        verifier:    verifier,
        oauth2Config: oauth2Config,
        keyFunc:     keyFunc,
    }, nil
}

// getPublicKey fetches and caches JWKS keys
var jwksCache = struct {
    keys map[string]*rsa.PublicKey
    exp  time.Time
    mu   sync.RWMutex
}{}

func getPublicKey(ctx context.Context, jwksURL string, token *jwt.Token) (interface{}, error) {
    kid, ok := token.Header["kid"].(string)
    if !ok {
        return nil, fmt.Errorf("missing kid in token header")
    }

    jwksCache.mu.RLock()
    if key, found := jwksCache.keys[kid]; found && time.Now().Before(jwksCache.exp) {
        jwksCache.mu.RUnlock()
        return key, nil
    }
    jwksCache.mu.RUnlock()

    // Fetch JWKS
    req, err := http.NewRequestWithContext(ctx, "GET", jwksURL, nil)
    if err != nil {
        return nil, err
    }
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var jwks struct {
        Keys []json.RawMessage `json:"keys"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
        return nil, err
    }

    jwksCache.mu.Lock()
    jwksCache.keys = make(map[string]*rsa.PublicKey)
    for _, k := range jwks.Keys {
        var key struct {
            Kid string `json:"kid"`
            Kty string `json:"kty"`
            N   string `json:"n"`
            E   string `json:"e"`
        }
        if err := json.Unmarshal(k, &key); err != nil {
            continue
        }
        if key.Kty != "RSA" {
            continue
        }
        pubKey, err := jwt.ParseRSAPublicKeyFromPEM(
            []byte(fmt.Sprintf("-----BEGIN PUBLIC KEY-----\n%s\n-----END PUBLIC KEY-----", key.N)),
        )
        // reconstruct properly
        // ...
    }
    jwksCache.exp = time.Now().Add(1 * time.Hour)
    jwksCache.mu.Unlock()

    return jwksCache.keys[kid], nil
}
```

#### 2.2 HTTP Middleware OIDC

File mới: [`pkg/middleware/oidc.go`](pkg/middleware/oidc.go)

```go
package middleware

import (
    "context"
    "net/http"

    "github.com/gin-gonic/gin"
    "goshop/pkg/oidc"
)

// OIDCAuth validates Bearer token từ Authentik và inject user info vào context
func OIDCAuth(validator *oidc.Validator) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            // Fallback sang JWT cũ nếu cần (dual-auth mode)
            // hoặc trả 401
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
            return
        }

        ctx := c.Request.Context()
        rawIDToken, err := validator.Verify(ctx, authHeader)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            return
        }

        // Extract claims
        var claims struct {
            Sub   string `json:"sub"`
            Email string `json:"email"`
            Name  string `json:"name"`
        }
        if err := rawIDToken.Claims(&claims); err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
            return
        }

        // Inject vào context — tương thích với code hiện tại dùng "userId"
        c.Set("userId", claims.Sub)
        c.Set("email", claims.Email)
        c.Set("name", claims.Name)
        c.Next()
    }
}
```

#### 2.3 gRPC Interceptor OIDC

File mới: [`pkg/middleware/oidc_interceptor.go`](pkg/middleware/oidc_interceptor.go)

```go
package middleware

import (
    "context"
    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"
    "goshop/pkg/oidc"
)

type OIDCInterceptor struct {
    validator *oidc.Validator
    ignored   []string
}

func NewOIDCInterceptor(validator *oidc.Validator, ignored []string) *OIDCInterceptor {
    return &OIDCInterceptor{validator: validator, ignored: ignored}
}

func (i *OIDCInterceptor) Unary() grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        for _, m := range i.ignored {
            if info.FullMethod == m {
                return handler(ctx, req)
            }
        }
        md, ok := metadata.FromIncomingContext(ctx)
        if !ok {
            return nil, errUnauthorized
        }
        tokens := md.Get("authorization")
        if len(tokens) == 0 {
            return nil, errUnauthorized
        }
        rawIDToken, err := i.validator.Verify(ctx, tokens[0])
        if err != nil {
            return nil, errUnauthorized
        }
        var claims struct{ Sub string }
        if err := rawIDToken.Claims(&claims); err != nil {
            return nil, errUnauthorized
        }
        ctx = context.WithValue(ctx, "userId", claims.Sub)
        return handler(ctx, req)
    }
}
```

---

### Phase 3 — Auth Handlers mới

#### 3.1 Login redirect → Authentik

File mới: [`internal/user/port/http/oidc_handlers.go`](internal/user/port/http/oidc_handlers.go)

```go
package http

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "goshop/pkg/oidc"
)

// Login redirects user đến Authentik authorization endpoint
func (h *OIDCHandler) Login(c *gin.Context) {
    url := h.validator.AuthCodeURL("goshop-state")
    c.Redirect(http.StatusFound, url)
}

// Callback nhận authorization code từ Authentik, đổi lấy token, tạo/update user local
func (h *OIDCHandler) Callback(c *gin.Context) {
    code := c.Query("code")
    if code == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "missing code"})
        return
    }

    ctx := c.Request.Context()
    oauth2Token, err := h.validator.Exchange(ctx, code)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to exchange token"})
        return
    }

    rawIDToken, ok := oauth2Token.Extra("id_token").(string)
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "no id_token"})
        return
    }

    // Verify và extract claims
    idToken, err := h.validator.Verify(ctx, rawIDToken)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid id_token"})
        return
    }

    var claims struct {
        Sub   string `json:"sub"`
        Email string `json:"email"`
        Name  string `json:"name"`
    }
    if err := idToken.Claims(&claims); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
        return
    }

    // Upsert user vào DB local
    user, err := h.service.UpsertUserFromOIDC(ctx, claims.Sub, claims.Email, claims.Name)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sync user"})
        return
    }

    // Issue JWT session token (hoặc dùng OIDC token trực tiếp)
    tokenData := map[string]interface{}{
        "id":    user.ID,
        "email": user.Email,
        "role":  user.Role,
    }
    accessToken := jtoken.GenerateAccessToken(tokenData)
    refreshToken := jtoken.GenerateRefreshToken(tokenData)

    // Redirect về FE với token (hoặc set httpOnly cookie)
    redirectURL := h.cfg.OIDCRedirectURL + "?token=" + accessToken
    c.Redirect(http.StatusFound, redirectURL)
}
```

#### 3.2 User Service — Upsert từ OIDC

File: [`internal/user/service/user.go`](internal/user/service/user.go) — thêm method

```go
type UserService interface {
    // ... existing methods ...
    UpsertUserFromOIDC(ctx context.Context, subject, email, name string) (*model.User, error)
}

func (s *userService) UpsertUserFromOIDC(ctx context.Context, subject, email, name string) (*model.User, error) {
    // 1. Tìm user theo email
    user, err := s.repo.GetUserByEmail(ctx, email)
    if err == nil {
        // User tồn tại — cập nhật subject nếu chưa có
        if user.ID != subject {
            // Có thể lưu subject vào trường riêng hoặc dùng subject làm ID
        }
        return user, nil
    }

    // 2. Tạo user mới — không có password (để Authentik quản lý)
    newUser := &model.User{
        ID:     subject, // Dùng OIDC sub làm primary key
        Email:  email,
        Role:   model.UserRoleCustomer,
    }
    if err := s.repo.Create(ctx, newUser); err != nil {
        return nil, err
    }
    return newUser, nil
}
```

#### 3.3 Cập nhật User Model

File: [`internal/user/model/user.go`](internal/user/model/user.go)

```go
type User struct {
    ID        string     `json:"id" gorm:"unique;not null;index;primary_key"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
    DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
    Email     string     `json:"email" gorm:"unique;not null;index:idx_user_email"`
    Password  string     `json:"password"` // nullable — empty cho OIDC users
    Role      UserRole   `json:"role"`
}
```

> **Lưu ý**: `Password` trở thành nullable. Users đăng ký qua OIDC sẽ có `Password = ""`.

---

### Phase 4 — Routes

File: [`internal/user/port/http/routes.go`](internal/user/port/http/routes.go) — điều kiện theo `auth_mode`

```go
func Routes(r *gin.RouterGroup, sqlDB dbs.Database, validator validation.Validation) {
    userRepo := repository.NewUserRepository(sqlDB)
    userSvc := service.NewUserService(validator, userRepo)
    userHandler := NewUserHandler(userSvc)

    cfg := config.GetConfig()

    authRoute := r.Group("/auth")

    switch cfg.AuthMode {
    case config.AuthModeOIDC:
        // OIDC flow — không có register/login bằng email/password
        oidcSvc := service.NewOIDCService(validator, userRepo)
        oidcHandler := NewOIDCHandler(oidcSvc, cfg)
        authRoute.GET("/login", oidcHandler.Login)          // redirect → Authentik
        authRoute.GET("/callback", oidcHandler.Callback)    // nhận code từ Authentik
    case config.AuthModeJWT:
        // JWT flow truyền thống
        authMiddleware := middleware.JWTAuth()
        refreshAuthMiddleware := middleware.JWTRefresh()
        authRoute.POST("/register", userHandler.Register)
        authRoute.POST("/login", userHandler.Login)
        authRoute.POST("/refresh", refreshAuthMiddleware, userHandler.RefreshToken)
        authRoute.GET("/me", authMiddleware, userHandler.GetMe)
        authRoute.PUT("/change-password", authMiddleware, userHandler.ChangePassword)
    }
}
```

Tương tự cho gRPC interceptor trong [`internal/server/grpc/server.go`](internal/server/grpc/server.go):

```go
func NewServer(validator validation.Validation, db dbs.Database, cache redis.Redis) *Server {
    cfg := config.GetConfig()
    var interceptor grpc.UnaryServerInterceptor

    switch cfg.AuthMode {
    case config.AuthModeOIDC:
        oidcValidator, _ := oidc.NewValidator(...)
        interceptor = middleware.NewOIDCInterceptor(oidcValidator, config.AuthIgnoreMethods).Unary()
    default:
        interceptor = middleware.NewAuthInterceptor(config.AuthIgnoreMethods).Unary()
    }

    grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor))
    // ...
}
```

---

### Phase 5 — K8s Manifests & CI/CD

#### 5.1 Helm values cho GoShop

Cập nhật Helm values của GoShop để inject OIDC env vars từ K8s Secret:

```yaml
# deploy/helm/goshop/values.yaml
env:
  OIDC_ISSUER: "https://auth.cunghoclaptrinh.online/application/o/goshop-api/"
  OIDC_CLIENT_ID: "goshop-api-client"
  OIDC_CLIENT_SECRET:
    valueFrom:
      secretKeyRef:
        name: goshop-oidc-secret
        key: client-secret
  OIDC_REDIRECT_URL: "https://goshop.cunghoclaptrinh.online/api/v1/auth/callback"
  OIDC_JWKS_URL: "https://auth.cunghoclaptrinh.online/application/o/goshop-api/jwks/"
```

Tạo K8s Secret trước khi deploy:

```bash
kubectl create secret generic goshop-oidc-secret \
  --from-literal=client-secret=<AUTHENTIK_CLIENT_SECRET> \
  -n goshop
```

#### 5.2 CI/CD — `.github/workflows/ci.yml`

Thêm job `authentik` để test OIDC flow với testcontainer:

```yaml
authentik:
  runs-on: ubuntu-latest
  services:
    authentik:
      image: ghcr.io/goauthentik/server:2024.12.3
      env:
        AUTHENTIK_SECRET_KEY: test-secret
      ports:
        - 9000:9000
  steps:
    - uses: actions/checkout@v4
    - name: Wait for Authentik
      run: |
        for i in $(seq 1 30); do
          curl -sf http://localhost:9000/ && break || sleep 2
        done
    - name: Run OIDC integration tests
      run: go test -tags=oidc ./tests/integration/oidc/...
```

---

### Phase 6 — Migration & Cleanup

Với `auth_mode` flag, việc migration trở nên đơn giản — chỉ cần đổi giá trị config, không cần deploy lại.

#### 6.1 Cutover sang OIDC

- Chuyển `auth_mode: jwt` → `auth_mode: oidc` trên staging, verify full flow
- Canary 1% trên production, theo dõi error rate + latency
- Tăng dần 10% → 50% → 100%
- Users đăng nhập bằng OIDC sẽ được `UpsertUserFromOIDC` tự động tạo trong DB

#### 6.2 JWT code removal (Phase 8 — sau 2–4 tuần OIDC ổn định)

- Xóa `bcrypt` dependency
- Xóa `auth_secret` khỏi config
- Xóa `POST /auth/login`, `POST /auth/register`, `PUT /auth/change-password`
- Xóa `POST /auth/refresh`
- Xóa [`pkg/jtoken/`](pkg/jtoken/jwt.go)
- Xóa [`pkg/middleware/auth.go`](pkg/middleware/auth.go) (JWT version)
- Xóa `AuthModeJWT` và `auth_mode` khỏi config (không cần nữa)
- Cập nhật `AuthIgnoreMethods` trong [`pkg/config/config.go`](pkg/config/config.go)

---

## 5. Cấu trúc thư mục mới

```
goshop/
├── pkg/
│   ├── oidc/
│   │   ├── validator.go          # NEW — OIDC provider + JWKS token verification
│   │   └── validator_test.go     # NEW
│   └── middleware/
│       ├── oidc.go               # NEW — OIDC HTTP middleware
│       └── oidc_interceptor.go   # NEW — OIDC gRPC interceptor
├── internal/
│   └── user/
│       ├── model/
│       │   └── user.go           # MODIFIED — Password nullable (empty cho OIDC users)
│       ├── service/
│       │   └── user.go           # MODIFIED — + UpsertUserFromOIDC(sub, email, name)
│       ├── port/
│       │   ├── http/
│       │   │   ├── oidc_handlers.go   # NEW — GET /auth/login (redirect), GET /auth/callback
│       │   │   └── routes.go          # MODIFIED — mount OIDC routes
│       │   └── grpc/
│       │       └── handlers.go        # MODIFIED — optional OIDC interceptor
│       └── repository/
│           └── user.go           # MODIFIED — GetUserByOIDCSubject (nếu cần)
├── deploy/
│   └── helm/
│       └── goshop/               # MODIFIED — + OIDC env vars từ K8s Secret
│           └── values.yaml
├── config.sample.yaml            # MODIFIED — + OIDC config
└── docs/
    ├── authentik-integration.md  # NEW — runbook
    └── authentik-integration-plan.md  # THIS FILE
```

---

## 6. Sequence diagram — OIDC Login Flow

```
User    GoShop FE    GoShop API    Authentik    Postgres
 |          |             |             |            |
 |--GET /auth/login------>|             |            |
 |          |             |--302 /authorize-------->|
 |          |             |             |<--login UI--|
 |<--302 to Authentik------|             |            |
 |--login (email/pass)---------------------------->|
 |          |             |             |--200 code--|
 |<--302 /callback?code=xxx----------------------|
 |--GET /auth/callback?code=xxx------------------>|
 |          |             |--POST /token------------>|
 |          |             |<--id+access+refresh-------|
 |          |             |--GET /userinfo----------->|
 |          |             |<--{sub, email, name}-------|
 |          |             |--UPSERT user-------------->|
 |          |             |<--user record--------------|
 |          |             |--200 {token, user}---------|
 |<--200 {access, refresh}--|             |            |
```

---

## 7. Testing strategy

| Test type | Scope | Tool |
|-----------|-------|------|
| Unit | `oidc.Validator.Verify`, `OIDCAuth` middleware | `testify` |
| Integration | Full OIDC flow với Authentik testcontainer | `testcontainers` |
| E2E | FE → API → Authentik → DB | Playwright / Cypress |
| Contract | Proto gRPC không đổi | `buf` |

---

## 8. Rollout plan

> Authentik đã được deploy lên K8s bởi repo khác. Phần này chỉ liệt kê công việc phía GoShop.

Feature flag `auth_mode` (`jwt` | `oidc`) cho phép chuyển đổi giữa 2 flow mà không cần deploy lại.

| Phase | Mô tả | Thời gian ước tính |
|-------|-------|-------------------|
| Phase 1 | Config + `auth_mode` flag + dependencies | 0.5 ngày |
| Phase 2 | OIDC middleware (HTTP + gRPC) | 1 ngày |
| Phase 3 | Auth handlers + user upsert | 1 ngày |
| Phase 4 | Routes + Helm values + CI | 0.5 ngày |
| Phase 5 | Testing | 1 ngày |
| Phase 6 | Canary: `auth_mode=oidc` cho 1% traffic | 2 ngày |
| Phase 7 | Full cutover: `auth_mode=oidc` 100% | 0.5 ngày |
| Phase 8 | Xóa JWT code (sau khi OIDC ổn định) | 1 ngày |

**Tổng: ~5–6 ngày dev**

### Rollout steps

1. **Deploy với `auth_mode: jwt`** — không đổi hành vi hiện tại
2. **Chuyển `auth_mode: oidc`** trên staging, test full flow
3. **Canary 1%** trên production, theo dõi error rate
4. **Tăng dần** 10% → 50% → 100%
5. **Xóa JWT code** sau 2–4 tuần OIDC ổn định

---

## 9. Risks & Mitigations

| Risk | Mitigation |
|------|------------|
| Authentik downtime | Có fallback sang JWT cũ trong dual-auth mode |
| JWKS cache stale | Cache 1h, force refresh khi `kid` không tìm thấy |
| User email thay đổi trên Authentik | Dùng `sub` (OIDC subject) làm stable ID, không dùng email |
| Migration data user cũ | Giữ nguyên password hash, cho user đăng nhập cả 2 cách |
| MFA bắt buộc trên Authentik | Cấu hình flow Authentik cho phép optional MFA |
