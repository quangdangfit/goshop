package oidc

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// Validator wraps the Authentik OIDC provider, its ID-token verifier, and the
// OAuth2 client configuration used for the authorization-code flow.
//
// go-oidc handles JWKS fetching and caching internally, so callers only need
// Verify() / Exchange() / AuthCodeURL().
type Validator struct {
	provider     *oidc.Provider
	verifier     *oidc.IDTokenVerifier
	oauth2Config *oauth2.Config
}

// NewValidator discovers the provider metadata at boot and builds a reusable
// validator. The call is blocking and will fail if the issuer is unreachable —
// the caller (typically server startup) should treat that as fatal.
func NewValidator(issuer, clientID, clientSecret, redirectURL string) (*Validator, error) {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, fmt.Errorf("failed to discover OIDC provider: %w", err)
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: clientID})

	return &Validator{
		provider: provider,
		verifier: verifier,
		oauth2Config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Endpoint:     provider.Endpoint(),
			Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
		},
	}, nil
}

// Verify validates an ID token's signature, issuer, audience, and expiry, and
// returns the parsed token so callers can extract claims.
func (v *Validator) Verify(ctx context.Context, rawIDToken string) (*oidc.IDToken, error) {
	return v.verifier.Verify(ctx, rawIDToken)
}

// AuthCodeURL builds the Authentik authorize URL for the given state value.
func (v *Validator) AuthCodeURL(state string) string {
	return v.oauth2Config.AuthCodeURL(state)
}

// Exchange trades an authorization code for an OAuth2 token (including the
// id_token in the "id_token" extra field).
func (v *Validator) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return v.oauth2Config.Exchange(ctx, code)
}
