package domain

import "context"

type User struct {
	ID           int32
	Username     string
	PasswordHash string
	Email        string
	IsAdmin      bool
}

// UserUsecase represent the usecases
type UserUsecase interface {
	Fetch(ctx context.Context) ([]User, error)
	GetByID(ctx context.Context, id int32) (User, error)
	Update(ctx context.Context, val *User) error
	Store(ctx context.Context, val *User) error
	Delete(ctx context.Context, id string) error
}
