package context

import (
	"context"
	"github.com/v8platform/ras-grpc-gw/internal/domain"
)

type identUser struct{}
type identClient struct{}
type identEndpoint struct{}

func UserToContext(ctx context.Context, user domain.User) context.Context {

	return context.WithValue(ctx, identUser{}, user)

}

func UserFromContext(ctx context.Context) (domain.User, bool) {

	u, ok := ctx.Value(identUser{}).(domain.User)
	return u, ok

}

func ClientToContext(ctx context.Context, value domain.Client) context.Context {

	return context.WithValue(ctx, identClient{}, value)

}

func ClientFromContext(ctx context.Context) (domain.Client, bool) {

	u, ok := ctx.Value(identClient{}).(domain.Client)
	return u, ok

}

func EndpointToContext(ctx context.Context, value string) context.Context {

	return context.WithValue(ctx, identEndpoint{}, value)

}

func EndpointFromContext(ctx context.Context) (string, bool) {

	u, ok := ctx.Value(identEndpoint{}).(string)
	return u, ok

}
