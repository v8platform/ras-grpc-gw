package service

import "github.com/v8platform/ras-grpc-gw/pkg/cache"

type UsersService struct {
	repo  interface{}
	cache cache.Cache
}

func NewUsersService(repo interface{}, cache cache.Cache) *UsersService {
	return &UsersService{
		repo:  repo,
		cache: cache,
	}
}
