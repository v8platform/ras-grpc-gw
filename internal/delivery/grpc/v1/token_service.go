package v1

import (
	"context"
	"github.com/v8platform/ras-grpc-gw/internal/domain"
	service2 "github.com/v8platform/ras-grpc-gw/internal/service"
	"github.com/v8platform/ras-grpc-gw/pkg/gen/access/service"
)

type TokenServerService interface {
	service.TokenServiceServer
}

type tokenServerService struct {
	service.UnimplementedTokenServiceServer
	services *service2.Services
}

func (a tokenServerService) Get(ctx context.Context, request *service.GetRequest) (*service.GetTokenResponse, error) {

	user, err := a.services.Users.GetByCredentials(ctx, request.GetUser(), request.GetPassword())
	if err != nil {
		return nil, err
	}

	tokens, err := a.services.Tokens.Get(ctx, user)
	if err != nil {
		return nil, err
	}

	return &service.GetTokenResponse{
		AccessToken:  string(tokens.Access),
		RefreshToken: string(tokens.Refresh),
	}, err
}

func (a tokenServerService) Refresh(ctx context.Context, request *service.RefreshRequest) (*service.GetTokenResponse, error) {

	tokens, err := a.services.Tokens.Refresh(ctx, domain.RefreshToken(request.RefreshToken))
	if err != nil {
		return nil, err
	}

	return &service.GetTokenResponse{
		AccessToken:  string(tokens.Access),
		RefreshToken: string(tokens.Refresh),
	}, err
}

func NewTokenServerService(services *service2.Services) TokenServerService {
	return &tokenServerService{
		services: services,
	}
}

func userClaimFromToken(struct{}) string {
	return "foobar"
}

// exampleAuthFunc is used by a middleware to authenticate requests
func exampleAuthFunc(ctx context.Context) (context.Context, error) {
	// token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	// if err != nil {
	// 	return nil, err
	// }
	//
	// tokenInfo, err := parseToken(token)
	// if err != nil {
	// 	return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	// }
	//
	// // WARNING: in production define your own type to avoid context collisions
	// newCtx := context.WithValue(ctx, "tokenInfo", tokenInfo)
	//
	// return newCtx, nil

	return nil, nil
}
