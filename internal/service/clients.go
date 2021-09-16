package service

import (
	"context"
	"github.com/lithammer/shortuuid/v3"
	"github.com/v8platform/ras-grpc-gw/internal/domain"
	"github.com/v8platform/ras-grpc-gw/internal/repository"
	"github.com/v8platform/ras-grpc-gw/pkg/cache"
)

// ApplicationsService реализует бизнес-логику работы
type ApplicationsService interface {
	Create(ctx context.Context, host string, name string) (domain.Application, error)
	GetByID(ctx context.Context, uuid string) (domain.Application, error)
	Remove(ctx context.Context, uuid string) error
	UpdateAuth(ctx context.Context, uuid string, user, password string) error
	Update(ctx context.Context, app domain.Application) (domain.Application, error)
}

type applicationService struct {
	services *Services
	r        repository.Clients
	cache    cache.Cache
}

func (c applicationService) Update(ctx context.Context, app domain.Application) (domain.Application, error) {
	err := c.r.Update(ctx, app)
	if err != nil {
		return domain.Application{}, err
	}
	// TODO Передалать
	newApp, err := c.r.GetByID(ctx, app.UUID)
	if err != nil {
		return domain.Application{}, err
	}
	return newApp, nil
}

func (c applicationService) UpdateAuth(ctx context.Context, uuid string, user, password string) error {

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

func (c applicationService) Create(ctx context.Context, host string, name string) (domain.Application, error) {

	client := domain.Application{
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

func (c applicationService) GetByID(ctx context.Context, uuid string) (domain.Application, error) {
	client, err := c.r.GetByID(ctx, uuid)
	if err != nil {
		return domain.Application{}, err
	}
	return client, nil
}

func (c applicationService) Remove(ctx context.Context, uuid string) error {
	return c.r.Delete(ctx, uuid)
}

func NewClientService(clients repository.Clients, cache cache.Cache, manager *Services) ApplicationsService {
	return &applicationService{
		r:        clients,
		cache:    cache,
		services: manager,
	}
}
