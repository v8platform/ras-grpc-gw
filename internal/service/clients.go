package service

import (
	"context"
	"github.com/v8platform/ras-grpc-gw/internal/domain"
	"github.com/v8platform/ras-grpc-gw/internal/repository"
	"github.com/v8platform/ras-grpc-gw/pkg/cache"
)

// ClientsService реализует бизнес-логику работы
type ClientsService interface {
	Register(ctx context.Context, client *domain.Client) error
	GetByShortUUID(ctx context.Context, uuid string) (*domain.Client, error)
	Remove(ctx context.Context, client *domain.Client) error
	GetUserClients(ctx context.Context, userId int32) ([]*domain.Client, error)
}

type clientService struct {
	r     repository.Clients
	cache cache.Cache
}

func (c clientService) GetUserClients(ctx context.Context, userId int32) ([]*domain.Client, error) {
	panic("implement me")
}

func (c clientService) Register(ctx context.Context, client *domain.Client) error {
	panic("implement me")
}

func (c clientService) GetByShortUUID(ctx context.Context, uuid string) (*domain.Client, error) {
	panic("implement me")
}

func (c clientService) Remove(ctx context.Context, client *domain.Client) error {
	panic("implement me")
}

func NewClientService(clients repository.Clients, cache cache.Cache) ClientsService {
	return &clientService{
		r:     clients,
		cache: cache,
	}
}
