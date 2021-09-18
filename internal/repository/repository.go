package repository

import (
	"context"
	"github.com/v8platform/ras-grpc-gw/internal/config"
	pudge "github.com/v8platform/ras-grpc-gw/internal/database/pudgedb"
	"github.com/v8platform/ras-grpc-gw/internal/domain"
	"github.com/v8platform/ras-grpc-gw/internal/repository/pudgedb"
)

type Users interface {
	Fetch(ctx context.Context) (res []domain.User, err error)
	GetByID(ctx context.Context, id string) (domain.User, error)
	GetByName(ctx context.Context, username string) (domain.User, error)
	GetByCredentials(ctx context.Context, username string, passwordHash string) (domain.User, error)
	Update(ctx context.Context, user domain.User) error
	Store(ctx context.Context, user domain.User) error
	Delete(ctx context.Context, id string) error
	AttachApplication(ctx context.Context, userId string, clientId string) (domain.User, error)
}

// Clients represent the repository
type Clients interface {
	Fetch(ctx context.Context) (res []domain.Application, err error)
	GetByUser(_ context.Context, user domain.User) ([]domain.Application, error)
	GetByID(ctx context.Context, id string) (domain.Application, error)
	Update(ctx context.Context, client domain.Application) error
	Store(ctx context.Context, client domain.Application) error
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

func CreateRepository(config config.DatabaseConfig) (*Repositories, error) {

	switch config.Engine.Name() {

	case "pudge":

		db, err := pudge.New(config.Engine.Config())
		if err != nil {
			return nil, err
		}

		return NewPudgeRepositories(db), nil

	case "postgres":
		panic("TODO Add support postgres")
	}

	return nil, nil
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
