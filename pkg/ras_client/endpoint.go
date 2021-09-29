package client

import (
	"github.com/google/uuid"
	protocolv1 "github.com/v8platform/protos/gen/ras/protocol/v1"
	"time"
)

var _ protocolv1.Endpoint = (*ChannelEndpoint)(nil)

type ChannelEndpoint struct {
	UUID    uuid.UUID
	ID      int32
	Version int32
	UsedAt  time.Time
}

func (c *ChannelEndpoint) GetUsedAt() time.Time {
	return c.UsedAt
}

func (c *ChannelEndpoint) SetUsedAt(tm time.Time) {
	c.UsedAt = tm
}

func (e ChannelEndpoint) GetVersion() int32 {
	return e.Version
}

func (e ChannelEndpoint) GetId() int32 {
	return e.ID
}

type Endpoint struct {
	UUID uuid.UUID

	Ver     int32
	version string

	ChangeAt time.Time
}

func (e *Endpoint) Version() string {
	return e.version
}
