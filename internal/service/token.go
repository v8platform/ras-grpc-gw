package service

import (
	"context"
	"github.com/v8platform/ras-grpc-gw/internal/domain"
	"github.com/v8platform/ras-grpc-gw/pkg/auth"
	"github.com/v8platform/ras-grpc-gw/pkg/cache"
)

type ValidateStatus int

// TokensService реализует бизнес-логику работы
type TokensService interface {
	Get(ctx context.Context, user *domain.User) (*domain.Token, error)
	Refresh(ctx context.Context, refresh *domain.Token) (*domain.Token, error)
	Validate(ctx context.Context, token *domain.Token) (ValidateStatus, error)
}

type tokensService struct {
	repo    interface{}
	manager auth.TokenManager
	cache   cache.Cache
}

func (t tokensService) Get(ctx context.Context, user *domain.User) (*domain.Token, error) {
	panic("implement me")
}

func (t tokensService) Refresh(ctx context.Context, refresh *domain.Token) (*domain.Token, error) {
	panic("implement me")
}

func (t tokensService) Validate(ctx context.Context, token *domain.Token) (ValidateStatus, error) {
	panic("implement me")
}

func NewAuthService(repo interface{}, cache cache.Cache) TokensService {
	return &tokensService{
		repo:  repo,
		cache: cache,
	}
}
