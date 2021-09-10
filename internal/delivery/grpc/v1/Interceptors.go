package v1

import (
	"context"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	appCtx "github.com/v8platform/ras-grpc-gw/internal/context"
	"github.com/v8platform/ras-grpc-gw/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewInterceptors(services *service.Services) []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		grpc_auth.UnaryServerInterceptor(authTokenFunc(services)),
		getClientFunc(services),
		getEndpointFunc(services),
	}
}

func authTokenFunc(services *service.Services) grpc_auth.AuthFunc {

	return func(ctx context.Context) (context.Context, error) {

		token, err := grpc_auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, err
		}

		tokenInfo, err := services.TokenManager.Validate(token, "access")
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
		}

		if len(tokenInfo) > 0 {

			user, err := services.Users.GetByUUID(ctx, tokenInfo)
			if err != nil {
				return nil, err
			}
			ctx = appCtx.UserToContext(ctx, user)

		}
		return ctx, nil
	}
}

func getClientFunc(services *service.Services) grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var client string
		token := metautils.ExtractIncoming(ctx).Get("x-client")
		if token != "" {
			tokenInfo, err := services.TokenManager.Validate(token, "access")
			if err != nil {
				return nil, status.Errorf(codes.Unavailable, "invalid client token: %v", err)
			}
			client = tokenInfo
		}

		if len(client) > 0 {

			client, err := services.Clients.GetByUUID(ctx, client)
			if err != nil {
				return nil, err
			}
			ctx = appCtx.ClientToContext(ctx, client)

		}

		h, err := handler(ctx, req)

		return h, err

	}
}

func getEndpointFunc(services *service.Services) grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		endpoint := metautils.ExtractIncoming(ctx).Get("x-endpoint")
		newCtx := context.WithValue(ctx, "x-endpoint", endpoint)
		h, err := handler(newCtx, req)
		return h, err

	}
}