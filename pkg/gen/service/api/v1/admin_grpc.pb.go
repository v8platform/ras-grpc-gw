// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package apiv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AccessClient is the client API for Access service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AccessClient interface {
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*Tokens, error)
	Refresh(ctx context.Context, in *RefreshRequest, opts ...grpc.CallOption) (*Tokens, error)
}

type accessClient struct {
	cc grpc.ClientConnInterface
}

func NewAccessClient(cc grpc.ClientConnInterface) AccessClient {
	return &accessClient{cc}
}

func (c *accessClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*Tokens, error) {
	out := new(Tokens)
	err := c.cc.Invoke(ctx, "/service.api.v1.Access/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accessClient) Refresh(ctx context.Context, in *RefreshRequest, opts ...grpc.CallOption) (*Tokens, error) {
	out := new(Tokens)
	err := c.cc.Invoke(ctx, "/service.api.v1.Access/Refresh", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AccessServer is the server API for Access service.
// All implementations must embed UnimplementedAccessServer
// for forward compatibility
type AccessServer interface {
	Login(context.Context, *LoginRequest) (*Tokens, error)
	Refresh(context.Context, *RefreshRequest) (*Tokens, error)
	mustEmbedUnimplementedAccessServer()
}

// UnimplementedAccessServer must be embedded to have forward compatible implementations.
type UnimplementedAccessServer struct {
}

func (UnimplementedAccessServer) Login(context.Context, *LoginRequest) (*Tokens, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedAccessServer) Refresh(context.Context, *RefreshRequest) (*Tokens, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Refresh not implemented")
}
func (UnimplementedAccessServer) mustEmbedUnimplementedAccessServer() {}

// UnsafeAccessServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AccessServer will
// result in compilation errors.
type UnsafeAccessServer interface {
	mustEmbedUnimplementedAccessServer()
}

func RegisterAccessServer(s grpc.ServiceRegistrar, srv AccessServer) {
	s.RegisterService(&Access_ServiceDesc, srv)
}

func _Access_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.api.v1.Access/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Access_Refresh_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessServer).Refresh(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.api.v1.Access/Refresh",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessServer).Refresh(ctx, req.(*RefreshRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Access_ServiceDesc is the grpc.ServiceDesc for Access service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Access_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "service.api.v1.Access",
	HandlerType: (*AccessServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _Access_Login_Handler,
		},
		{
			MethodName: "Refresh",
			Handler:    _Access_Refresh_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service/api/v1/admin.proto",
}

// ConfigClient is the client API for Config service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConfigClient interface {
	GetConfig(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ServerConfig, error)
	UpdateConfig(ctx context.Context, in *ServerConfig, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetClientStatus(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ClientStatus, error)
	GetConnections(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetConnectionsResponse, error)
	GetEndpoints(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetEndpointsResponse, error)
	ConfigureEndpoint(ctx context.Context, in *ConfigureEndpointRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type configClient struct {
	cc grpc.ClientConnInterface
}

func NewConfigClient(cc grpc.ClientConnInterface) ConfigClient {
	return &configClient{cc}
}

func (c *configClient) GetConfig(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ServerConfig, error) {
	out := new(ServerConfig)
	err := c.cc.Invoke(ctx, "/service.api.v1.Config/GetConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configClient) UpdateConfig(ctx context.Context, in *ServerConfig, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/service.api.v1.Config/UpdateConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configClient) GetClientStatus(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ClientStatus, error) {
	out := new(ClientStatus)
	err := c.cc.Invoke(ctx, "/service.api.v1.Config/GetClientStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configClient) GetConnections(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetConnectionsResponse, error) {
	out := new(GetConnectionsResponse)
	err := c.cc.Invoke(ctx, "/service.api.v1.Config/GetConnections", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configClient) GetEndpoints(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetEndpointsResponse, error) {
	out := new(GetEndpointsResponse)
	err := c.cc.Invoke(ctx, "/service.api.v1.Config/GetEndpoints", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configClient) ConfigureEndpoint(ctx context.Context, in *ConfigureEndpointRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/service.api.v1.Config/ConfigureEndpoint", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConfigServer is the server API for Config service.
// All implementations must embed UnimplementedConfigServer
// for forward compatibility
type ConfigServer interface {
	GetConfig(context.Context, *emptypb.Empty) (*ServerConfig, error)
	UpdateConfig(context.Context, *ServerConfig) (*emptypb.Empty, error)
	GetClientStatus(context.Context, *emptypb.Empty) (*ClientStatus, error)
	GetConnections(context.Context, *emptypb.Empty) (*GetConnectionsResponse, error)
	GetEndpoints(context.Context, *emptypb.Empty) (*GetEndpointsResponse, error)
	ConfigureEndpoint(context.Context, *ConfigureEndpointRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedConfigServer()
}

// UnimplementedConfigServer must be embedded to have forward compatible implementations.
type UnimplementedConfigServer struct {
}

func (UnimplementedConfigServer) GetConfig(context.Context, *emptypb.Empty) (*ServerConfig, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConfig not implemented")
}
func (UnimplementedConfigServer) UpdateConfig(context.Context, *ServerConfig) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateConfig not implemented")
}
func (UnimplementedConfigServer) GetClientStatus(context.Context, *emptypb.Empty) (*ClientStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetClientStatus not implemented")
}
func (UnimplementedConfigServer) GetConnections(context.Context, *emptypb.Empty) (*GetConnectionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConnections not implemented")
}
func (UnimplementedConfigServer) GetEndpoints(context.Context, *emptypb.Empty) (*GetEndpointsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEndpoints not implemented")
}
func (UnimplementedConfigServer) ConfigureEndpoint(context.Context, *ConfigureEndpointRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConfigureEndpoint not implemented")
}
func (UnimplementedConfigServer) mustEmbedUnimplementedConfigServer() {}

// UnsafeConfigServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConfigServer will
// result in compilation errors.
type UnsafeConfigServer interface {
	mustEmbedUnimplementedConfigServer()
}

func RegisterConfigServer(s grpc.ServiceRegistrar, srv ConfigServer) {
	s.RegisterService(&Config_ServiceDesc, srv)
}

func _Config_GetConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServer).GetConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.api.v1.Config/GetConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServer).GetConfig(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Config_UpdateConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServerConfig)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServer).UpdateConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.api.v1.Config/UpdateConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServer).UpdateConfig(ctx, req.(*ServerConfig))
	}
	return interceptor(ctx, in, info, handler)
}

func _Config_GetClientStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServer).GetClientStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.api.v1.Config/GetClientStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServer).GetClientStatus(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Config_GetConnections_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServer).GetConnections(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.api.v1.Config/GetConnections",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServer).GetConnections(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Config_GetEndpoints_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServer).GetEndpoints(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.api.v1.Config/GetEndpoints",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServer).GetEndpoints(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Config_ConfigureEndpoint_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigureEndpointRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServer).ConfigureEndpoint(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.api.v1.Config/ConfigureEndpoint",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServer).ConfigureEndpoint(ctx, req.(*ConfigureEndpointRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Config_ServiceDesc is the grpc.ServiceDesc for Config service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Config_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "service.api.v1.Config",
	HandlerType: (*ConfigServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetConfig",
			Handler:    _Config_GetConfig_Handler,
		},
		{
			MethodName: "UpdateConfig",
			Handler:    _Config_UpdateConfig_Handler,
		},
		{
			MethodName: "GetClientStatus",
			Handler:    _Config_GetClientStatus_Handler,
		},
		{
			MethodName: "GetConnections",
			Handler:    _Config_GetConnections_Handler,
		},
		{
			MethodName: "GetEndpoints",
			Handler:    _Config_GetEndpoints_Handler,
		},
		{
			MethodName: "ConfigureEndpoint",
			Handler:    _Config_ConfigureEndpoint_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service/api/v1/admin.proto",
}
