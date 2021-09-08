package domain

import "context"

type Client struct {
	UserID        string
	UUID          string
	Host          string
	AgentUser     string
	AgentPassword string
}

// ClientUsecase represent the usecases
type ClientUsecase interface {
	Fetch(ctx context.Context, userID int32) ([]Client, error)
	GetByUUID(ctx context.Context, uuid string) (Client, error)
	Update(ctx context.Context, client *Client) error
	Store(ctx context.Context, client *Client) error
	Delete(ctx context.Context, id string) error
}

// Clients represent the repository
type ClientRepository interface {
	Fetch(ctx context.Context) (res []Client, err error)
	GetByUUID(ctx context.Context, id string) (Client, error)
	Update(ctx context.Context, ar *Client) error
	Store(ctx context.Context, a *Client) error
	Delete(ctx context.Context, id string) error
}
