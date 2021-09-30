package client

import (
	"github.com/google/uuid"
)

type Endpoint struct {
	UUID uuid.UUID

	Ver     int32
	version string
}

func (e *Endpoint) Version() string {
	return e.version
}
