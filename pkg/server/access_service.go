package server

import "github.com/v8platform/ras-grpc-gw/pkg/gen/access/service"

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
