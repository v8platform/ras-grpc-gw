package cond

import (
	"context"
	clientv1 "github.com/v8platform/protos/gen/ras/client/v1"
)

type Func func(ctx context.Context, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}) bool

func (f Func) check(ctx context.Context, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}) bool {
	return f(ctx, endpoint, info, req)
}
func (f Func) Handler(interceptor ...func(ctx context.Context, channel clientv1.Channel, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}, handler clientv1.InterceptorHandler) (interface{}, error)) func(ctx context.Context, channel clientv1.Channel, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}, handler clientv1.InterceptorHandler) (interface{}, error) {
	return build(f, interceptor...)
}
func (f Func) IsFalse(ctx context.Context, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}) bool {
	return !f.check(ctx, endpoint, info, req)
}

func (f Func) IsTrue(ctx context.Context, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}) bool {
	return f.check(ctx, endpoint, info, req)
}

type checkCond interface {
	check(ctx context.Context, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}) bool
}

type cond interface {
	checkCond
	Handler(interceptor ...func(ctx context.Context, channel clientv1.Channel, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}, handler clientv1.InterceptorHandler) (interface{}, error)) func(ctx context.Context, channel clientv1.Channel, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}, handler clientv1.InterceptorHandler) (interface{}, error)
	IsTrue(ctx context.Context, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}) bool
	IsFalse(ctx context.Context, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}) bool
}

type Or []cond

func (c Or) check(ctx context.Context, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}) bool {
	for _, cond := range c {
		if cond.check(ctx, endpoint, info, req) {
			return true
		}
	}
	return false
}

func (c Or) Handler(interceptor ...func(ctx context.Context, channel clientv1.Channel, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}, handler clientv1.InterceptorHandler) (interface{}, error)) func(ctx context.Context, channel clientv1.Channel, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}, handler clientv1.InterceptorHandler) (interface{}, error) {
	return build(c, interceptor...)
}

func (c Or) IsFalse(ctx context.Context, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}) bool {
	return !c.check(ctx, endpoint, info, req)
}

func (c Or) IsTrue(ctx context.Context, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}) bool {
	return c.check(ctx, endpoint, info, req)
}

type Xor []cond

func (c Xor) check(ctx context.Context, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}) bool {
	for _, cond := range c {
		if cond.check(ctx, endpoint, info, req) {
			return false
		}
	}
	return true
}
func (c Xor) Handler(interceptor ...func(ctx context.Context, channel clientv1.Channel, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}, handler clientv1.InterceptorHandler) (interface{}, error)) func(ctx context.Context, channel clientv1.Channel, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}, handler clientv1.InterceptorHandler) (interface{}, error) {
	return build(c, interceptor...)
}

func (c Xor) IsFalse(ctx context.Context, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}) bool {
	return !c.check(ctx, endpoint, info, req)
}
func (c Xor) IsTrue(ctx context.Context, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}) bool {
	return c.check(ctx, endpoint, info, req)
}

type Not [1]cond

func (c Not) check(ctx context.Context, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}) bool {
	for _, cond := range c {
		if cond.check(ctx, endpoint, info, req) {
			return false
		}
	}
	return true
}

func (c Not) Handler(interceptor ...func(ctx context.Context, channel clientv1.Channel, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}, handler clientv1.InterceptorHandler) (interface{}, error)) func(ctx context.Context, channel clientv1.Channel, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}, handler clientv1.InterceptorHandler) (interface{}, error) {
	return build(c, interceptor...)
}

func (c Not) IsFalse(ctx context.Context, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}) bool {
	return !c.check(ctx, endpoint, info, req)
}
func (c Not) IsTrue(ctx context.Context, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}) bool {
	return c.check(ctx, endpoint, info, req)
}

type And []cond

func (c And) check(ctx context.Context, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}) bool {
	for _, cond := range c {
		if !cond.check(ctx, endpoint, info, req) {
			return false
		}
	}
	return true
}

func (c And) Handler(interceptor ...func(ctx context.Context, channel clientv1.Channel, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}, handler clientv1.InterceptorHandler) (interface{}, error)) func(ctx context.Context, channel clientv1.Channel, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}, handler clientv1.InterceptorHandler) (interface{}, error) {
	return build(c, interceptor...)
}

func (c And) IsFalse(ctx context.Context, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}) bool {
	return !c.check(ctx, endpoint, info, req)
}
func (c And) IsTrue(ctx context.Context, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}) bool {
	return c.check(ctx, endpoint, info, req)
}

func build(condition cond, interceptor ...func(ctx context.Context, channel clientv1.Channel, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}, handler clientv1.InterceptorHandler) (interface{}, error)) func(ctx context.Context, channel clientv1.Channel, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}, handler clientv1.InterceptorHandler) (interface{}, error) {
	n := len(interceptor)

	return func(ctx context.Context, channel clientv1.Channel, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}, handler clientv1.InterceptorHandler) (interface{}, error) {
		if condition.IsFalse(ctx, endpoint, info, req) {
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
