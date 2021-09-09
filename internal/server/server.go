package server

import (
	"fmt"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/v8platform/ras-grpc-gw/internal/service"
	"log"
	"net"
	"time"

	"github.com/grpc-ecosystem/go-grpc-prometheus"
	_ "github.com/lithammer/shortuuid/v3"
	"github.com/v8platform/ras-grpc-gw/pkg/ras_client"
	"google.golang.org/grpc"
)

import "github.com/grpc-ecosystem/go-grpc-middleware"

type RegisterServerHandler func(server *grpc.Server)

func NewServer(services *service.Services, handlers ...RegisterServerHandler) *Server {
	return &Server{
		services: services,
		handlers: handlers,
	}
}

type Server struct {
	services *service.Services
	handlers []RegisterServerHandler
}

type EndpointInfo struct {
	uuid       string
	client     *ClientInfo
	EndpointId string
}

type ClientInfo struct {
	uuid        string
	conn        *ras_client.ClientConn
	IdleTimeout time.Duration
}

func (s *Server) Serve(host string) error {

	listener, err := net.Listen("tcp", host)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", host, err)
	}

	// srv := NewRasClientServiceServer()
	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			// grpc_auth.UnaryServerInterceptor(myAuthFunction),
		)),
	)

	for _, handler := range s.handlers {
		handler(server)
	}
	// ras_service.RegisterAuthServiceServer(server, srv)
	// ras_service.RegisterClustersServiceServer(server, srv)
	// ras_service.RegisterSessionsServiceServer(server, srv)
	// ras_service.RegisterInfobasesServiceServer(server, srv)
	//
	// accessSrv := NewAccessServer()
	//
	// access_service.RegisterClientServiceServer(server, accessSrv)
	// access_service.RegisterTokenServiceServer(server, accessSrv)

	log.Println("Listening on", host)
	if err := server.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve gRPC server: %w", err)
	}

	return nil
}
