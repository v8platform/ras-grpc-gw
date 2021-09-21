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
