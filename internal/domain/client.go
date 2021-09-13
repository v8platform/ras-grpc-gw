package domain

import "time"

type Client struct {
	UUID          string
	Name          string
	Host          string
	IdleTimeout   time.Duration
	AgentUser     string
	AgentPassword string
	Endpoints     []string
}

type Endpoint struct {
	UUID string
}
