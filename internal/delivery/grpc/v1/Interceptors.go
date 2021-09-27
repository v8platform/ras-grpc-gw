package v1

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	clientv1 "github.com/v8platform/protos/gen/ras/client/v1"
	messagesv1 "github.com/v8platform/protos/gen/ras/messages/v1"
	appCtx "github.com/v8platform/ras-grpc-gw/internal/context"
	"github.com/v8platform/ras-grpc-gw/internal/service"
	client2 "github.com/v8platform/ras-grpc-gw/pkg/ras_client"
	"github.com/v8platform/ras-grpc-gw/pkg/ras_client/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

func NewInterceptors(services *service.Services) []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		// grpc_auth.UnaryServerInterceptor(authTokenFunc(services)),
		// getClientFunc(services),
		getEndpointFunc(services),
		getRequestOpts(services),
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

func ClusterAuthInterceptor(user, password string) client2.Interceptor {

	return interceptor.New(
		interceptor.AND(
			interceptor.GetClusterId,
			interceptor.IsEndpoint,
		),
		setClusterAuthInterceptor(user, password),
	)
}

func setClusterAuthInterceptor(user, password string) client2.Interceptor {

	type getClusterId interface {
		GetClusterID() string
	}

	return func(ctx context.Context, channel clientv1.Channel, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}, handler clientv1.InterceptorHandler) (interface{}, error) {

		if endpoint == nil {
			return handler(ctx, channel, endpoint, req)
		}

		tReq := req.(getClusterId)
		clusterId := tReq.GetClusterID()

		_, err := clientv1.AuthenticateClusterHandler(ctx, channel, endpoint, &messagesv1.ClusterAuthenticateRequest{
			ClusterId: clusterId,
			User:      user,
			Password:  password,
		}, nil)

		if err != nil {
			fmt.Println(err)
		}

		return handler(ctx, channel, endpoint, req)
	}
}

func setClusterIDToRequestInterceptor(clusterId uuid.UUID) client2.Interceptor {
	return func(ctx context.Context, channel clientv1.Channel, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}, handler clientv1.InterceptorHandler) (interface{}, error) {

		switch tReq := req.(type) {
		case *messagesv1.GetInfobasesSummaryRequest:
			tReq.ClusterId = clusterId.String()
		}

		return handler(ctx, channel, endpoint, req)
	}
}

func getRequestOpts(services *service.Services) grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req)
		}

		var opts []interface{}
		var requestInterceptors []client2.Interceptor
		for key, values := range md {

			if strings.HasPrefix(key, "x-req-") {
				switch strings.Trim(key, "x-req-") {
				case "cluster-id":
					for _, value := range values {
						switch len(value) {
						case 36:
							clusterId := uuid.MustParse(value)

							requestInterceptors = append(requestInterceptors, setClusterIDToRequestInterceptor(clusterId))
							break
						default:
							// Получить кластер по имени
						}
					}
				case "cluster-auth":
					for _, value := range values {
						switch len(value) {
						case 36:
							clusterId := uuid.MustParse(value)
							requestInterceptors = append(requestInterceptors, ClusterAuthInterceptor(clusterId))
							break
						default:
							// Получить кластер по имени
						}
					}
				}
			}

		}
		requestInterceptors = append(requestInterceptors)

		if len(requestInterceptors) > 0 {
			opts = append(opts, client2.RequestInterceptor(requestInterceptors...))
		}

		if len(opts) > 0 {

			ctx = appCtx.RequestOptsToContext(ctx, opts)

		}

		h, err := handler(ctx, req)
		return h, err

	}
}

func getEndpointFunc(services *service.Services) grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		endpoint := metautils.ExtractIncoming(ctx).Get("x-endpoint")

		if len(endpoint) > 0 {

			ctx = appCtx.EndpointToContext(ctx, endpoint)

		}

		h, err := handler(ctx, req)
		return h, err

	}
}
