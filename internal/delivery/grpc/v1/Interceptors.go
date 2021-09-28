package v1

import (
	"context"
	"encoding/base64"
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
	"github.com/v8platform/ras-grpc-gw/pkg/ras_client/md"
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

func ClusterAuthInterceptor() client2.Interceptor {

	return interceptor.New(
		interceptor.AND(
			interceptor.GetClusterId,
			interceptor.IsEndpoint,
		),
		setClusterAuthInterceptor(),
	)
}

func setClusterAuthInterceptor() client2.Interceptor {

	return func(ctx context.Context, channel clientv1.Channel, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}, handler clientv1.InterceptorHandler) (interface{}, error) {
		type getClusterId interface {
			GetClusterID() string
		}

		reqMd := md.ExtractMetadata(ctx)
		clusterId := reqMd.Get("cluster-id")

		if len(clusterId) == 0 {
			tReq := req.(getClusterId)
			clusterId = tReq.GetClusterID()
		}

		user := reqMd.Get("cluster-user")
		password := reqMd.Get("cluster-password")

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

func InfobaseAuthInterceptor() client2.Interceptor {

	return interceptor.New(
		interceptor.AND(
			interceptor.GetClusterId,
			interceptor.IsEndpoint,
		),
		setInfobaseAuthInterceptor(),
	)
}

func setInfobaseAuthInterceptor() client2.Interceptor {

	return func(ctx context.Context, channel clientv1.Channel, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}, handler clientv1.InterceptorHandler) (interface{}, error) {
		type getClusterId interface {
			GetClusterID() string
		}

		reqMd := md.ExtractMetadata(ctx)
		clusterId := reqMd.Get("cluster-id")

		if len(clusterId) == 0 {
			tReq := req.(getClusterId)
			clusterId = tReq.GetClusterID()
		}

		user := reqMd.Get("infobase-user")
		password := reqMd.Get("infobase-password")

		_, err := clientv1.AuthenticateInfobaseHandler(ctx, channel, endpoint, &messagesv1.AuthenticateInfobaseRequest{
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

func SetClusterIDToRequestInterceptor() client2.Interceptor {
	return func(ctx context.Context, channel clientv1.Channel, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}, handler clientv1.InterceptorHandler) (interface{}, error) {
		reqMd := md.ExtractMetadata(ctx)
		clusterId := reqMd.Get("cluster-id")
		if len(clusterId) == 0 {
			return handler(ctx, channel, endpoint, req)
		}

		switch tReq := req.(type) {
		case *messagesv1.GetInfobasesSummaryRequest:
			tReq.ClusterId = clusterId
		case *messagesv1.GetInfobasesRequest:
			tReq.ClusterId = clusterId
		case *messagesv1.GetSessionsRequest:
			tReq.ClusterId = clusterId
		case *messagesv1.GetInfobaseInfoRequest:
			tReq.ClusterId = clusterId
		}

		return handler(ctx, channel, endpoint, req)
	}
}

func AnnotateRequestMetadata(_ *service.Services) md.AnnotationHandler {

	return func(ctx context.Context, req interface{}) md.RequestMetadata {

		grpcmd, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return md.Pairs()
		}

		var pairs []string

		for key, values := range grpcmd {

			if strings.HasPrefix(key, "x-req-") {
				switch strings.Trim(key, "x-req-") {
				case "cluster-id":
					var clusterId string

					for _, value := range values {

						switch len(value) {
						case 36:

							_ = uuid.MustParse(value)
							clusterId = value

							break
						default:
							// Получить кластер по имени
						}
					}
					if len(clusterId) > 0 {
						pairs = append(pairs, "cluster-id", clusterId)
					}

				case "cluster-auth":
					var username, password string

					for _, value := range values {
						if len(value) == 0 {
							continue
						}
						var ok bool
						username, password, ok = parseBasicAuth(value)
						if !ok {
							continue
						}
					}

					if len(username) > 0 {
						pairs = append(pairs, "cluster-user", username)
						pairs = append(pairs, "cluster-password", password)
					}
				case "infobase-auth":
					var username, password string

					for _, value := range values {
						if len(value) == 0 {
							continue
						}
						var ok bool
						username, password, ok = parseBasicAuth(value)
						if !ok {
							continue
						}
					}

					if len(username) > 0 {
						pairs = append(pairs, "infobase-user", username)
						pairs = append(pairs, "infobase-password", password)
					}
				}
			}

		}

		return md.Pairs(pairs...)

	}
}

// parseBasicAuth parses an Basic Authentication string.
// "QWxhZGRpbjpvcGVuIHNlc2FtZQ==" returns ("Aladdin", "open sesame", true).
func parseBasicAuth(auth string) (username, password string, ok bool) {

	c, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
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
