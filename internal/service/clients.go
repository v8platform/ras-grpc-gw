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
	Create(ctx context.Context, host string, name string) (domain.Client, error)
	GetByID(ctx context.Context, uuid string) (domain.Client, error)
	Remove(ctx context.Context, uuid string) error
	UpdateAuth(ctx context.Context, uuid string, user, password string) error
}

type clientService struct {
	services *Services
	r        repository.Clients
	cache    cache.Cache
}

func (c clientService) UpdateAuth(ctx context.Context, uuid string, user, password string) error {

	client, err := c.r.GetByID(ctx, uuid)
	if err != nil {
		return err
	}

	client.AgentUser = user
	client.AgentPassword = password

	err = c.r.Update(ctx, client)

	if err != nil {
		return err
	}
	return nil

}

func (c clientService) Create(ctx context.Context, host string, name string) (domain.Client, error) {

	client := domain.Client{
		Host: host,
		Name: name,
		UUID: shortuuid.New(),
	}

	err := c.r.Store(ctx, client)

	if err != nil {
		return client, err
	}

	return client, nil

}

func (c clientService) GetByID(ctx context.Context, uuid string) (domain.Client, error) {
	client, err := c.r.GetByID(ctx, uuid)
	if err != nil {
		return domain.Client{}, err
	}
	return client, nil
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
