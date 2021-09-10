package pudgedb

import (
	"context"
	"github.com/lithammer/shortuuid/v3"
	"github.com/v8platform/ras-grpc-gw/internal/database/pudgedb"
	"github.com/v8platform/ras-grpc-gw/internal/domain"
	"github.com/v8platform/ras-grpc-gw/pkg/hash"
	"testing"
)

func TestUsersRepository_Store(t *testing.T) {
	type fields struct {
		db *pudgedb.Db
	}
	type args struct {
		ctx context.Context
		cal *domain.User
	}

	db := pudgedb.New("../../../pudgedb")
	hasher := hash.NewSHA1Hasher("salt")
	passwordHash, _ := hasher.Hash("pwd")
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"simple",
			fields{&db},
			args{
				ctx: context.Background(),
				cal: &domain.User{
					ID:           0,
					UUID:         shortuuid.New(),
					Username:     "admin",
					PasswordHash: passwordHash,
					Email:        "admin@mail",
					IsAdmin:      true,
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UsersRepository{
				db: tt.fields.db,
			}
			if err := u.Store(tt.args.ctx, tt.args.cal); (err != nil) != tt.wantErr {
				t.Errorf("Store() error = %v, wantErr %v", err, tt.wantErr)
			}
			var user *domain.User
			var err error
			if user, err = u.GetByCredentials(tt.args.ctx, tt.args.cal.Username, tt.args.cal.PasswordHash); (err != nil) != tt.wantErr {
				t.Errorf("GetByCredentials() error = %v, wantErr %v", err, tt.wantErr)
			}

			t.Log(user)

		})
	}
}
