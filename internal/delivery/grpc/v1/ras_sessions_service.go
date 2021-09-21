package v1

import (
	"context"
	clientv1 "github.com/v8platform/protos/gen/ras/client/v1"
	v1 "github.com/v8platform/protos/gen/ras/messages/v1"
	"github.com/v8platform/ras-grpc-gw/internal/service"
	apiv1 "github.com/v8platform/ras-grpc-gw/pkg/gen/service/api/v1"
)

type rasSessionsServiceServer struct {
	apiv1.UnimplementedSessionsServiceServer
	services *service.Services
	clients  Client
}

func (r rasSessionsServiceServer) GetSessions(ctx context.Context, request *v1.GetSessionsRequest) (*v1.GetSessionsResponse, error) {
	endpoint, err := r.clients.GetEndpoint(ctx)
	if err != nil {
		return nil, err
	}
	s := clientv1.NewSessionsService(endpoint)
	return s.GetSessions(ctx, request)
}

func (r rasSessionsServiceServer) GetInfobaseSessions(ctx context.Context, request *v1.GetInfobaseSessionsRequest) (*v1.GetInfobaseSessionsResponse, error) {
	endpoint, err := r.clients.GetEndpoint(ctx)
	if err != nil {
		return nil, err
	}
	s := clientv1.NewInfobasesService(endpoint)
	return s.GetSessions(ctx, request)
}

func NewSessionsServiceServer(services *service.Services, clients Client) apiv1.SessionsServiceServer {
	return &rasSessionsServiceServer{
		services: services,
		clients:  clients,
	}
}
