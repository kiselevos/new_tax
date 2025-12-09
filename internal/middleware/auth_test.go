package middleware

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const privateMethod = "/tax.TaxService/CalculatePrivate"
const publicMethod = "/tax.TaxService/CalculatePublic"

func allowHandler(ctx context.Context, req interface{}) (interface{}, error) {
	return "OK", nil
}

func TestAuth_PublicMethodAllowed(t *testing.T) {
	interceptor := Auth("secret")

	_, err := interceptor(context.Background(), nil,
		&grpc.UnaryServerInfo{FullMethod: publicMethod},
		allowHandler)

	require.NoError(t, err)
}

func TestAuth_PrivateWithInternalHeader(t *testing.T) {
	interceptor := Auth("secret")

	md := metadata.Pairs("x-internal", "true")
	ctx := metadata.NewIncomingContext(context.Background(), md)

	resp, err := interceptor(ctx, nil,
		&grpc.UnaryServerInfo{FullMethod: privateMethod},
		allowHandler)

	require.NoError(t, err)
	require.Equal(t, "OK", resp)
}

func TestAuth_PrivateWithCorrectAPIKey(t *testing.T) {
	interceptor := Auth("secret")

	md := metadata.Pairs("x-api-key", "secret")
	ctx := metadata.NewIncomingContext(context.Background(), md)

	resp, err := interceptor(ctx, nil,
		&grpc.UnaryServerInfo{FullMethod: privateMethod},
		allowHandler)

	require.NoError(t, err)
	require.Equal(t, "OK", resp)
}

func TestAuth_PrivateWithIncorrectAPIKey(t *testing.T) {
	interceptor := Auth("secret")

	md := metadata.Pairs("x-api-key", "wrong")
	ctx := metadata.NewIncomingContext(context.Background(), md)

	_, err := interceptor(ctx, nil,
		&grpc.UnaryServerInfo{FullMethod: privateMethod},
		allowHandler)

	require.Error(t, err)
	st, _ := status.FromError(err)
	require.Equal(t, codes.PermissionDenied, st.Code())
}

func TestAuth_PrivateWithoutAnyAuth(t *testing.T) {
	interceptor := Auth("secret")

	_, err := interceptor(context.Background(), nil,
		&grpc.UnaryServerInfo{FullMethod: privateMethod},
		allowHandler)

	require.Error(t, err)
	st, _ := status.FromError(err)

	require.Equal(t, codes.Unauthenticated, st.Code())
}
