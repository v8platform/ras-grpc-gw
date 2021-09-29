package client

import (
	"context"
	"fmt"
	clientv1 "github.com/v8platform/protos/gen/ras/client/v1"
	messagesv1 "github.com/v8platform/protos/gen/ras/messages/v1"
	"github.com/v8platform/ras-grpc-gw/pkg/ras_client/cond"
	"github.com/v8platform/ras-grpc-gw/pkg/ras_client/md"
	"os"
)

type Interceptor clientv1.Interceptor
type InterceptorHandler clientv1.InterceptorHandler

func ChainInterceptor(interceptors ...Interceptor) Interceptor {
	n := len(interceptors)

	return func(ctx context.Context, channel clientv1.Channel, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}, handler clientv1.InterceptorHandler) (interface{}, error) {
		chainer := func(currentInter clientv1.Interceptor, currentHandler clientv1.InterceptorHandler) clientv1.InterceptorHandler {
			return func(currentCtx context.Context, channel clientv1.Channel, endpoint clientv1.Endpoint, currentReq interface{}) (interface{}, error) {
				return currentInter(currentCtx, channel, endpoint, info, currentReq, currentHandler)
			}
		}

		chainedHandler := handler
		for i := n - 1; i >= 0; i-- {
			chainedHandler = chainer(clientv1.Interceptor(interceptors[i]), chainedHandler)
		}

		return chainedHandler(ctx, channel, endpoint, req)
	}
}

func SetClusterIDToRequestInterceptor() Interceptor {
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
		case *messagesv1.GetWorkingServersRequest:
			tReq.ClusterId = clusterId
		}

		return handler(ctx, channel, endpoint, req)
	}
}

func InfobaseAuthInterceptor() Interceptor {
	return cond.And{
		cond.Cond(HasClusterId),
		cond.Cond(IsEndpointRequest),
	}.Handler(infobaseAuthInterceptor)
}

func infobaseAuthInterceptor(ctx context.Context, channel clientv1.Channel, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}, handler clientv1.InterceptorHandler) (interface{}, error) {

	if cond.Condition(
		cond.And{
			cond.Cond(HasClusterId),
			cond.Cond(IsEndpointRequest),
		}, endpoint, info, req) {
		return handler(ctx, channel, endpoint, req)
	}

	type getClusterId interface {
		GetClusterID() string
	}

	reqMd := md.ExtractMetadata(ctx)
	clusterId := reqMd.Get("cluster-id")

	if len(clusterId) == 0 {
		tReq := req.(getClusterId)
		clusterId = tReq.GetClusterID()
	}

	if !(reqMd.Has("infobase-user") && reqMd.Has("infobase-password")) {
		return handler(ctx, channel, endpoint, req)
	}

	// user := reqMd.Get("infobase-user")
	// password := reqMd.Get("infobase-password")

	// TODO Для тестов
	user := os.Getenv("IB_USER")
	password := os.Getenv("IB_PWD")

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
