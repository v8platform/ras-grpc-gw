package pudgedb

import (
	"context"
	"github.com/v8platform/ras-grpc-gw/internal/database/pudgedb"
	"github.com/v8platform/ras-grpc-gw/internal/domain"
)

type ClientsRepository struct {
	db *pudgedb.Db
}

func (r *ClientsRepository) table() (*pudgedb.Table, error) {

	return r.db.Table("clients")

}

func (r *ClientsRepository) get(key interface{}, value interface{}) error {

	var (
		table *pudgedb.Table
		err   error
	)

	if table, err = r.table(); err != nil {
		return err
	}

	return table.Get(key, value)
}

func (r *ClientsRepository) set(key interface{}, value interface{}) error {

	var (
		table *pudgedb.Table
		err   error
	)

	if table, err = r.table(); err != nil {
		return err
	}

	return table.Set(key, value)
}

func (r *ClientsRepository) delete(key interface{}) error {
	var (
		table *pudgedb.Table
		err   error
	)

	if table, err = r.table(); err != nil {
		return err
	}

	return table.Delete(key)
}

func (r *ClientsRepository) Fetch(_ context.Context) (res []domain.Client, err error) {
	panic("implement me")
}

func (r *ClientsRepository) GetByUser(_ context.Context, user domain.User) ([]domain.Client, error) {

	var clients []domain.Client

	for _, clientId := range user.Clients {

		var client domain.Client

		err := r.get(clientId, &client)
		if err != nil {
			continue
		}

		clients = append(clients, client)

	}

	return clients, nil

}

func (r *ClientsRepository) GetByID(ctx context.Context, id string) (domain.Client, error) {
	var client domain.Client

	err := r.get(id, &client)
	if err != nil {
		return domain.Client{}, err
	}
	return client, nil
}

func (r *ClientsRepository) Update(ctx context.Context, client domain.Client) error {
	err := r.set(client.UUID, client)
	if err != nil {
		return err
	}
	return nil
}

func (r *ClientsRepository) Store(_ context.Context, client domain.Client) error {

	err := r.set(client.UUID, client)
	if err != nil {
		return err
	}
	return nil
}

func (r *ClientsRepository) Delete(ctx context.Context, id string) error {
	err := r.delete(id)
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
