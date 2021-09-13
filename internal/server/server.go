package server

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	_ "github.com/lithammer/shortuuid/v3"
	"github.com/v8platform/ras-grpc-gw/internal/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

type RegisterServerHandler func(server *grpc.Server)

func NewServer(opts ...Option) *Server {

	server := &Server{}

	for _, opt := range opts {

		switch opt.Ident() {
		case identInterceptor{}:
			server.unaryInterceptors = append(server.unaryInterceptors, opt.Value().(grpc.UnaryServerInterceptor))
		case identOption{}:
			server.options = append(server.options, opt.Value().(grpc.ServerOption))
		case identHandler{}:
			server.handlers = append(server.handlers, opt.Value().([]RegisterServerHandler)...)
		default:
			fmt.Printf("get unknown option %v", opt.Ident())
		}
	}

	return server
}

type Server struct {
	services          *service.Services
	handlers          []RegisterServerHandler
	unaryInterceptors []grpc.UnaryServerInterceptor
	options           []grpc.ServerOption
}

func (s *Server) Serve(host string) error {

	listener, err := net.Listen("tcp", host)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", host, err)
	}

	interceptors := []grpc.UnaryServerInterceptor{
		grpc_recovery.UnaryServerInterceptor(),
		grpc_prometheus.UnaryServerInterceptor,
	}

	for _, interceptor := range s.unaryInterceptors {
		interceptors = append(interceptors, interceptor)
	}

	s.options = append(s.options, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(interceptors...)))

	server := grpc.NewServer(s.options...)

	for _, handler := range s.handlers {
		handler(server)
	}

	go func() {
		log.Fatal(Run("dns:///" + host))
	}()

	log.Println("Listening on", host)
	if err := server.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve gRPC server: %w", err)
	}

	return nil
}
