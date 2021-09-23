package client

import "context"

type endpointValue struct{}

func EndpointFromContext(ctx context.Context) string {

	return ctx.Value(endpointValue{}).(string)

}

func EndpointToContext(ctx context.Context, uuid string) context.Context {

	return context.WithValue(ctx, endpointValue{}, uuid)

}
