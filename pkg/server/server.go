package server

import (
	"context"
	"fmt"
	clientv1 "github.com/v8platform/protos/gen/ras/client/v1"
	messagesv1 "github.com/v8platform/protos/gen/ras/messages/v1"
	protocolv1 "github.com/v8platform/protos/gen/ras/protocol/v1"
	ras_service "github.com/v8platform/protos/gen/ras/service/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"sync"
)

func NewRASServer(rasAddr string) *RASServer {
	return &RASServer{
		rasAddr: rasAddr,
	}
}

type RASServer struct {
	rasAddr string
}

func (s *RASServer) Serve(host string) error {

	listener, err := net.Listen("tcp", host)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", host, err)
	}

	srv := NewRasClientServiceServer(s.rasAddr)
	server := grpc.NewServer()
	ras_service.RegisterRASServiceServer(server, srv)

	log.Println("Listening on", host)
	if err := server.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve gRPC server: %w", err)
	}

	return nil
}

func NewRasClientServiceServer(rasAddr string) ras_service.RASServiceServer {
	return &rasClientServiceServer{
		Host:              rasAddr,
		ClientServiceImpl: clientv1.NewClientService(rasAddr),
		once:              &sync.Once{},
	}
}

type rasClientServiceServer struct {
	ras_service.UnimplementedRASServiceServer
	clientv1.ClientServiceImpl
	rasService clientv1.RasServiceImpl
	endpoint   clientv1.EndpointServiceImpl
	Host       string
	once       *sync.Once
}

func (s rasClientServiceServer) AuthenticateCluster(ctx context.Context, request *messagesv1.ClusterAuthenticateRequest) (*emptypb.Empty, error) {
	if err := s.initOnce(); err != nil {
		return nil, err
	}

	auth := clientv1.NewAuthService(s.endpoint)
	return auth.AuthenticateCluster(request)

}

func (s rasClientServiceServer) AuthenticateInfobase(ctx context.Context, request *messagesv1.AuthenticateInfobaseRequest) (*emptypb.Empty, error) {
	if err := s.initOnce(); err != nil {
		return nil, err
	}

	auth := clientv1.NewAuthService(s.endpoint)
	return auth.AuthenticateInfobase(request)
}

func (s rasClientServiceServer) AuthenticateAgent(ctx context.Context, request *messagesv1.AuthenticateAgentRequest) (*emptypb.Empty, error) {
	if err := s.initOnce(); err != nil {
		return nil, err
	}

	auth := clientv1.NewAuthService(s.endpoint)
	return auth.AuthenticateAgent(request)
}

func (s rasClientServiceServer) GetClusters(ctx context.Context, request *messagesv1.GetClustersRequest) (*messagesv1.GetClustersResponse, error) {
	if err := s.initOnce(); err != nil {
		return nil, err
	}

	service := clientv1.NewClustersService(s.endpoint)
	return service.GetClusters(request)
}

func (s rasClientServiceServer) GetClusterInfo(ctx context.Context, request *messagesv1.GetClusterInfoRequest) (*messagesv1.GetClusterInfoResponse, error) {
	if err := s.initOnce(); err != nil {
		return nil, err
	}

	service := clientv1.NewClustersService(s.endpoint)
	return service.GetClusterInfo(request)
}

func (s rasClientServiceServer) GetSessions(ctx context.Context, request *messagesv1.GetSessionsRequest) (*messagesv1.GetSessionsResponse, error) {
	if err := s.initOnce(); err != nil {
		return nil, err
	}

	service := clientv1.NewSessionsService(s.endpoint)
	return service.GetSessions(request)
}

func (s rasClientServiceServer) GetShortInfobases(ctx context.Context, request *messagesv1.GetInfobasesShortRequest) (*messagesv1.GetInfobasesShortResponse, error) {
	if err := s.initOnce(); err != nil {
		return nil, err
	}

	service := clientv1.NewInfobasesService(s.endpoint)
	return service.GetShortInfobases(request)
}

func (s *rasClientServiceServer) GetInfobaseSessions(ctx context.Context, request *messagesv1.GetInfobaseSessionsRequest) (*messagesv1.GetInfobaseSessionsResponse, error) {
	service := clientv1.NewInfobasesService(s.endpoint)
	return service.GetSessions(request)
}

func (s *rasClientServiceServer) initOnce() (err error) {
	s.once.Do(func() {
		err = s.init()
	})

	return err
}

func (s *rasClientServiceServer) init() error {

	_, err := s.Negotiate(protocolv1.NewNegotiateMessage())
	if err != nil {
		return err
	}

	_, err = s.Connect(&protocolv1.ConnectMessage{})
	if err != nil {
		return err
	}

	endpointOpenAck, err := s.EndpointOpen(&protocolv1.EndpointOpen{
		Service: "v8.service.Admin.Cluster",
		Version: "10.0",
	})

	if err != nil {
		supporsedVarsion := s.DetectSupportedVersion(err)
		if len(supporsedVarsion) > 0 {
			endpointOpenAck, err = s.EndpointOpen(&protocolv1.EndpointOpen{
				Service: "v8.service.Admin.Cluster",
				Version: supporsedVarsion,
			})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	endpoint, err := s.NewEndpoint(endpointOpenAck)
	if err != nil {
		return err
	}

	s.endpoint = clientv1.NewEndpointService(s, endpoint)

	return nil

}
