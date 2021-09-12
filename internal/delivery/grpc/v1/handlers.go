package v1

import (
	ras_service "github.com/v8platform/protos/gen/ras/service/api/v1"
	"github.com/v8platform/ras-grpc-gw/internal/server"
	"github.com/v8platform/ras-grpc-gw/internal/service"
	access_service "github.com/v8platform/ras-grpc-gw/pkg/gen/access/service"
	"google.golang.org/grpc"
)

func NewHandlers(services *service.Services) []server.RegisterServerHandler {

	return []server.RegisterServerHandler{
		func(server *grpc.Server) {
			access_service.RegisterAuthServiceServer(server, NewAuthServerService(services))
		},
		func(server *grpc.Server) {
			access_service.RegisterClientServiceServer(server, NewClientServerService(services))
		},
	}

}

func NewRasHandlers(services *service.Services) []server.RegisterServerHandler {

	clients := NewRasClientsStorage()

	return []server.RegisterServerHandler{
		func(server *grpc.Server) {
			ras_service.RegisterAuthServiceServer(server, NewRasAuthServer(services, clients))
		},
		func(server *grpc.Server) {
			ras_service.RegisterClustersServiceServer(server, NewRasClustersServer(services, clients))
		},
	}

}
