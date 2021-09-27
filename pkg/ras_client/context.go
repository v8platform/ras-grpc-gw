package client

import "context"

type endpointValue struct{}

func EndpointFromContext(ctx context.Context) *ChannelEndpoint {
	val := ctx.Value(endpointValue{})
	if val == nil {
		return nil
	}
	return val.(*ChannelEndpoint)

}

func EndpointToContext(ctx context.Context, endpoint *ChannelEndpoint) context.Context {

	return context.WithValue(ctx, endpointValue{}, endpoint)

}

type channelValue struct{}

func ChannelFromContext(ctx context.Context) *Channel {
	val := ctx.Value(channelValue{})
	if val == nil {
		return nil
	}
	return val.(*Channel)

}

func ChannelToContext(ctx context.Context, c *Channel) context.Context {

	return context.WithValue(ctx, channelValue{}, c)

}
