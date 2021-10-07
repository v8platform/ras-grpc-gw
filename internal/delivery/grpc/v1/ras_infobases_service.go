package v1

import (
	"context"
	messagesv1 "github.com/v8platform/protos/gen/ras/messages/v1"
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

	return s.client.GetInfobases(ctx, request)
}

func (s *rasInfobasesServiceServer) GetInfobasesSummary(ctx context.Context, request *messagesv1.GetInfobasesSummaryRequest) (*messagesv1.GetInfobasesSummaryResponse, error) {
	return s.client.GetInfobasesSummary(ctx, request)
}

func (s *rasInfobasesServiceServer) GetInfobase(ctx context.Context, request *messagesv1.GetInfobaseInfoRequest) (*messagesv1.GetInfobaseInfoResponse, error) {
	return s.client.GetInfobase(ctx, request)
}

func (s *rasInfobasesServiceServer) CreateInfobase(ctx context.Context, request *messagesv1.CreateInfobaseRequest) (*messagesv1.CreateInfobaseResponse, error) {
	return s.client.CreateInfobase(ctx, request)
}

func (s *rasInfobasesServiceServer) UpdateInfobase(ctx context.Context, request *messagesv1.UpdateInfobaseRequest) (*emptypb.Empty, error) {
	return s.client.UpdateInfobase(ctx, request)
}

func (s *rasInfobasesServiceServer) DeleteInfobase(ctx context.Context, request *messagesv1.DropInfobaseRequest) (*emptypb.Empty, error) {
	return s.client.DropInfobase(ctx, request)
}

func (s *rasInfobasesServiceServer) UpdateInfobaseSummary(ctx context.Context, request *messagesv1.UpdateInfobaseSummaryRequest) (*emptypb.Empty, error) {
	return s.client.UpdateInfobaseSummary(ctx, request)
}

func newInfobasesServiceServer(services *service.Services, client client.Client) apiv1.InfobasesServer {
	return &rasInfobasesServiceServer{
		services: services,
		client:   client,
	}
}
