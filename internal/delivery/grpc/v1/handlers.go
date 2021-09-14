package v1

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/v8platform/ras-grpc-gw/internal/service"
	apiv1 "github.com/v8platform/ras-grpc-gw/pkg/gen/service/api/v1"
	"google.golang.org/grpc"
)

func RegisterServerServices(services *service.Services) (func(server *grpc.Server), func(ctx context.Context, mux *runtime.ServeMux) error) {
	users := NewUsersServiceServer(services)

	clientsStorage := NewRasClientsStorage()

	// auth := NewAuthServiceServer(services, clientsStorage)
	clusters := NewClustersServiceServer(services, clientsStorage)
	return func(server *grpc.Server) {
			apiv1.RegisterUsersServiceServer(server, users)
			apiv1.RegisterClustersServiceServer(server, clusters)
		},
		func(ctx context.Context, mux *runtime.ServeMux) error {
			if err := apiv1.RegisterUsersServiceHandlerServer(ctx, mux, users); err != nil {
				return err
			}

			if err := apiv1.RegisterClustersServiceHandlerServer(ctx, mux, clusters); err != nil {
				return err
			}

			// if err := apiv1.RegisterAuthServiceHandlerServer(ctx, mux, auth); err != nil {
			// 	return err
			// }

			return nil
		}
}
