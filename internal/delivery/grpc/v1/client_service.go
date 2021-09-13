package v1

import (
	"context"
	context2 "github.com/v8platform/ras-grpc-gw/internal/context"
	service2 "github.com/v8platform/ras-grpc-gw/internal/service"
	"github.com/v8platform/ras-grpc-gw/pkg/gen/access/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ClientServerService interface {
	service.ClientServiceServer
}

type clientServerService struct {
	service.UnimplementedClientServiceServer
	services *service2.Services
}

func (c clientServerService) Register(ctx context.Context, request *service.RegisterRequest) (*service.RegisterResponse, error) {

	user, ok := context2.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "unknow user")
	}

	client, err := c.services.Users.RegisterClient(ctx, user.UUID, request.GetHost(), request.GetUuid())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &service.RegisterResponse{
		Uuid: client.UUID,
	}, nil
}

func (c clientServerService) Get(ctx context.Context, request *service.GetClientRequest) (*service.Tokens, error) {

	tokens, err := c.services.Tokens.Get(ctx, request.GetUuid())
	if err != nil {
		return nil, err
	}

	return &service.Tokens{
		AccessToken:  string(tokens.Access),
		RefreshToken: string(tokens.Refresh),
	}, nil
}

func (c clientServerService) Refresh(ctx context.Context, request *service.RefreshRequest) (*service.Tokens, error) {
	panic("implement me")
}

func (c clientServerService) GetClients(ctx context.Context, request *service.GetClientsRequest) (*service.GetClientsResponse, error) {
	panic("implement me")
}

func (c clientServerService) ResetClients(ctx context.Context, request *service.ResetClientsRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func NewClientServerService(services *service2.Services) ClientServerService {
	return &clientServerService{
		services: services,
	}
}
