package client

import (
	"context"
	clientv1 "github.com/v8platform/protos/gen/ras/client/v1"
	messagesv1 "github.com/v8platform/protos/gen/ras/messages/v1"
)

type Interceptor clientv1.Interceptor
type InterceptorHandler clientv1.InterceptorHandler

type InterceptorCond struct {
	Regexp   string
	Requests []string
	Services []string
}

func (c InterceptorCond) Cond(info *clientv1.RequestInfo) bool {
	return false
}

func NewInterceptor(data InterceptorCond, h Interceptor) Interceptor {
	return func(ctx context.Context, channel clientv1.Channel, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}, handler clientv1.InterceptorHandler) (interface{}, error) {
		if data.Cond(info) {
			return h(ctx, channel, endpoint, info, req, handler)
		}
		return handler(ctx, channel, endpoint, req)
	}
}

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

func AuthInterceptor(service clientv1.AuthService) Interceptor {

	return NewInterceptor(InterceptorCond{Requests: []string{"GetInfobase"}}, func(ctx context.Context, channel clientv1.Channel, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}, handler clientv1.InterceptorHandler) (interface{}, error) {

		var (
			hasAuth   bool
			clusterId string
		)

		type GetCluster interface {
			GetCluster() string
		}

		switch req.(type) {
		case GetCluster:
			clusterId = req.(GetCluster).GetCluster()
		default:
			return handler(ctx, channel, endpoint, req)
		}

		if hasAuth {
			_, err := service.AuthenticateInfobase(ctx, &messagesv1.AuthenticateInfobaseRequest{
				ClusterId: clusterId,
				User:      "user",
				Password:  "pwd",
			})

			if err != nil {
				return nil, err
			}
		}

		return handler(ctx, channel, endpoint, req)
	})
}
