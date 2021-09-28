package v1

import (
	"context"
	v1 "github.com/v8platform/protos/gen/ras/messages/v1"
	"github.com/v8platform/ras-grpc-gw/internal/service"
	apiv1 "github.com/v8platform/ras-grpc-gw/pkg/gen/service/api/v1"
	client "github.com/v8platform/ras-grpc-gw/pkg/ras_client"
	"google.golang.org/protobuf/types/known/emptypb"
)

type rasSessionsServiceServer struct {
	apiv1.UnimplementedSessionsServer
	services *service.Services
	cc       client.Client
}

func (r rasSessionsServiceServer) GetSessions(ctx context.Context, request *v1.GetSessionsRequest) (*v1.GetSessionsResponse, error) {
	return r.cc.GetSessions(ctx, request)
}

func (r rasSessionsServiceServer) GetInfobaseSessions(ctx context.Context, request *v1.GetInfobaseSessionsRequest) (*v1.GetInfobaseSessionsResponse, error) {
	return r.cc.GetInfobaseSessions(ctx, request)
}

func (r rasSessionsServiceServer) TerminateSession(ctx context.Context, request *v1.TerminateSessionRequest) (*emptypb.Empty, error) {
	return r.cc.TerminateSession(ctx, request)
}

func NewSessionsServiceServer(services *service.Services, cc client.Client) apiv1.SessionsServer {
	return &rasSessionsServiceServer{
		services: services,
		cc:       cc,
	}
}
