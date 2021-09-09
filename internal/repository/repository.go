package repository

import (
	"context"
	pudge "github.com/v8platform/ras-grpc-gw/internal/database/pudgedb"
	"github.com/v8platform/ras-grpc-gw/internal/domain"
	"github.com/v8platform/ras-grpc-gw/internal/repository/pudgedb"
)

type Users interface {
	Fetch(ctx context.Context) (res []domain.User, err error)
	GetByID(ctx context.Context, id int32) (*domain.User, error)
	GetByCredentials(ctx context.Context, user string, passwordHash string) (*domain.User, error)
	Update(ctx context.Context, cal *domain.User) error
	Store(ctx context.Context, cal *domain.User) error
	Delete(ctx context.Context, id string) error
}

// Clients represent the repository
type Clients interface {
	Fetch(ctx context.Context) (res []domain.Client, err error)
	GetByUUID(ctx context.Context, id string) (domain.Client, error)
	Update(ctx context.Context, ar *domain.Client) error
	Store(ctx context.Context, a *domain.Client) error
	Delete(ctx context.Context, id string) error
}

var _ Clients = (*pudgedb.ClientsRepository)(nil)
var _ Users = (*pudgedb.UsersRepository)(nil)

func NewPudgeRepositories(db *pudge.Db) *Repositories {

	return NewRepositories(
		pudgedb.NewUsersRepository(db),
		pudgedb.NewClientsRepository(db),
	)

}

type Repositories struct {
	Users   Users
	Clients Clients
}

func NewRepositories(users Users, clients Clients) *Repositories {

	return &Repositories{
		Users:   users,
		Clients: clients,
	}
}
