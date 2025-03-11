package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/narravabrion/go-cms-server/internal/models"
)

type UserStore struct {
	redisDB *redis.Client
}

const UserExp = time.Minute

func (rs *UserStore) Get(ctx context.Context, userID int64) (*models.User, error) {
	cacheKey := fmt.Sprintf("user-%v", userID)
	data, err := rs.redisDB.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	var user models.User
	if data != "" {
		err := json.Unmarshal([]byte(data), &user)
		if err != nil {
			return nil, err
		}
	}
	return &user, nil
}
func (rs *UserStore) Set(ctx context.Context, user *models.User) error {
	cacheKey := fmt.Sprintf("user-%v", user.ID)
	json, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return rs.redisDB.SetEX(ctx, cacheKey, json, UserExp).Err()
}
