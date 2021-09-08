package server

import (
	"context"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/v8platform/ras-grpc-gw/pkg/gen/access/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AccessServer interface {
	service.TokenServiceServer
	service.ClientServiceServer
	ValidateToken(token string) (bool, error)
	ValidateHash(token string, hash string) (bool, error)
}

type accessServer struct {
	service.UnimplementedTokenServiceServer
	service.UnimplementedClientServiceServer
}

func (a accessServer) ValidateToken(token string) (bool, error) {
	panic("implement me")
}

func (a accessServer) ValidateHash(token string, hash string) (bool, error) {
	panic("implement me")
}

func NewAccessServer() AccessServer {
	return &accessServer{}
}

func parseToken(token string) (struct{}, error) {
	return struct{}{}, nil
}

func userClaimFromToken(struct{}) string {
	return "foobar"
}

// exampleAuthFunc is used by a middleware to authenticate requests
func exampleAuthFunc(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	tokenInfo, err := parseToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	// WARNING: in production define your own type to avoid context collisions
	newCtx := context.WithValue(ctx, "tokenInfo", tokenInfo)

	return newCtx, nil
}
