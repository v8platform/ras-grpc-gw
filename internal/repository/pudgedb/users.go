package pudgedb

import (
	"context"
	"github.com/v8platform/ras-grpc-gw/internal/database/pudgedb"
	"github.com/v8platform/ras-grpc-gw/internal/domain"
	"strings"
)

type UsersRepository struct {
	db *pudgedb.Db
}

func (r *UsersRepository) GetByName(_ context.Context, username string) (domain.User, error) {
	table, err := r.table()
	if err != nil {
		return domain.User{}, err
	}

	keys, err := table.Keys(nil, 0, 0, false)
	if err != nil {
		return domain.User{}, err
	}

	for _, key := range keys {
		var user domain.User
		err := table.Get(key, &user)
		if err != nil {
			continue
		}

		if strings.ToLower(user.Username) == strings.ToLower(username) {
			return user, nil
		}

	}

	return domain.User{}, domain.ErrUserNotFound
}

func (r *UsersRepository) AttachClient(_ context.Context, userId string, clientId string) (domain.User, error) {

	table, err := r.table()
	if err != nil {
		return domain.User{}, err
	}
	var user domain.User

	if err := table.Get(userId, &user); err == nil {
		return user, nil
	}

	user.Clients = append(user.Clients, clientId)

	err = table.Set(user.UUID, user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *UsersRepository) table() (*pudgedb.Table, error) {

	return r.db.Table("users")

}

func (r *UsersRepository) GetByID(_ context.Context, id string) (domain.User, error) {
	table, err := r.table()
	if err != nil {
		return domain.User{}, err
	}
	var user domain.User

	if err := table.Get(id, &user); err == nil {
		return user, nil
	}

	return domain.User{}, domain.ErrUserNotFound
}

func (r *UsersRepository) GetByCredentials(_ context.Context, username string, passwordHash string) (domain.User, error) {

	table, err := r.table()
	if err != nil {
		return domain.User{}, err
	}

	keys, err := table.Keys(nil, 0, 0, false)
	if err != nil {
		return domain.User{}, err
	}

	for _, key := range keys {
		var user domain.User
		err := table.Get(key, &user)
		if err != nil {
			continue
		}

		if strings.ToLower(user.Username) == strings.ToLower(username) &&
			strings.EqualFold(user.PasswordHash, passwordHash) {
			return user, nil
		}

	}

	return domain.User{}, domain.ErrUserNotFound
}

func (r *UsersRepository) Fetch(_ context.Context) (res []domain.User, err error) {

	table, err := r.table()
	if err != nil {
		return nil, err
	}

	keys, err := table.Keys(nil, 0, 0, false)
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		var user domain.User
		err := table.Get(key, &user)
		if err != nil {
			continue
		}
		res = append(res, user)
	}

	return
}

func (r *UsersRepository) Update(_ context.Context, user domain.User) error {

	table, err := r.table()
	if err != nil {
		return err
	}

	err = table.Set(user.UUID, user)
	if err != nil {
		return err
	}

	return nil
}

func (r *UsersRepository) Store(_ context.Context, user domain.User) error {
	table, err := r.table()
	if err != nil {
		return err
	}

	err = table.Set(user.UUID, user)
	if err != nil {
		return err
	}
	return nil
}

func (r *UsersRepository) Delete(_ context.Context, id string) error {

	table, err := r.table()
	if err != nil {
		return err
	}

	err = table.Delete(id)
	if err != nil {
		return err
	}
	return nil

}

func NewUsersRepository(db *pudgedb.Db) *UsersRepository {
	return &UsersRepository{
		db: db,
	}
}
