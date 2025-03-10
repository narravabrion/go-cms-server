package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/narravabrion/go-cms-server/internal/models"
)


type Storage struct  {
	Users interface {
		Get(context.Context, int64) (*models.User, error)
		Set(context.Context, *models.User) error
	}
}

func NewRedisStorage(redisDB *redis.Client) Storage {
	return Storage{
		Users: &UserStore{redisDB: redisDB},
	}
}