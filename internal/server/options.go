package server

import (
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/lestrrat-go/option"
	"google.golang.org/grpc"
)

type Option = option.Interface

type identHandler struct{}
type identInterceptor struct{}
type identOption struct{}

type HandlerOption interface {
	Option
	handlerOption()
}

type handlerOption struct {
	Option
}

func newHandlerOption(n interface{}, v interface{}) HandlerOption {
	return &handlerOption{option.New(n, v)}
}
func (*handlerOption) handlerOption() {}

func WithHandlers(h ...RegisterServerHandler) HandlerOption {
	return newHandlerOption(identHandler{}, h)
}

type UnaryInterceptorOption interface {
	Option
	interceptorOption()
}

type unaryInterceptorOption struct {
	Option
}

func newUnaryInterceptorOption(n interface{}, v interface{}) UnaryInterceptorOption {
	return &unaryInterceptorOption{option.New(n, v)}
}
func (*unaryInterceptorOption) interceptorOption() {}

func WithAuth(auth grpc_auth.AuthFunc) UnaryInterceptorOption {
	return newUnaryInterceptorOption(identInterceptor{}, grpc_auth.UnaryServerInterceptor(auth))
}

func WithInterceptor(unaryInterceptor grpc.UnaryServerInterceptor) UnaryInterceptorOption {
	return newUnaryInterceptorOption(identInterceptor{}, unaryInterceptor)
}

type ServerOption interface {
	Option
	serverOption()
}

type serverOption struct {
	Option
}

func newServerOption(n interface{}, v interface{}) ServerOption {
	return &serverOption{option.New(n, v)}
}
func (*serverOption) serverOption() {}

func WithChainInterceptor(unaryInterceptor ...grpc.UnaryServerInterceptor) ServerOption {
	return newServerOption(identOption{}, grpc.ChainUnaryInterceptor(unaryInterceptor...))
}
