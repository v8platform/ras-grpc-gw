package v1

import (
	"context"
	"github.com/v8platform/ras-grpc-gw/internal/domain"
	service2 "github.com/v8platform/ras-grpc-gw/internal/service"
	"github.com/v8platform/ras-grpc-gw/pkg/gen/access/service"
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

	// TODO Получить user_id из контекста
	userId := ""

	client := &domain.Client{
		UserID:        userId,
		UUID:          request.GetUuid(),
		Host:          request.GetHost(),
		AgentUser:     "",
		AgentPassword: "",
	}

	uuid, err := c.services.Clients.Register(ctx, client)
	if err != nil {
		return nil, err
	}

	return &service.RegisterResponse{
		Uuid: uuid,
	}, nil
}

func (c clientServerService) Get(ctx context.Context, request *service.GetClientRequest) (*service.Tokens, error) {
	panic("implement me")
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
