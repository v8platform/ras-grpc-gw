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

func (r *ClientsRepository) table() (db *pudge.Db, err error) {

	return r.db.Table("clients")

}

func (r *ClientsRepository) get(key interface{}, value interface{}) error {

	var (
		table *pudge.Db
		err   error
	)

	if table, err = r.table(); err != nil {
		return err
	}

	return table.Get(key, value)
}

func (r *ClientsRepository) set(key interface{}, value interface{}) error {

	var (
		table *pudge.Db
		err   error
	)

	if table, err = r.table(); err != nil {
		return err
	}

	return table.Set(key, value)
}

func (r *ClientsRepository) delete(key interface{}) error {
	var (
		table *pudge.Db
		err   error
	)

	if table, err = r.table(); err != nil {
		return err
	}

	return table.Delete(key)
}

func (r *ClientsRepository) Fetch(ctx context.Context) (res []domain.Client, err error) {
	panic("implement me")
}

func (r *ClientsRepository) fetch(ctx context.Context, filter func(domain.Client) bool) ([]domain.Client, error) {
	var (
		table  *pudge.Db
		err    error
		result []domain.Client
	)

	if table, err = r.table(); err != nil {
		return nil, err
	}

	keys, err := table.Keys(nil, 0, 0, false)
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		var val domain.Client
		err := table.Get(key, val)
		if err != nil {
			return nil, err
		}

		if filter(val) {
			result = append(result, val)
		}
	}

	return result, nil
}

func (r *ClientsRepository) fetchOne(ctx context.Context, filter func(domain.Client) bool) (domain.Client, error) {

	var (
		table  *pudge.Db
		err    error
		result domain.Client
	)

	if table, err = r.table(); err != nil {
		return domain.Client{}, err
	}

	keys, err := table.Keys(nil, 0, 0, false)
	if err != nil {
		return domain.Client{}, err
	}

	for _, key := range keys {
		var val domain.Client
		err := table.Get(key, val)
		if err != nil {
			return domain.Client{}, err
		}

		if filter(val) {
			result = val
			break
		}
	}

	return result, nil
}

func (r *ClientsRepository) GetByUser(ctx context.Context, userId string) ([]domain.Client, error) {

	return nil, nil

}

func (r *ClientsRepository) GetByUUID(ctx context.Context, id string) (domain.Client, error) {
	var client domain.Client

	err := r.get(id, &client)
	if err != nil {
		return domain.Client{}, err
	}
	return client, nil
}

func (r *ClientsRepository) Update(ctx context.Context, client *domain.Client) error {
	err := r.set(client.UUID, *client)
	if err != nil {
		return err
	}
	return nil
}

func (r *ClientsRepository) Store(_ context.Context, client *domain.Client) error {

	err := r.set(client.UUID, *client)
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
