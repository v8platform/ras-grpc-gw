package interceptor

import (
	"context"
	clientv1 "github.com/v8platform/protos/gen/ras/client/v1"
	client "github.com/v8platform/ras-grpc-gw/pkg/ras_client"
)

type Cond func(endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}) bool

func OR(conds ...Cond) Cond {
	return func(endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}) bool {
		for _, cond := range conds {
			if cond(endpoint, info, req) {
				return true
			}
		}
		return false
	}
}

func AND(conds ...Cond) Cond {
	return func(endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}) bool {
		for _, cond := range conds {
			if !cond(endpoint, info, req) {
				return false
			}
		}
		return true
	}
}

type getClusterId interface {
	GetClusterID() string
}

func GetClusterId(_ clientv1.Endpoint, _ *clientv1.RequestInfo, req interface{}) bool {
	_, ok := req.(getClusterId)
	return ok
}

func IsEndpoint(endpoint clientv1.Endpoint, _ *clientv1.RequestInfo, _ interface{}) bool {
	return endpoint != nil
}

func New(condition Cond, interceptor client.Interceptor) client.Interceptor {

	return func(ctx context.Context, channel clientv1.Channel, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}, handler clientv1.InterceptorHandler) (interface{}, error) {
		if condition(endpoint, info, req) {
			return handler(ctx, channel, endpoint, req)
		}
		return interceptor(ctx, channel, endpoint, info, req, handler)
	}
}
