package v1

import (
	"context"
	context2 "github.com/v8platform/ras-grpc-gw/internal/context"
	"github.com/v8platform/ras-grpc-gw/internal/service"
	"github.com/v8platform/ras-grpc-gw/pkg/gen/service/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ClientServerService interface {
	apiv1.ApplicationsServiceServer
}

type clientServerService struct {
	apiv1.UnimplementedApplicationsServiceServer
	services *service.Services
}

func (c clientServerService) Register(ctx context.Context, request *apiv1.RegisterRequest) (*apiv1.RegisterResponse, error) {

	user, ok := context2.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "unknow user")
	}

	client, err := c.services.Users.RegisterClient(ctx, user.UUID, request.GetHost(), request.GetUuid())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &apiv1.RegisterResponse{
		Uuid: client.UUID,
	}, nil
}

func (c clientServerService) Get(ctx context.Context, request *apiv1.GetClientRequest) (*apiv1.Tokens, error) {

	tokens, err := c.services.Tokens.Get(ctx, request.GetUuid())
	if err != nil {
		return nil, err
	}

	return &apiv1.Tokens{
		AccessToken:  string(tokens.Access),
		RefreshToken: string(tokens.Refresh),
	}, nil
}

func (c clientServerService) Refresh(ctx context.Context, request *apiv1.RefreshRequest) (*apiv1.Tokens, error) {
	panic("implement me")
}

func (c clientServerService) GetClients(ctx context.Context, request *apiv1.GetClientsRequest) (*apiv1.GetClientsResponse, error) {
	panic("implement me")
}

func (c clientServerService) ResetClients(ctx context.Context, request *apiv1.ResetClientsRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func NewClientServerService(services *service.Services) ClientServerService {
	return &clientServerService{
		services: services,
	}
}
