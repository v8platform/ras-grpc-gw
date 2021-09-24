package client

import (
	protocolv1 "github.com/v8platform/protos/gen/ras/protocol/v1"
)

var _ protocolv1.Endpoint = (*channelEndpoint)(nil)

type channelEndpoint struct {
	ID      int32
	Version int32
}

func (e channelEndpoint) GetVersion() int32 {
	return e.Version
}

func (e channelEndpoint) GetId() int32 {
	return e.ID
}

func newChannelEndpoint(id, version int32) channelEndpoint {
	return channelEndpoint{
		ID:      id,
		Version: version,
	}
}

type Endpoint struct {
	UUID string

	Ver     int32
	version string
}

func (e *Endpoint) Version() string {
	return e.version
}
