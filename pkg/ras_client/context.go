package client

import "context"

type endpointValue struct{}

func EndpointFromContext(ctx context.Context) string {

	return ctx.Value(endpointValue{}).(string)

}

func EndpointToContext(ctx context.Context, uuid string) context.Context {

	return context.WithValue(ctx, endpointValue{}, uuid)

}

type channelValue struct{}

func ChannelFromContext(ctx context.Context) *Channel {

	return ctx.Value(channelValue{}).(*Channel)

}

func ChannelToContext(ctx context.Context, c *Channel) context.Context {

	return context.WithValue(ctx, channelValue{}, c)

}
