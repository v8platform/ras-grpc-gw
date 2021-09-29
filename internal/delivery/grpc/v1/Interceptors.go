package v1

import (
	"context"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	clientv1 "github.com/v8platform/protos/gen/ras/client/v1"
	appCtx "github.com/v8platform/ras-grpc-gw/internal/context"
	"github.com/v8platform/ras-grpc-gw/internal/service"
	client2 "github.com/v8platform/ras-grpc-gw/pkg/ras_client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func NewInterceptors(services *service.Services) []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		// grpc_auth.UnaryServerInterceptor(authTokenFunc(services)),
		// getClientFunc(services),
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
			return nil, status.Errorf(codes.Unauthenticated, "invalid Auth token: %v", err)
		}

		if len(tokenInfo) > 0 {

			// user, err := services.Users.GetByUUID(ctx, tokenInfo)
			// if err != nil {
			// 	return nil, err
			// }
			// ctx = appCtx.UserToContext(ctx, user)

		}
		return ctx, nil
	}
}

func SendEndpointID(ctx context.Context, channel clientv1.Channel, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}, handler clientv1.InterceptorHandler) (interface{}, error) {

	if endpoint == nil {
		return handler(ctx, channel, endpoint, req)
	}

	cnEndpoint := endpoint.(*client2.ChannelEndpoint)

	err := grpc.SetHeader(ctx, metadata.Pairs(
		"x-endpoint", cnEndpoint.UUID.String()))
	if err != nil {
		return nil, err
	}

	return handler(ctx, channel, endpoint, req)
}

func getEndpointFunc(_ *service.Services) grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		endpoint := metautils.ExtractIncoming(ctx).Get("x-endpoint")

		if len(endpoint) > 0 {

			ctx = appCtx.EndpointToContext(ctx, endpoint)

		}

		h, err := handler(ctx, req)
		return h, err

	}
}
