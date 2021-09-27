package context

import (
	"context"
)

type identEndpoint struct{}

func EndpointToContext(ctx context.Context, value string) context.Context {

	return context.WithValue(ctx, identEndpoint{}, value)

}

func EndpointFromContext(ctx context.Context) (string, bool) {

	u, ok := ctx.Value(identEndpoint{}).(string)
	return u, ok

}

func RequestOptsToContext(ctx context.Context, value []interface{}) context.Context {

	return context.WithValue(ctx, identEndpoint{}, value)

}

func RequestOptsFromContext(ctx context.Context) ([]interface{}, bool) {

	u, ok := ctx.Value(identEndpoint{}).([]interface{})
	return u, ok

}
