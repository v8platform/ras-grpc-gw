package pudgedb

import (
	"context"
	"github.com/recoilme/pudge"
	"github.com/v8platform/ras-grpc-gw/internal/database/pudgedb"
	"github.com/v8platform/ras-grpc-gw/internal/domain"
)

type ClientsRepository struct {
	db *pudgedb.Db
}

func (u ClientsRepository) clients() string {
	return u.db.GetPath("clients")
}

func (c ClientsRepository) Fetch(ctx context.Context) (res []domain.Client, err error) {
	panic("implement me")
}

func (c ClientsRepository) GetByUUID(ctx context.Context, id string) (domain.Client, error) {
	var client domain.Client
	err := pudge.Get(c.clients(), id, &client)
	if err != nil {
		return domain.Client{}, err
	}
	return client, nil
}

func (c ClientsRepository) Update(ctx context.Context, client *domain.Client) error {
	err := pudge.Set(c.clients(), client.UUID, *client)
	if err != nil {
		return err
	}
	return nil
}

func (c ClientsRepository) Store(_ context.Context, client *domain.Client) error {

	err := pudge.Set(c.clients(), client.UUID, *client)
	if err != nil {
		return err
	}
	return nil
}

func (c ClientsRepository) Delete(ctx context.Context, id string) error {
	err := pudge.Delete(c.clients(), id)
	if err != nil {
		return err
	}
	return nil
}

func NewClientsRepository(db *pudgedb.Db) *ClientsRepository {
	return &ClientsRepository{
		db: db,
	}
}
