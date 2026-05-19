package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"goshop/pkg/apperror"
	"goshop/pkg/oidc"
)

// OIDCInterceptor validates Authentik-issued Bearer ID tokens on incoming
// gRPC requests. Methods listed in `ignored` (e.g. login/register) are
// allowed through unauthenticated.
type OIDCInterceptor struct {
	validator *oidc.Validator
	ignored   []string
}

// NewOIDCInterceptor creates a new OIDC gRPC interceptor.
func NewOIDCInterceptor(validator *oidc.Validator, ignored []string) *OIDCInterceptor {
	return &OIDCInterceptor{
		validator: validator,
		ignored:   ignored,
	}
}

// Unary returns a grpc.UnaryServerInterceptor.
func (i *OIDCInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		for _, m := range i.ignored {
			if info.FullMethod == m {
				return handler(ctx, req)
			}
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, apperror.ErrUnauthorized
		}

		tokens := md.Get("authorization")
		if len(tokens) == 0 {
			return nil, apperror.ErrUnauthorized
		}

		idToken, err := i.validator.Verify(ctx, tokens[0])
		if err != nil {
			return nil, apperror.Wrap(apperror.ErrUnauthorized, err)
		}

		var claims struct {
			Sub    string   `json:"sub"`
			Groups []string `json:"groups,omitempty"`
		}
		if err := idToken.Claims(&claims); err != nil {
			return nil, apperror.Wrap(apperror.ErrUnauthorized, err)
		}

		role := mapAuthentikGroupsToRole(claims.Groups)

		// Inject userId/role into the context so downstream gRPC handlers and
		// the AdminOnly interceptor keep working unchanged.
		ctx = context.WithValue(ctx, "userId", claims.Sub)
		ctx = context.WithValue(ctx, "role", role)
		return handler(ctx, req)
	}
}
