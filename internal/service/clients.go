package service

import (
	"context"
	"github.com/lithammer/shortuuid/v3"
	"github.com/v8platform/ras-grpc-gw/internal/domain"
	"github.com/v8platform/ras-grpc-gw/internal/repository"
	"github.com/v8platform/ras-grpc-gw/pkg/cache"
)

// ClientsService реализует бизнес-логику работы
type ClientsService interface {
	Register(ctx context.Context, client *domain.Client) (string, error)
	GetByUUID(ctx context.Context, uuid string) (*domain.Client, error)
	Remove(ctx context.Context, uuid string) error
}

type clientService struct {
	services *Services
	r        repository.Clients
	cache    cache.Cache
}

func (c clientService) Register(ctx context.Context, client *domain.Client) (string, error) {

	if len(client.UUID) == 0 {
		client.UUID = shortuuid.New()
	}

	err := c.r.Store(ctx, client)

	if err != nil {
		return "", err
	}

	return client.UUID, nil

}

func (c clientService) GetByUUID(ctx context.Context, uuid string) (*domain.Client, error) {
	client, err := c.r.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (c clientService) Remove(ctx context.Context, uuid string) error {
	return c.r.Delete(ctx, uuid)
}

func NewClientService(clients repository.Clients, cache cache.Cache, manager *Services) ClientsService {
	return &clientService{
		r:        clients,
		cache:    cache,
		services: manager,
	}
}
