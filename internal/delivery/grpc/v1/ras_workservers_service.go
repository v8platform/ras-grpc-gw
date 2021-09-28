package v1

import (
	"context"
	v1 "github.com/v8platform/protos/gen/ras/messages/v1"
	"github.com/v8platform/ras-grpc-gw/internal/service"
	apiv1 "github.com/v8platform/ras-grpc-gw/pkg/gen/service/api/v1"
	client "github.com/v8platform/ras-grpc-gw/pkg/ras_client"
	"google.golang.org/protobuf/types/known/emptypb"
)

type rasWorkServersServiceServer struct {
	apiv1.UnimplementedServersServer
	services *service.Services
	cc       client.Client
}

func (r rasWorkServersServiceServer) GetWorkingServers(ctx context.Context, request *v1.GetWorkingServersRequest) (*v1.GetWorkingServersResponse, error) {
	return r.cc.GetWorkingServers(ctx, request)
}

func (r rasWorkServersServiceServer) GetWorkingServerInfo(ctx context.Context, request *v1.GetWorkingServerInfoRequest) (*v1.GetWorkingServerInfoResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (r rasWorkServersServiceServer) AddWorkingServer(ctx context.Context, request *v1.AddWorkingServerRequest) (*v1.AddWorkingServerResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (r rasWorkServersServiceServer) DeleteWorkingServer(ctx context.Context, request *v1.DeleteWorkingServerRequest) (*emptypb.Empty, error) {
	// TODO implement me
	panic("implement me")
}

func NewWorkServersServiceServer(services *service.Services, cc client.Client) apiv1.ServersServer {
	return &rasWorkServersServiceServer{
		services: services,
		cc:       cc,
	}
}
