package pudgedb

import (
	"context"
	"crypto/sha256"
	"github.com/recoilme/pudge"
	"github.com/v8platform/ras-grpc-gw/internal/database/pudgedb"
	"github.com/v8platform/ras-grpc-gw/internal/domain"
)

type UsersRepository struct {
	db *pudgedb.Db
}

func (u UsersRepository) pathUsers() string {
	return u.db.GetPath("users")
}

func (u UsersRepository) byCredentials() string {
	return u.db.GetPath("byCredentials")
}

func (u UsersRepository) GetByID(ctx context.Context, id int32) (*domain.User, error) {
	panic("implement me")
}

func (u UsersRepository) GetByCredentials(ctx context.Context, username string, passwordHash string) (*domain.User, error) {

	digest := sha256.Sum256([]byte(username + passwordHash))
	var key string

	err := pudge.Get(u.byCredentials(), digest, &key)
	if err != nil {
		return nil, err
	}

	var user domain.User
	err = pudge.Get(u.pathUsers(), key, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u UsersRepository) Fetch(ctx context.Context) (res []domain.User, err error) {

	keys, _ := pudge.Keys(u.pathUsers(), 0, -1, 0, true)
	for _, key := range keys {
		var p domain.User
		err := pudge.Get(u.pathUsers(), key, &p)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}

	return
}

func (u UsersRepository) Update(ctx context.Context, cal *domain.User) error {

	err := pudge.Set(u.pathUsers(), cal.UUID, *cal)
	if err != nil {
		return err
	}

	return nil
}

func (u UsersRepository) Store(ctx context.Context, cal *domain.User) error {

	err := pudge.Set(u.pathUsers(), cal.UUID, *cal)
	if err != nil {
		return err
	}
	digest := sha256.Sum256([]byte(cal.Username + cal.PasswordHash))

	err = pudge.Set(u.byCredentials(), digest, cal.UUID)
	if err != nil {
		return err
	}

	return nil
}

func (u UsersRepository) Delete(ctx context.Context, id string) error {

	err := pudge.Delete(u.pathUsers(), id)
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
