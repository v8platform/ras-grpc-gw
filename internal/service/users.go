package service

import (
	"context"
	"github.com/lithammer/shortuuid/v3"
	"github.com/v8platform/ras-grpc-gw/internal/domain"
	"github.com/v8platform/ras-grpc-gw/internal/repository"
	"github.com/v8platform/ras-grpc-gw/pkg/cache"
	"github.com/v8platform/ras-grpc-gw/pkg/hash"
)

var _ UsersService = (*usersService)(nil)

// UsersService реализует бизнес-логику работы
type UsersService interface {
	GetByCredentials(ctx context.Context, user string, password string) (domain.User, error)
	GetUserClients(ctx context.Context, userId int32) ([]domain.Client, error)
	GetByUUID(ctx context.Context, uuid string) (domain.User, error)
	RegisterClient(ctx context.Context, userID string, host, name string) (domain.Client, error)
	Register(ctx context.Context, user string, password string) (domain.User, error)
}

type usersService struct {
	services *Services
	r        repository.Users
	cache    cache.Cache
	hasher   hash.PasswordHasher
}

func (u usersService) Register(ctx context.Context, username string, password string) (domain.User, error) {

	_, err := u.r.GetByName(ctx, username)
	if err == nil {
		return domain.User{}, domain.ErrUserExists
	}
	passwordHash, err := u.hasher.Hash(password)
	if err != nil {
		return domain.User{}, err
	}

	user := domain.User{
		UUID:         shortuuid.New(),
		Username:     username,
		PasswordHash: passwordHash,
		Email:        "",
		IsAdmin:      false,
		Clients:      []string{},
	}
	err = u.r.Store(ctx, user)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil

}

func (u usersService) RegisterClient(ctx context.Context, userID string, host, name string) (domain.Client, error) {

	client, err := u.services.Clients.Create(ctx, host, name)
	if err != nil {
		return client, err
	}

	_, err = u.r.AttachClient(ctx, userID, client.UUID)
	if err != nil {
		return domain.Client{}, err
	}

	return client, nil
}

func (u usersService) GetByUUID(ctx context.Context, uuid string) (domain.User, error) {

	user, err := u.r.GetByID(ctx, uuid)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u usersService) GetByCredentials(ctx context.Context, username string, password string) (domain.User, error) {

	passwordHash, err := u.hasher.Hash(password)
	if err != nil {
		return domain.User{}, err
	}

	user, err := u.r.GetByCredentials(ctx, username, passwordHash)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (u usersService) GetUserClients(ctx context.Context, userId int32) ([]domain.Client, error) {
	panic("implement me")
}

func NewUsersService(repository repository.Users, cache cache.Cache, hasher hash.PasswordHasher, manager *Services) *usersService {
	return &usersService{
		r:        repository,
		cache:    cache,
		services: manager,
		hasher:   hasher,
	}
}
