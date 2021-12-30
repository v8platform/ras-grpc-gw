package client

import (
	"context"
	"github.com/google/uuid"
	clientv1 "github.com/v8platform/protos/gen/ras/client/v1"
	"github.com/v8platform/ras-grpc-gw/pkg/ras_client/md"
	"time"
)

func (c *client) Invoke(ctx context.Context, needEndpoint bool, req interface{}, handler clientv1.InvokeHandler, opts ...interface{}) (reply interface{}, err error) {

	var (
		channel         *Channel
		channelEndpoint *ChannelEndpoint
		interceptor     Interceptor
		endpointUUID    uuid.UUID
		timeout         time.Duration
		endpointConfig  EndpointConfig
		interceptors    []Interceptor
	)

	requestOptions := combine(c.endpointOptions, opts)

	for _, option := range requestOptions {
		switch option.Ident() {
		case endpointIdent{}:
			getEndpoint := option.Value().(func(context.Context) (uuid.UUID, bool))
			endpointUUID, _ = getEndpoint(ctx)
		case interceptorsIdent{}:
			interceptor = option.Value().(Interceptor)
			interceptors = append(interceptors, interceptor)
		case requestTimeoutIdent{}:
			timeout = option.Value().(time.Duration)
		}
	}
	if timeout > 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	ctx, err = md.AnnotateContext(ctx, c.metadataAnnotators, req)
	if err != nil {
		return nil, err
	}

	channel = ChannelFromContext(ctx)

	if channel == nil {
		channel, err = c.getChannel(ctx)
		if err != nil {
			return nil, err
		}
		defer c.putChannel(ctx, channel)
		ctx = ChannelToContext(ctx, channel)
	}

	if needEndpoint {

		channelEndpoint = EndpointFromContext(ctx)

		if channelEndpoint != nil && !channel.IsChannelEndpoint(channelEndpoint) {

			endpoint, ok := c.endpoints[endpointUUID]
			// TODO Переписать
			if !ok {
				if endpointUUID == uuid.Nil {
					endpointUUID = uuid.New()
				}

				endpoint = &Endpoint{
					UUID:    endpointUUID,
					Ver:     10,
					version: "10.0",
				}
				c.mu.Lock()
				c.endpoints[endpointUUID] = endpoint
				c.endpointConfig[endpointUUID] = c.defaultEndpointConfig.copy()
				c.mu.Unlock()
			}

			channelEndpoint, err = c.initEndpoint(ctx, channel, endpoint)

			if err != nil {
				return nil, err
			}

			ctx = EndpointToContext(ctx, channelEndpoint)
		} else {

			endpointUUID = uuid.New()

			endpoint := &Endpoint{
				UUID:    endpointUUID,
				Ver:     10,
				version: "10.0",
			}
			c.mu.Lock()
			c.endpoints[endpointUUID] = endpoint
			c.endpointConfig[endpointUUID] = c.defaultEndpointConfig.copy()
			c.mu.Unlock()
			channelEndpoint, err = c.initEndpoint(ctx, channel, endpoint)

			ctx = EndpointToContext(ctx, channelEndpoint)
		}
		interceptors = append(interceptors,
			AddDefaultInfobaseAuthInterceptor(endpointConfig),
			AddDefaultClusterAuthInterceptor(endpointConfig),
		)
	}

	return handler(ctx, channel, channelEndpoint, req, clientv1.Interceptor(ChainInterceptor(interceptors...)))
}
