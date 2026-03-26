package middleware

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"goshop/pkg/config"
	"goshop/pkg/jtoken"
)

func init() {
	config.LoadConfig()
}

func makeToken(id string) string {
	return jtoken.GenerateAccessToken(map[string]interface{}{
		"id":    id,
		"email": "test@example.com",
		"role":  "customer",
	})
}

func fakeHandler(ctx context.Context, req interface{}) (interface{}, error) {
	return "ok", nil
}

func TestUnary(t *testing.T) {
	tests := []struct {
		name      string
		ignored   []string
		method    string
		ctx       func() context.Context
		wantResp  interface{}
		wantError bool
	}{
		{
			name:    "ignored method passes through",
			ignored: []string{"/user.UserService/Login"},
			method:  "/user.UserService/Login",
			ctx: func() context.Context {
				return context.Background()
			},
			wantResp:  "ok",
			wantError: false,
		},
		{
			name:    "authorized with valid token",
			ignored: []string{},
			method:  "/product.ProductService/GetProductByID",
			ctx: func() context.Context {
				token := makeToken("user-123")
				md := metadata.Pairs("token", token)
				return metadata.NewIncomingContext(context.Background(), md)
			},
			wantResp:  "ok",
			wantError: false,
		},
		{
			name:    "missing token",
			ignored: []string{},
			method:  "/product.ProductService/GetProductByID",
			ctx: func() context.Context {
				return context.Background()
			},
			wantResp:  nil,
			wantError: true,
		},
		{
			name:    "invalid token",
			ignored: []string{},
			method:  "/product.ProductService/GetProductByID",
			ctx: func() context.Context {
				md := metadata.Pairs("token", "invalid.token.value")
				return metadata.NewIncomingContext(context.Background(), md)
			},
			wantResp:  nil,
			wantError: true,
		},
		{
			name:    "no metadata in context",
			ignored: []string{},
			method:  "/product.ProductService/GetProductByID",
			ctx: func() context.Context {
				return context.Background()
			},
			wantResp:  nil,
			wantError: true,
		},
		{
			name:    "empty token in metadata",
			ignored: []string{},
			method:  "/product.ProductService/GetProductByID",
			ctx: func() context.Context {
				md := metadata.MD{}
				return metadata.NewIncomingContext(context.Background(), md)
			},
			wantResp:  nil,
			wantError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ai := NewAuthInterceptor(tc.ignored)
			interceptor := ai.Unary()

			info := &grpc.UnaryServerInfo{FullMethod: tc.method}
			resp, err := interceptor(tc.ctx(), nil, info, fakeHandler)
			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantResp, resp)
		})
	}
}

func TestAuthorize_UserIDPropagated(t *testing.T) {
	ai := NewAuthInterceptor([]string{})
	interceptor := ai.Unary()

	token := makeToken("user-abc")
	md := metadata.Pairs("token", token)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	var capturedUserID string
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		capturedUserID, _ = ctx.Value("userId").(string)
		return nil, nil
	}

	info := &grpc.UnaryServerInfo{FullMethod: "/any/Method"}
	_, _ = interceptor(ctx, nil, info, handler)
	assert.Equal(t, "user-abc", capturedUserID)
}
