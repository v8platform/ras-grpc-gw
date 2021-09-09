package service

import (
	"context"
	"github.com/v8platform/ras-grpc-gw/internal/domain"
	"github.com/v8platform/ras-grpc-gw/pkg/auth"
	"time"
)

// TokensService реализует бизнес-логику работы
type TokensService interface {
	Get(ctx context.Context, issuer string) (domain.Tokens, error)
	Refresh(ctx context.Context, refresh domain.RefreshToken) (domain.Tokens, error)
}

type tokensService struct {
	services   *Services
	tokens     auth.TokenManager
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func (t tokensService) Get(ctx context.Context, issuer string) (domain.Tokens, error) {

	return t.createTokens(issuer)

}

func (t tokensService) Refresh(ctx context.Context, refresh domain.RefreshToken) (domain.Tokens, error) {
	issuer, err := t.tokens.Validate(string(refresh), "refresh")
	if err != nil {
		return domain.Tokens{}, err
	}

	return t.createTokens(issuer)
}

func (t tokensService) createTokens(issuer string) (domain.Tokens, error) {

	access, err := t.tokens.Generate(issuer, "access", 10*time.Minute)
	if err != nil {
		return domain.Tokens{}, err
	}

	refresh, err := t.tokens.Generate(issuer, "refresh", 1*time.Hour)
	if err != nil {
		return domain.Tokens{}, err
	}

	return domain.Tokens{
		Access:  domain.AccessToken(access),
		Refresh: domain.RefreshToken(refresh),
	}, nil
}

func NewTokenService(tokenManager auth.TokenManager, manager *Services) TokensService {
	return &tokensService{
		tokens:   tokenManager,
		services: manager,
	}
}
