package client

import (
	"github.com/google/uuid"
	protocolv1 "github.com/v8platform/protos/gen/ras/protocol/v1"
)

var _ protocolv1.Endpoint = (*ChannelEndpoint)(nil)

type ChannelEndpoint struct {
	UUID    uuid.UUID
	ID      int32
	Version int32
}

func (e ChannelEndpoint) GetVersion() int32 {
	return e.Version
}

func (e ChannelEndpoint) GetId() int32 {
	return e.ID
}

func newChannelEndpoint(id, version int32) *ChannelEndpoint {
	return &ChannelEndpoint{
		ID:      id,
		Version: version,
	}
}

type Endpoint struct {
	UUID uuid.UUID

	Ver     int32
	version string
}

func (e *Endpoint) Version() string {
	return e.version
}
