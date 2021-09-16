package v1

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/v8platform/ras-grpc-gw/internal/service"
	apiv1 "github.com/v8platform/ras-grpc-gw/pkg/gen/service/api/v1"
	"google.golang.org/grpc"
)

func RegisterServerServices(services *service.Services) (func(server *grpc.Server), func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error) {
	users := NewUsersServiceServer(services)

	clientsStorage := NewRasClientsStorage()

	// auth := NewAuthServiceServer(services, clientsStorage)
	clusters := NewClustersServiceServer(services, clientsStorage)
	return func(server *grpc.Server) {
			apiv1.RegisterUsersServiceServer(server, users)
			apiv1.RegisterApplicationsServiceServer(server, NewApplicationsServerService(services))
			apiv1.RegisterClustersServiceServer(server, clusters)
			apiv1.RegisterAuthServiceServer(server, NewAuthServiceServer(services, clientsStorage))
			apiv1.RegisterSessionsServiceServer(server, NewSessionsServiceServer(services, clientsStorage))
		},
		func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
			if err := apiv1.RegisterUsersServiceHandler(ctx, mux, conn); err != nil {
				return err
			}
			if err := apiv1.RegisterApplicationsServiceHandler(ctx, mux, conn); err != nil {
				return err
			}
			if err := apiv1.RegisterClustersServiceHandler(ctx, mux, conn); err != nil {
				return err
			}
			if err := apiv1.RegisterAuthServiceHandler(ctx, mux, conn); err != nil {
				return err
			}
			if err := apiv1.RegisterSessionsServiceHandler(ctx, mux, conn); err != nil {
				return err
			}

			// if err := apiv1.RegisterAuthServiceHandlerServer(ctx, mux, auth); err != nil {
			// 	return err
			// }

			return nil
		}
}
