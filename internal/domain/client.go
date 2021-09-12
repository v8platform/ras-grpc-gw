package domain

import "github.com/lithammer/shortuuid/v3"

type Client struct {
	UserID        string
	UUID          string
	Host          string
	AgentUser     string
	AgentPassword string
	Endpoints     []string
}

func m() {
	shortuuid.New()
}

type Endpoint struct {
	UUID string
}
