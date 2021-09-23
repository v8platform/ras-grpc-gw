package client

import (
	protocolv1 "github.com/v8platform/protos/gen/ras/protocol/v1"
)

var _ protocolv1.Endpoint = (*Endpoint)(nil)

type Endpoint struct {
	UUID    string
	ID      int32
	Version int32
}

func (e *Endpoint) GetVersion() int32 {
	return e.Version
}

func (e *Endpoint) GetId() int32 {
	return e.ID
}
