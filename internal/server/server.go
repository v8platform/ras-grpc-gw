package server

import (
	"context"
	"fmt"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"time"

	"github.com/grpc-ecosystem/go-grpc-prometheus"
	_ "github.com/lithammer/shortuuid/v3"
	"github.com/spf13/cast"
	clientv1 "github.com/v8platform/protos/gen/ras/client/v1"
	messagesv1 "github.com/v8platform/protos/gen/ras/messages/v1"
	ras_service "github.com/v8platform/protos/gen/ras/service/api/v1"
	access_service "github.com/v8platform/ras-grpc-gw/pkg/gen/access/service"
	"github.com/v8platform/ras-grpc-gw/pkg/ras_client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

import "github.com/grpc-ecosystem/go-grpc-middleware"

func NewRASServer(rasAddr string) *RASServer {
	return &RASServer{
		rasAddr: rasAddr,
	}
}

type RASServer struct {
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

func (s *RASServer) Serve(host string) error {

	listener, err := net.Listen("tcp", host)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", host, err)
	}

	srv := NewRasClientServiceServer()
	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			// grpc_auth.UnaryServerInterceptor(myAuthFunction),
		)),
	)
	ras_service.RegisterAuthServiceServer(server, srv)
	ras_service.RegisterClustersServiceServer(server, srv)
	ras_service.RegisterSessionsServiceServer(server, srv)
	ras_service.RegisterInfobasesServiceServer(server, srv)

	accessSrv := NewAccessServer()

	access_service.RegisterClientServiceServer(server, accessSrv)
	access_service.RegisterTokenServiceServer(server, accessSrv)

	log.Println("Listening on", host)
	if err := server.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve gRPC server: %w", err)
	}

	return nil
}

func NewRasClientServiceServer() ras_service.RASServiceServer {
	return &rasClientServiceServer{}
}

type rasClientServiceServer struct {
	ras_service.UnimplementedRASServiceServer

	idxClients   map[string]*ClientInfo
	idxEndpoints map[string]*EndpointInfo
}

func (s *rasClientServiceServer) endpointFromContext(ctx context.Context) (*EndpointInfo, error) {

	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return nil, status.Errorf(codes.DataLoss, "Server: failed to get metadata")
	}

	if t, ok := md["endpoint_id"]; ok {

		for _, e := range t {

			if endpoint, ok := s.idxEndpoints[e]; ok {
				return endpoint, nil
			}
		}
	}

	return nil, nil
}

func (s *rasClientServiceServer) clientFromContext(ctx context.Context) (*ClientInfo, error) {

	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return nil, status.Errorf(codes.DataLoss, "Server: failed to get metadata")
	}

	if t, ok := md["client_id"]; ok {

		for _, e := range t {

			if client, ok := s.idxClients[e]; ok {
				return client, nil
			}
		}
	}

	return nil, nil
}

func (s *rasClientServiceServer) getEndpoint(ctx context.Context) clientv1.EndpointServiceImpl {

	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return nil, status.Errorf(codes.DataLoss, "Client: failed to get metadata")
	}

	if t, ok := md["endpoint_id"]; ok {

		for _, e := range t {
			if endpoint, ok := c.getEndpoint(e); ok {
				return clientv1.NewEndpointService(c, endpoint), nil

			}
		}
	}

	endpoint, err := c.turnEndpoint(ctx)
	if err != nil {
		return nil, err
	}

}

func (s *rasClientServiceServer) AuthenticateCluster(ctx context.Context, request *messagesv1.ClusterAuthenticateRequest) (*emptypb.Empty, error) {
	endpoint, err := s.client.GetEndpoint(ctx)
	if err != nil {
		return nil, err
	}
	auth := clientv1.NewAuthService(endpoint)

	return auth.AuthenticateCluster(ctx, request)

}

func (s *rasClientServiceServer) withEndpoint(ctx context.Context, fn func(clientv1.EndpointServiceImpl) error) (err error) {
	var endpoint clientv1.EndpointServiceImpl
	endpoint, err = s.client.GetEndpoint(ctx)

	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			header := metadata.New(map[string]string{
				"endpoint_id": cast.ToString(endpoint),
				//"host": cast.ToString(s.ras_client.),
			})

			_ = grpc.SendHeader(ctx, header)
		}
	}()

	return fn(endpoint)
}

func (s *rasClientServiceServer) AuthenticateInfobase(ctx context.Context, request *messagesv1.AuthenticateInfobaseRequest) (*emptypb.Empty, error) {

	var resp *emptypb.Empty
	var err error

	err = s.withEndpoint(ctx, func(endpoint clientv1.EndpointServiceImpl) error {
		auth := clientv1.NewAuthService(endpoint)
		resp, err = auth.AuthenticateInfobase(ctx, request)
		if err != nil {
			return err
		}
		return nil

	})

	return resp, err
}

func (s *rasClientServiceServer) AuthenticateAgent(ctx context.Context, request *messagesv1.AuthenticateAgentRequest) (*emptypb.Empty, error) {
	endpoint, err := s.client.GetEndpoint(ctx)
	if err != nil {
		return nil, err
	}
	auth := clientv1.NewAuthService(endpoint)

	return auth.AuthenticateAgent(ctx, request)
}

func (s *rasClientServiceServer) GetClusters(ctx context.Context, request *messagesv1.GetClustersRequest) (*messagesv1.GetClustersResponse, error) {
	endpoint, err := s.client.GetEndpoint(ctx)
	if err != nil {
		return nil, err
	}
	service := clientv1.NewClustersService(endpoint)
	return service.GetClusters(ctx, request)
}

func (s *rasClientServiceServer) GetClusterInfo(ctx context.Context, request *messagesv1.GetClusterInfoRequest) (*messagesv1.GetClusterInfoResponse, error) {
	endpoint, err := s.client.GetEndpoint(ctx)
	if err != nil {
		return nil, err
	}
	service := clientv1.NewClustersService(endpoint)
	return service.GetClusterInfo(ctx, request)
}

func (s *rasClientServiceServer) GetSessions(ctx context.Context, request *messagesv1.GetSessionsRequest) (*messagesv1.GetSessionsResponse, error) {
	endpoint, err := s.client.GetEndpoint(ctx)
	if err != nil {
		return nil, err
	}

	service := clientv1.NewSessionsService(endpoint)
	return service.GetSessions(ctx, request)
}

func (s *rasClientServiceServer) GetShortInfobases(ctx context.Context, request *messagesv1.GetInfobasesShortRequest) (*messagesv1.GetInfobasesShortResponse, error) {
	endpoint, err := s.client.GetEndpoint(ctx)
	if err != nil {
		return nil, err
	}
	service := clientv1.NewInfobasesService(endpoint)
	return service.GetShortInfobases(ctx, request)
}

func (s *rasClientServiceServer) GetInfobaseSessions(ctx context.Context, request *messagesv1.GetInfobaseSessionsRequest) (*messagesv1.GetInfobaseSessionsResponse, error) {
	endpoint, err := s.client.GetEndpoint(ctx)
	if err != nil {
		return nil, err
	}
	service := clientv1.NewInfobasesService(endpoint)
	return service.GetSessions(ctx, request)
}
