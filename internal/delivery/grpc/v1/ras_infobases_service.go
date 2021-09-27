package v1

import (
	"context"
	messagesv1 "github.com/v8platform/protos/gen/ras/messages/v1"
	appCtx "github.com/v8platform/ras-grpc-gw/internal/context"
	"github.com/v8platform/ras-grpc-gw/internal/service"
	apiv1 "github.com/v8platform/ras-grpc-gw/pkg/gen/service/api/v1"
	client "github.com/v8platform/ras-grpc-gw/pkg/ras_client"
	"google.golang.org/protobuf/types/known/emptypb"
)

type rasInfobasesServiceServer struct {
	apiv1.UnimplementedInfobasesServer
	services *service.Services
	client   client.Client
}

func (s *rasInfobasesServiceServer) GetInfobases(ctx context.Context, request *messagesv1.GetInfobasesRequest) (*messagesv1.GetInfobasesResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s *rasInfobasesServiceServer) GetInfobasesSummary(ctx context.Context, request *messagesv1.GetInfobasesSummaryRequest) (*messagesv1.GetInfobasesSummaryResponse, error) {
	opts, _ := appCtx.RequestOptsFromContext(ctx)
	return s.client.GetInfobasesSummary(ctx, request, opts...)

}

func (s *rasInfobasesServiceServer) GetInfobase(ctx context.Context, request *messagesv1.GetInfobaseInfoRequest) (*messagesv1.GetInfobaseInfoResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s *rasInfobasesServiceServer) CreateInfobase(ctx context.Context, request *messagesv1.CreateInfobaseRequest) (*messagesv1.CreateInfobaseResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s *rasInfobasesServiceServer) UpdateInfobase(ctx context.Context, request *messagesv1.UpdateInfobaseRequest) (*emptypb.Empty, error) {
	// TODO implement me
	panic("implement me")
}

func (s *rasInfobasesServiceServer) DeleteInfobase(ctx context.Context, request *messagesv1.DropInfobaseRequest) (*emptypb.Empty, error) {
	// TODO implement me
	panic("implement me")
}

func (s *rasInfobasesServiceServer) UpdateInfobaseSummary(ctx context.Context, request *messagesv1.UpdateInfobaseSummaryRequest) (*emptypb.Empty, error) {
	// TODO implement me
	panic("implement me")
}

func NewInfobasesServiceServer(services *service.Services, client client.Client) apiv1.InfobasesServer {
	return &rasInfobasesServiceServer{
		services: services,
		client:   client,
	}
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
