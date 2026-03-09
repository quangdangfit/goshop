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

func TestUnary_IgnoredMethod(t *testing.T) {
	ai := NewAuthInterceptor([]string{"/user.UserService/Login"})
	interceptor := ai.Unary()

	info := &grpc.UnaryServerInfo{FullMethod: "/user.UserService/Login"}
	resp, err := interceptor(context.Background(), nil, info, fakeHandler)
	assert.NoError(t, err)
	assert.Equal(t, "ok", resp)
}

func TestUnary_Authorized(t *testing.T) {
	ai := NewAuthInterceptor([]string{})
	interceptor := ai.Unary()

	token := makeToken("user-123")
	md := metadata.Pairs("token", token)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	info := &grpc.UnaryServerInfo{FullMethod: "/product.ProductService/GetProductByID"}
	resp, err := interceptor(ctx, nil, info, fakeHandler)
	assert.NoError(t, err)
	assert.Equal(t, "ok", resp)
}

func TestUnary_MissingToken(t *testing.T) {
	ai := NewAuthInterceptor([]string{})
	interceptor := ai.Unary()

	info := &grpc.UnaryServerInfo{FullMethod: "/product.ProductService/GetProductByID"}
	resp, err := interceptor(context.Background(), nil, info, fakeHandler)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestUnary_InvalidToken(t *testing.T) {
	ai := NewAuthInterceptor([]string{})
	interceptor := ai.Unary()

	md := metadata.Pairs("token", "invalid.token.value")
	ctx := metadata.NewIncomingContext(context.Background(), md)

	info := &grpc.UnaryServerInfo{FullMethod: "/product.ProductService/GetProductByID"}
	resp, err := interceptor(ctx, nil, info, fakeHandler)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestUnary_NoMetadata(t *testing.T) {
	ai := NewAuthInterceptor([]string{})
	interceptor := ai.Unary()

	info := &grpc.UnaryServerInfo{FullMethod: "/product.ProductService/GetProductByID"}
	// plain context with no metadata
	resp, err := interceptor(context.Background(), nil, info, fakeHandler)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestUnary_EmptyTokenInMetadata(t *testing.T) {
	ai := NewAuthInterceptor([]string{})
	interceptor := ai.Unary()

	// Metadata exists but "token" key has no values
	md := metadata.MD{}
	ctx := metadata.NewIncomingContext(context.Background(), md)

	info := &grpc.UnaryServerInfo{FullMethod: "/product.ProductService/GetProductByID"}
	resp, err := interceptor(ctx, nil, info, fakeHandler)
	assert.Error(t, err)
	assert.Nil(t, resp)
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
