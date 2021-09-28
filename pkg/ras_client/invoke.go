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
	)

	requestOptions := combine(c.endpointOptions, opts)

	for _, option := range requestOptions {
		switch option.Ident() {
		case endpointIdent{}:
			getEndpoint := option.Value().(func(context.Context) (uuid.UUID, bool))
			endpointUUID, _ = getEndpoint(ctx)
		case interceptorsIdent{}:
			interceptor = option.Value().(Interceptor)
		case requestTimeoutIdent{}:
			timeout = option.Value().(time.Duration)
		}
	}
	if timeout > 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
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
		idxEndpoints := channel.Endpoints()

		channelEndpoint = EndpointFromContext(ctx)
		if channelEndpoint != nil {
			if len(idxEndpoints) > 0 {
				id, ok := idxEndpoints[channelEndpoint.UUID]
				if !(ok || id == channelEndpoint.ID) {
					return nil, ErrNotChannelEndpoint
				}
			}
		} else {

			endpoint, ok := c.endpoints[endpointUUID]
			if !ok {
				if endpointUUID == uuid.Nil {
					endpointUUID = uuid.New()
				}

				endpoint = &Endpoint{
					UUID:    endpointUUID,
					Ver:     10,
					version: "10.0",
				}

				c.endpoints[endpointUUID] = endpoint
				c.endpointConfig[endpointUUID] = &EndpointConfig{}

			}

			id, err := c.initEndpoint(ctx, channel, endpoint)
			if err != nil {
				return nil, err
			}
			channelEndpoint = &ChannelEndpoint{
				UUID:    endpointUUID,
				ID:      id,
				Version: endpoint.Ver,
			}
			ctx = EndpointToContext(ctx, channelEndpoint)
		}
	}

	ctx, err = md.AnnotateContext(ctx, c.metadataAnnotators, req)

	if err != nil {
		return nil, err
	}

	if len(c.Interceptors) > 0 {
		interceptor = ChainInterceptor(ChainInterceptor(c.Interceptors...), interceptor)
	}

	return handler(ctx, channel, channelEndpoint, req, clientv1.Interceptor(interceptor))
}
