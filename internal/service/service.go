package service

import (
	"github.com/v8platform/ras-grpc-gw/internal/repository"
	"github.com/v8platform/ras-grpc-gw/pkg/auth"
	"github.com/v8platform/ras-grpc-gw/pkg/cache"
	"github.com/v8platform/ras-grpc-gw/pkg/hash"
)

type builder struct {
	services []func(services *Services)
}

func (b *builder) ClientsService(clients repository.Clients, cache cache.Cache) {

	b.services = append(b.services, func(services *Services) {
		services.Clients = NewClientService(clients, cache, services)
	})
}

func (b *builder) UsersService(users repository.Users, cache cache.Cache, hasher hash.PasswordHasher) {

	b.services = append(b.services, func(services *Services) {
		services.Users = NewUsersService(users, cache, hasher, services)
	})
}

func (b *builder) TokenService(tokenManager auth.TokenManager) {

	b.services = append(b.services, func(services *Services) {
		services.Tokens = NewTokenService(tokenManager, services)
	})

}

func (b *builder) Build() (*Services, error) {

	services := &Services{}

	for _, initService := range b.services {
		initService(services)
	}

	err := services.checkServices()
	if err != nil {
		return nil, err
	}

	return services, nil
}

func NewServices(options Options) (*Services, error) {

	b := builder{}
	b.TokenService(options.TokenManager)
	b.UsersService(options.Repositories.Users, options.Cache, options.Hasher)
	b.ClientsService(options.Repositories.Clients, options.Cache)

	services, err := b.Build()
	if err != nil {
		return nil, err
	}

	services.Cache = options.Cache
	services.Hasher = options.Hasher
	services.TokenManager = options.TokenManager

	return services, nil
}

type Options struct {
	Repositories *repository.Repositories
	Cache        cache.Cache
	Hasher       hash.PasswordHasher
	TokenManager auth.TokenManager
}

type Services struct {
	Tokens  TokensService
	Clients ClientsService
	Users   UsersService

	Cache        cache.Cache
	Hasher       hash.PasswordHasher
	TokenManager auth.TokenManager
}

func (m *Services) checkServices() error {

	// TODO Сделать проверку что все сервисы инициализированы
	return nil

}
