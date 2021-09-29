package cond

import (
	"context"
	clientv1 "github.com/v8platform/protos/gen/ras/client/v1"
	client "github.com/v8platform/ras-grpc-gw/pkg/ras_client"
)

type Cond func(endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}) bool

func (fn Cond) Condition(endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}) bool {
	return fn(endpoint, info, req)
}
func (fn Cond) Handler(interceptor ...client.Interceptor) client.Interceptor {
	return build(fn, interceptor...)
}

type cond interface {
	Handler(interceptor ...client.Interceptor) client.Interceptor
	Condition(endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}) bool
}

type Or []cond

func (c Or) Condition(endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}) bool {
	for _, cond := range c {
		if cond.Condition(endpoint, info, req) {
			return true
		}
	}
	return false
}

func (c Or) Handler(interceptor ...client.Interceptor) client.Interceptor {
	return build(c, interceptor...)
}

type Xor []cond

func (c Xor) Condition(endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}) bool {
	for _, cond := range c {
		if cond.Condition(endpoint, info, req) {
			return false
		}
	}
	return true
}
func (c Xor) Handler(interceptor ...client.Interceptor) client.Interceptor {
	return build(c, interceptor...)
}

type Not [1]cond

func (c Not) Condition(endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}) bool {
	for _, cond := range c {
		if cond.Condition(endpoint, info, req) {
			return false
		}
	}
	return true
}
func (c Not) Handler(interceptor ...client.Interceptor) client.Interceptor {
	return build(c, interceptor...)
}

type And []cond

func (c And) Condition(endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}) bool {
	for _, cond := range c {
		if !cond.Condition(endpoint, info, req) {
			return false
		}
	}
	return true
}

func (c And) Handler(interceptor ...client.Interceptor) client.Interceptor {
	return build(c, interceptor...)
}

func build(condition cond, interceptor ...client.Interceptor) client.Interceptor {
	n := len(interceptor)

	return func(ctx context.Context, channel clientv1.Channel, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}, handler clientv1.InterceptorHandler) (interface{}, error) {
		if condition.Condition(endpoint, info, req) {
			return handler(ctx, channel, endpoint, req)
		}
		chainer := func(currentInter clientv1.Interceptor, currentHandler clientv1.InterceptorHandler) clientv1.InterceptorHandler {
			return func(currentCtx context.Context, channel clientv1.Channel, endpoint clientv1.Endpoint, currentReq interface{}) (interface{}, error) {
				return currentInter(currentCtx, channel, endpoint, info, currentReq, currentHandler)
			}
		}

		chainedHandler := handler
		for i := n - 1; i >= 0; i-- {
			chainedHandler = chainer(clientv1.Interceptor(interceptor[i]), chainedHandler)
		}

		return chainedHandler(ctx, channel, endpoint, req)
	}
}
func Condition(condition cond, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}) bool {
	return condition.Condition(endpoint, info, req)
}

func handler(condition cond, interceptor client.Interceptor, ctx context.Context, channel clientv1.Channel, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}, handler clientv1.InterceptorHandler) (interface{}, error) {
	if condition.Condition(endpoint, info, req) {
		return handler(ctx, channel, endpoint, req)
	}
	return interceptor(ctx, channel, endpoint, info, req, handler)
}
