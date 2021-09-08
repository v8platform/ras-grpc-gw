package server

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Authorization unary interceptor function to handle authorize per RPC call
func endpointInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	// Skip authorize when GetJWT is requested
	if info.FullMethod != "/proto.EventStoreService/GetJWT" {
		if err := authorize(ctx); err != nil {
			return nil, err
		}
	}

	// Calls the handler
	h, err := handler(ctx, req)

	return h, err
}

func withServerUnaryInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(endpointInterceptor)
}

// authorize function authorizes the token received from Metadata
func authorize(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
	}

	authHeader, ok := md["authorization"]
	if !ok {
		return status.Errorf(codes.Unauthenticated, "Authorization token is not supplied")
	}

	token := authHeader[0]
	// validateToken function validates the token
	err := validateToken(token)

	if err != nil {
		return status.Errorf(codes.Unauthenticated, err.Error())
	}
	return nil
}
