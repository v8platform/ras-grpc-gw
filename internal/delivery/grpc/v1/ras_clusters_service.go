package v1

import (
	"context"
	clientv1 "github.com/v8platform/protos/gen/ras/client/v1"
	messagesv1 "github.com/v8platform/protos/gen/ras/messages/v1"
	ras_service "github.com/v8platform/protos/gen/ras/service/api/v1"
	"github.com/v8platform/ras-grpc-gw/internal/service"
)

//
// import (
// 	"context"
// 	"github.com/spf13/cast"
// 	clientv1 "github.com/v8platform/protos/gen/ras/client/v1"
// 	messagesv1 "github.com/v8platform/protos/gen/ras/messages/v1"
// 	ras_service "github.com/v8platform/protos/gen/ras/service/api/v1"
// 	"github.com/v8platform/ras-grpc-gw/internal/service"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/metadata"
// 	"google.golang.org/grpc/status"
// )

type rasClustersServiceServer struct {
	ras_service.UnimplementedClustersServiceServer
	services *service.Services
	clients  ClientsStorage
}

func NewRasClustersServer(services *service.Services, clients ClientsStorage) ras_service.ClustersServiceServer {
	return &rasClustersServiceServer{
		services: services,
		clients:  clients,
	}
}

func (s *rasClustersServiceServer) GetClusters(ctx context.Context, request *messagesv1.GetClustersRequest) (*messagesv1.GetClustersResponse, error) {
	endpoint, err := s.clients.GetEndpoint(ctx)
	if err != nil {
		return nil, err
	}
	clusters := clientv1.NewClustersService(endpoint)
	return clusters.GetClusters(ctx, request)
}

func (s *rasClustersServiceServer) GetClusterInfo(ctx context.Context, request *messagesv1.GetClusterInfoRequest) (*messagesv1.GetClusterInfoResponse, error) {
	endpoint, err := s.clients.GetEndpoint(ctx)
	if err != nil {
		return nil, err
	}
	clusters := clientv1.NewClustersService(endpoint)
	return clusters.GetClusterInfo(ctx, request)
}

// func (s *rasClientServiceServer) GetSessions(ctx context.Context, request *messagesv1.GetSessionsRequest) (*messagesv1.GetSessionsResponse, error) {
// 	endpoint, err := s.client.GetEndpoint(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	service := clientv1.NewSessionsService(endpoint)
// 	return service.GetSessions(ctx, request)
// }
//
// func (s *rasClientServiceServer) GetShortInfobases(ctx context.Context, request *messagesv1.GetInfobasesShortRequest) (*messagesv1.GetInfobasesShortResponse, error) {
// 	endpoint, err := s.client.GetEndpoint(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	service := clientv1.NewInfobasesService(endpoint)
// 	return service.GetShortInfobases(ctx, request)
// }
//
// func (s *rasClientServiceServer) GetInfobaseSessions(ctx context.Context, request *messagesv1.GetInfobaseSessionsRequest) (*messagesv1.GetInfobaseSessionsResponse, error) {
// 	endpoint, err := s.client.GetEndpoint(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	service := clientv1.NewInfobasesService(endpoint)
// 	return service.GetSessions(ctx, request)
// }
//
