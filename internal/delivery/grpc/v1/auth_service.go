package v1

import (
	"context"
	"github.com/v8platform/ras-grpc-gw/internal/domain"
	service2 "github.com/v8platform/ras-grpc-gw/internal/service"
	"github.com/v8platform/ras-grpc-gw/pkg/gen/access/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServerService interface {
	service.AuthServiceServer
}

type authServerService struct {
	service.UnimplementedAuthServiceServer
	services *service2.Services
}

func (a authServerService) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {

	return ctx, nil

}

func (a authServerService) SingIn(ctx context.Context, request *service.GetRequest) (*service.Tokens, error) {

	user, err := a.services.Users.GetByCredentials(ctx, request.GetUser(), request.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "invalid user or password")
	}

	tokens, err := a.services.Tokens.Get(ctx, user.UUID)
	if err != nil {
		return nil, err
	}

	return &service.Tokens{
		AccessToken:  string(tokens.Access),
		RefreshToken: string(tokens.Refresh),
	}, err
}

func (a authServerService) Refresh(ctx context.Context, request *service.RefreshRequest) (*service.Tokens, error) {

	tokens, err := a.services.Tokens.Refresh(ctx, domain.RefreshToken(request.RefreshToken))
	if err != nil {
		return nil, err
	}

	return &service.Tokens{
		AccessToken:  string(tokens.Access),
		RefreshToken: string(tokens.Refresh),
	}, err
}

func NewAuthServerService(services *service2.Services) AuthServerService {
	return &authServerService{
		services: services,
	}
}