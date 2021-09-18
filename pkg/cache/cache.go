package cache

import (
	"github.com/v8platform/ras-grpc-gw/internal/config"
)

type Cache interface {
	Connect()
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	Clear(key string)
	HealthCheck() (bool, error)
}

type memoryCache struct {
	Cache
}

func NewMemoryCache() Cache {

	return &memoryCache{}

}

func New(config config.CacheConfig) (Cache, error) {

	switch config.Engine.Name() {

	case "memory":

		return NewMemoryCache(), nil

	case "redis":
		panic("TODO Add support redis")
	}

	return nil, nil
}
