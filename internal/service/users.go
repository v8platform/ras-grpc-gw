package service

import (
	"context"
	"github.com/v8platform/ras-grpc-gw/internal/domain"
	"github.com/v8platform/ras-grpc-gw/internal/repository"
	"github.com/v8platform/ras-grpc-gw/pkg/cache"
	"github.com/v8platform/ras-grpc-gw/pkg/hash"
)

var _ UsersService = (*usersService)(nil)


// UsersService реализует бизнес-логику работы
type UsersService interface {
	GetByCredentials(ctx context.Context, user string, password string) (*domain.User, error)
	GetUserClients(ctx context.Context, userId int32) ([]*domain.Client, error)
}

type usersService struct {
	services *Services
	r  repository.Users
	cache cache.Cache
	hasher hash.PasswordHasher
}

func (u usersService) GetByCredentials(ctx context.Context, username string, password string) (*domain.User, error) {
	passwordHash, err := u.hasher.Hash(password)
	if err != nil {
		return nil, err
	}

	user, err := u.r.GetByCredentials(ctx, username, passwordHash)
	if err != nil {
	}
	return user, nil
}

func (u usersService) GetUserClients(ctx context.Context, userId int32) ([]*domain.Client, error) {
	panic("implement me")
}

func NewUsersService(repository repository.Users, cache cache.Cache, hasher hash.PasswordHasher, manager *Services) *usersService {
	return &usersService{
		r:  repository,
		cache: cache,
		services: manager,
		hasher: hasher,
	}
}
