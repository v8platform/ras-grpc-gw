package v1

import (
	"context"

	clientv1 "github.com/v8platform/protos/gen/ras/client/v1"
	messagesv1 "github.com/v8platform/protos/gen/ras/messages/v1"
	"github.com/v8platform/ras-grpc-gw/internal/service"
	apiv1 "github.com/v8platform/ras-grpc-gw/pkg/gen/service/api/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type rasAuthServiceServer struct {
	apiv1.UnimplementedAuthServiceServer
	services *service.Services
	clients  ClientsStorage
}

func NewAuthServiceServer(services *service.Services, clients ClientsStorage) apiv1.AuthServiceServer {
	return &rasAuthServiceServer{
		services: services,
		clients:  clients,
	}
}

func (s *rasAuthServiceServer) AuthenticateCluster(ctx context.Context, request *messagesv1.ClusterAuthenticateRequest) (*emptypb.Empty, error) {

	endpoint, err := s.clients.GetEndpoint(ctx)
	if err != nil {
		return nil, err
	}
	auth := clientv1.NewAuthService(endpoint)

	return auth.AuthenticateCluster(ctx, request)

}

func (s *rasAuthServiceServer) AuthenticateInfobase(ctx context.Context, request *messagesv1.AuthenticateInfobaseRequest) (*emptypb.Empty, error) {

	endpoint, err := s.clients.GetEndpoint(ctx)
	if err != nil {
		return nil, err
	}
	auth := clientv1.NewAuthService(endpoint)

	return auth.AuthenticateInfobase(ctx, request)

}

func (s *rasAuthServiceServer) AuthenticateAgent(ctx context.Context, request *messagesv1.AuthenticateAgentRequest) (*emptypb.Empty, error) {
	endpoint, err := s.clients.GetEndpoint(ctx)
	if err != nil {
		return nil, err
	}
	auth := clientv1.NewAuthService(endpoint)

	return auth.AuthenticateAgent(ctx, request)
}
