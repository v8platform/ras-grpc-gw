package client

import "errors"

var (
	ErrClosed         = errors.New("ras-client: client is closed")
	ErrUnknownMessage = errors.New("ras-client: unknown message packet")
	ErrTimeout        = errors.New("ras-client: client timeout")
)
