package v1

//
// import (
// 	"context"
//
// 	"github.com/v8platform/ras-grpc-gw/internal/domain"
// 	"github.com/v8platform/ras-grpc-gw/internal/service"
// 	apiv1 "github.com/v8platform/ras-grpc-gw/pkg/gen/service/api/v1"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// 	"google.golang.org/protobuf/types/known/emptypb"
// )
//
// type UsersServiceServer interface {
// 	apiv1.UsersServiceServer
// }
//
// type usersServiceServer struct {
// 	apiv1.UnimplementedUsersServiceServer
// 	services *service.Services
// }
//
// func (a usersServiceServer) Login(ctx context.Context, request *apiv1.LoginRequest) (*apiv1.Tokens, error) {
// 	user, err := a.services.Users.GetByCredentials(ctx, request.GetUser(), request.GetPassword())
// 	if err != nil {
// 		return nil, status.Errorf(codes.PermissionDenied, "invalid user or password")
// 	}
//
// 	tokens, err := a.services.Tokens.Get(ctx, user.UUID)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &apiv1.Tokens{
// 		AccessToken:  string(tokens.Access),
// 		RefreshToken: string(tokens.Refresh),
// 	}, err
// }
//
// func (a usersServiceServer) CreateUser(ctx context.Context, request *apiv1.CreateUserRequest) (*emptypb.Empty, error) {
// 	_, err := a.services.Users.Register(ctx, request.GetUser(), request.GetPassword())
// 	if err != nil {
// 		return nil, status.Errorf(codes.InvalidArgument, err.Error())
// 	}
// 	return &emptypb.Empty{}, err
// }
//
// func (a usersServiceServer) Refresh(ctx context.Context, request *apiv1.RefreshRequest) (*apiv1.Tokens, error) {
// 	tokens, err := a.services.Tokens.Refresh(ctx, domain.RefreshToken(request.RefreshToken))
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &apiv1.Tokens{
// 		AccessToken:  string(tokens.Access),
// 		RefreshToken: string(tokens.Refresh),
// 	}, err
// }
//
// func (a usersServiceServer) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
//
// 	return ctx, nil
//
// }
//
// func NewUsersServiceServer(services *service.Services) UsersServiceServer {
// 	return &usersServiceServer{
// 		services: services,
// 	}
// }
