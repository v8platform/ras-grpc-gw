package v1

import (
	"context"
	v1 "github.com/v8platform/protos/gen/ras/messages/v1"
	"github.com/v8platform/ras-grpc-gw/internal/service"
	apiv1 "github.com/v8platform/ras-grpc-gw/pkg/gen/service/api/v1"
	client "github.com/v8platform/ras-grpc-gw/pkg/ras_client"
	"google.golang.org/protobuf/types/known/emptypb"
)

type rasWorkingServersServiceServer struct {
	apiv1.UnimplementedServersServer
	services *service.Services
	cc       client.Client
}

func (r rasWorkingServersServiceServer) GetWorkingServers(ctx context.Context, request *v1.GetWorkingServersRequest) (*v1.GetWorkingServersResponse, error) {
	return r.cc.GetWorkingServers(ctx, request)
}

func (r rasWorkingServersServiceServer) GetWorkingServerInfo(ctx context.Context, request *v1.GetWorkingServerInfoRequest) (*v1.GetWorkingServerInfoResponse, error) {
	return r.cc.GetWorkingServerInfo(ctx, request)
}

func (r rasWorkingServersServiceServer) AddWorkingServer(ctx context.Context, request *v1.AddWorkingServerRequest) (*v1.AddWorkingServerResponse, error) {
	return r.cc.AddWorkingServer(ctx, request)
}

func (r rasWorkingServersServiceServer) DeleteWorkingServer(ctx context.Context, request *v1.DeleteWorkingServerRequest) (*emptypb.Empty, error) {
	return r.cc.DeleteWorkingServer(ctx, request)
}

func newWorkingServerServiceServer(services *service.Services, cc client.Client) apiv1.ServersServer {
	return &rasWorkingServersServiceServer{
		services: services,
		cc:       cc,
	}
}
