package v1

//
// import (
// 	"context"
// 	"github.com/ungerik/go-dry"
// 	context2 "github.com/v8platform/ras-grpc-gw/internal/context"
// 	"github.com/v8platform/ras-grpc-gw/internal/domain"
// 	"github.com/v8platform/ras-grpc-gw/internal/service"
// 	"github.com/v8platform/ras-grpc-gw/pkg/gen/service/api/v1"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// )
//
// type ApplicationsServerService interface {
// 	apiv1.ApplicationsServiceServer
// }
//
// type applicationsServerService struct {
// 	apiv1.UnimplementedApplicationsServiceServer
// 	services *service.Services
// }
//
// func (c *applicationsServerService) GetApplications(ctx context.Context, request *apiv1.GetClientsRequest) (*apiv1.GetClientsResponse, error) {
// 	user, ok := context2.UserFromContext(ctx)
// 	if !ok {
// 		return nil, status.Error(codes.PermissionDenied, "unknown user")
// 	}
//
// 	var list []*apiv1.ApplicationInfo
//
// 	for _, appId := range user.Applications {
//
// 		app, err := c.services.Applications.GetByID(ctx, appId)
// 		if err != nil {
// 			continue
// 			// return nil, status.Error(codes.Internal, err.Error())
// 		}
//
// 		list = append(list, &apiv1.ApplicationInfo{
// 			Host:      app.Host,
// 			Uuid:      app.UUID,
// 			Connected: false,
// 			LastUsed:  nil,
// 			IdleAt:    nil,
// 			Name:      app.Name,
// 		})
// 	}
//
// 	return &apiv1.GetClientsResponse{
// 		Apps: list,
// 	}, nil
// }
//
// func (c *applicationsServerService) UpdateApplication(ctx context.Context, request *apiv1.UpdateApplicationRequest) (*apiv1.UpdateApplicationResponse, error) {
//
// 	user, ok := context2.UserFromContext(ctx)
// 	if !ok {
// 		return nil, status.Error(codes.PermissionDenied, "unknown user")
// 	}
//
// 	if !dry.StringInSlice(request.GetInfo().GetUuid(), user.Applications) {
// 		return nil, status.Error(codes.InvalidArgument, "application not found")
// 	}
//
// 	info := request.GetInfo()
//
// 	app, err := c.services.Applications.Update(ctx, domain.Application{
// 		UUID: info.GetUuid(),
// 		Name: info.GetName(),
// 		Host: info.GetHost(),
// 		// IdleTimeout:   info.GetIdleTimeout(),
// 		// AgentUser:     ,
// 		// AgentPassword: "",
// 		// Endpoints:     nil,
// 	})
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}
//
// 	return &apiv1.UpdateApplicationResponse{
// 		Info: &apiv1.ApplicationInfo{
// 			Host: app.Host,
// 			// IdleTimeout: app.IdleTimeout,
// 			Uuid: app.UUID,
// 			// Connected:   false,
// 			// LastUsed:    nil,
// 			// IdleAt:      nil,
// 			Name: app.Name,
// 		},
// 	}, nil
// }
//
// func (c *applicationsServerService) GetApplication(ctx context.Context, request *apiv1.GetApplicationRequest) (*apiv1.GetApplicationResponse, error) {
// 	panic("implement me")
// }
//
// func (c *applicationsServerService) DeleteApplication(ctx context.Context, request *apiv1.DeleteApplicationRequest) (*apiv1.DeleteApplicationResponse, error) {
// 	panic("implement me")
// }
//
// func (c *applicationsServerService) Register(ctx context.Context, request *apiv1.RegisterRequest) (*apiv1.RegisterResponse, error) {
//
// 	user, ok := context2.UserFromContext(ctx)
// 	if !ok {
// 		return nil, status.Error(codes.PermissionDenied, "unknown user")
// 	}
//
// 	application, err := c.services.Users.RegisterApplication(ctx, user.UUID, request.GetHost(), request.GetName())
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}
//
// 	if request.GetAgentUser() != "" {
// 		err := c.services.Applications.UpdateAuth(ctx, application.UUID, request.GetAgentUser(), request.GetAgentPwd())
// 		if err != nil {
// 			return nil, status.Error(codes.Internal, err.Error())
// 		}
// 	}
//
// 	return &apiv1.RegisterResponse{
// 		Uuid: application.UUID,
// 	}, nil
// }
//
// func NewApplicationsServerService(services *service.Services) ApplicationsServerService {
// 	return &applicationsServerService{
// 		services: services,
// 	}
// }
