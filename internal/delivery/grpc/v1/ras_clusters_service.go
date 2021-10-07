package v1

import (
	"context"
	messagesv1 "github.com/v8platform/protos/gen/ras/messages/v1"
	appCtx "github.com/v8platform/ras-grpc-gw/internal/context"
	"github.com/v8platform/ras-grpc-gw/internal/service"
	apiv1 "github.com/v8platform/ras-grpc-gw/pkg/gen/service/api/v1"
	client "github.com/v8platform/ras-grpc-gw/pkg/ras_client"
)

type rasClustersServiceServer struct {
	apiv1.UnimplementedClustersServer
	services *service.Services
	client   client.Client
}

func newClustersServiceServer(services *service.Services, client client.Client) apiv1.ClustersServer {
	return &rasClustersServiceServer{
		services: services,
		client:   client,
	}
}

func (s *rasClustersServiceServer) GetClusters(ctx context.Context, request *messagesv1.GetClustersRequest) (*messagesv1.GetClustersResponse, error) {
	opts, _ := appCtx.RequestOptsFromContext(ctx)
	return s.client.GetClusters(ctx, request, opts...)
}

func (s *rasClustersServiceServer) GetClusterInfo(ctx context.Context, request *messagesv1.GetClusterInfoRequest) (*messagesv1.GetClusterInfoResponse, error) {
	return s.client.GetClusterInfo(ctx, request)
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
