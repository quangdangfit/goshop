package config

import (
	"log"
	"os"
	"time"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

const (
	ProductionEnv = "production"

	DatabaseTimeout    = 5 * time.Second
	ProductCachingTime = 1 * time.Minute
)

// AuthMode selects which authentication scheme the server runs in.
type AuthMode string

const (
	AuthModeJWT  AuthMode = "jwt"  // legacy email + bcrypt + HS256 JWT
	AuthModeOIDC AuthMode = "oidc" // OIDC via Authentik
)

var AuthIgnoreMethods = []string{
	"/user.UserService/Login",
	"/user.UserService/Register",
}

type Schema struct {
	Environment   string `env:"environment"`
	HttpPort      int    `env:"http_port"`
	GrpcPort      int    `env:"grpc_port"`
	AuthSecret    string `env:"auth_secret"`
	DatabaseURI   string `env:"database_uri"`
	RedisURI      string `env:"redis_uri"`
	RedisPassword string `env:"redis_password"`
	RedisDB       int    `env:"redis_db"`

	CORSAllowedOrigins     string `env:"cors_allowed_origins" envDefault:"*"`
	RateLimitRequests      int    `env:"rate_limit_requests" envDefault:"100"`
	RateLimitWindowSeconds int    `env:"rate_limit_window_seconds" envDefault:"60"`

	StripeSecretKey      string `env:"stripe_secret_key"`
	StripeWebhookSecret  string `env:"stripe_webhook_secret"`
	StripePublishableKey string `env:"stripe_publishable_key"`
	StripeAPIBase        string `env:"stripe_api_base"` // override for stripe-mock in tests

	SMTPHost     string `env:"smtp_host"`
	SMTPPort     int    `env:"smtp_port" envDefault:"25"`
	SMTPUser     string `env:"smtp_user"`
	SMTPPassword string `env:"smtp_password"`
	EmailFrom    string `env:"email_from"`

	// OIDC / Authentik
	AuthMode         AuthMode `env:"auth_mode" envDefault:"jwt"`
	OIDCIssuer       string   `env:"oidc_issuer"`
	OIDCClientID     string   `env:"oidc_client_id"`
	OIDCClientSecret string   `env:"oidc_client_secret"`
	OIDCRedirectURL  string   `env:"oidc_redirect_url"`
	OIDCJWKSURL      string   `env:"oidc_jwks_url"`
	OIDCScopes       string   `env:"oidc_scopes" envDefault:"openid,email,profile"`
	// FrontendBaseURL is where the OIDC callback redirects the browser to after
	// minting the GoShop JWT pair (e.g. https://goshop.example.com). The FE must
	// expose a /auth/callback route that reads access_token + refresh_token
	// from the query string.
	FrontendBaseURL string `env:"frontend_base_url" envDefault:"http://localhost:5173"`
}

var (
	cfg Schema
)

// ConfigFile is the path (relative to working dir or absolute) to load before
// env-var parsing. Override by setting CONFIG_FILE in the environment.
var ConfigFile = "config.yaml"

func LoadConfig() *Schema {
	path := ConfigFile
	if v := os.Getenv("CONFIG_FILE"); v != "" {
		path = v
	}

	if err := godotenv.Load(path); err != nil {
		log.Printf("Error on load configuration file, error: %v", err)
	}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Error on parsing configuration file, error: %v", err)
	}

	return &cfg
}

func GetConfig() *Schema {
	return &cfg
}
