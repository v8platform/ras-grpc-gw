// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package apiv1

import (
	context "context"
	v1 "github.com/v8platform/protos/gen/ras/messages/v1"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AuthServiceClient is the client API for AuthService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthServiceClient interface {
	AuthenticateCluster(ctx context.Context, in *v1.ClusterAuthenticateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	AuthenticateInfobase(ctx context.Context, in *v1.AuthenticateInfobaseRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	AuthenticateAgent(ctx context.Context, in *v1.AuthenticateAgentRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type authServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthServiceClient(cc grpc.ClientConnInterface) AuthServiceClient {
	return &authServiceClient{cc}
}

func (c *authServiceClient) AuthenticateCluster(ctx context.Context, in *v1.ClusterAuthenticateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/service.api.v1.AuthService/AuthenticateCluster", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) AuthenticateInfobase(ctx context.Context, in *v1.AuthenticateInfobaseRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/service.api.v1.AuthService/AuthenticateInfobase", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) AuthenticateAgent(ctx context.Context, in *v1.AuthenticateAgentRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/service.api.v1.AuthService/AuthenticateAgent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthServiceServer is the server API for AuthService service.
// All implementations must embed UnimplementedAuthServiceServer
// for forward compatibility
type AuthServiceServer interface {
	AuthenticateCluster(context.Context, *v1.ClusterAuthenticateRequest) (*emptypb.Empty, error)
	AuthenticateInfobase(context.Context, *v1.AuthenticateInfobaseRequest) (*emptypb.Empty, error)
	AuthenticateAgent(context.Context, *v1.AuthenticateAgentRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedAuthServiceServer()
}

// UnimplementedAuthServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAuthServiceServer struct {
}

func (UnimplementedAuthServiceServer) AuthenticateCluster(context.Context, *v1.ClusterAuthenticateRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AuthenticateCluster not implemented")
}
func (UnimplementedAuthServiceServer) AuthenticateInfobase(context.Context, *v1.AuthenticateInfobaseRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AuthenticateInfobase not implemented")
}
func (UnimplementedAuthServiceServer) AuthenticateAgent(context.Context, *v1.AuthenticateAgentRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AuthenticateAgent not implemented")
}
func (UnimplementedAuthServiceServer) mustEmbedUnimplementedAuthServiceServer() {}

// UnsafeAuthServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthServiceServer will
// result in compilation errors.
type UnsafeAuthServiceServer interface {
	mustEmbedUnimplementedAuthServiceServer()
}

func RegisterAuthServiceServer(s grpc.ServiceRegistrar, srv AuthServiceServer) {
	s.RegisterService(&AuthService_ServiceDesc, srv)
}

func _AuthService_AuthenticateCluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.ClusterAuthenticateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).AuthenticateCluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.api.v1.AuthService/AuthenticateCluster",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).AuthenticateCluster(ctx, req.(*v1.ClusterAuthenticateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_AuthenticateInfobase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.AuthenticateInfobaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).AuthenticateInfobase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.api.v1.AuthService/AuthenticateInfobase",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).AuthenticateInfobase(ctx, req.(*v1.AuthenticateInfobaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_AuthenticateAgent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.AuthenticateAgentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).AuthenticateAgent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.api.v1.AuthService/AuthenticateAgent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).AuthenticateAgent(ctx, req.(*v1.AuthenticateAgentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AuthService_ServiceDesc is the grpc.ServiceDesc for AuthService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuthService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "service.api.v1.AuthService",
	HandlerType: (*AuthServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AuthenticateCluster",
			Handler:    _AuthService_AuthenticateCluster_Handler,
		},
		{
			MethodName: "AuthenticateInfobase",
			Handler:    _AuthService_AuthenticateInfobase_Handler,
		},
		{
			MethodName: "AuthenticateAgent",
			Handler:    _AuthService_AuthenticateAgent_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service/api/v1/auth_service.proto",
}
