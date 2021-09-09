package pudgedb

import (
	"context"
	"github.com/v8platform/ras-grpc-gw/internal/database/pudgedb"
	"github.com/v8platform/ras-grpc-gw/internal/domain"
)

type UsersRepository struct {
	db *pudgedb.Db
}

func (u UsersRepository) GetByID(ctx context.Context, id int32) (*domain.User, error) {
	panic("implement me")
}

func (u UsersRepository) GetByCredentials(ctx context.Context, user string, passwordHash string) (*domain.User, error) {
	panic("implement me")
}

func (u UsersRepository) Fetch(ctx context.Context) (res []domain.User, err error) {
	panic("implement me")
}

func (u UsersRepository) Update(ctx context.Context, cal *domain.User) error {
	panic("implement me")
}

func (u UsersRepository) Store(ctx context.Context, cal *domain.User) error {
	panic("implement me")
}

func (u UsersRepository) Delete(ctx context.Context, id string) error {
	panic("implement me")
}

func NewUsersRepository(db *pudgedb.Db) *UsersRepository {
	return &UsersRepository{
		db: db,
	}
}
