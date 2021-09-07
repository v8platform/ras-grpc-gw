package domain

import "context"

type AccessToken string

type Token struct {
}

// TokenUsecase represent the article's usecases
type TokenUsecase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]Token, string, error)
	GetByID(ctx context.Context, id int64) (Token, error)
	Update(ctx context.Context, ar *Token) error
	GetByTitle(ctx context.Context, title string) (Token, error)
	Store(context.Context, *Token) error
	Delete(ctx context.Context, id int64) error
}

// TokenRepository represent the article's repository contract
type TokenRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []Token, nextCursor string, err error)
	GetByID(ctx context.Context, id int64) (Token, error)
	GetByTitle(ctx context.Context, title string) (Token, error)
	Update(ctx context.Context, ar *Token) error
	Store(ctx context.Context, a *Token) error
	Delete(ctx context.Context, id int64) error
}
