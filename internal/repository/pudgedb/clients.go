package pudgedb

import (
	"context"
	"github.com/v8platform/ras-grpc-gw/internal/database/pudgedb"
	"github.com/v8platform/ras-grpc-gw/internal/domain"
)

type ClientsRepository struct {
	db *pudgedb.Db
}

func (c ClientsRepository) Fetch(ctx context.Context) (res []domain.Client, err error) {
	panic("implement me")
}

func (c ClientsRepository) GetByUUID(ctx context.Context, id string) (domain.Client, error) {
	panic("implement me")
}

func (c ClientsRepository) Update(ctx context.Context, ar *domain.Client) error {
	panic("implement me")
}

func (c ClientsRepository) Store(ctx context.Context, a *domain.Client) error {
	panic("implement me")
}

func (c ClientsRepository) Delete(ctx context.Context, id string) error {
	panic("implement me")
}

func NewClientsRepository(db *pudgedb.Db) *ClientsRepository {
	return &ClientsRepository{
		db: db,
	}
}
