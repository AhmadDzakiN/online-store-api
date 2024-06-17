package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"time"
)

type cache struct {
	redisClient *redis.Client
}

type ICache interface {
	Write(ctx context.Context, key string, data interface{}, ttl time.Duration) (err error)
	Get(ctx context.Context, key string, data interface{}) (err error)
	Delete(ctx context.Context, key string) (err error)
	DeleteByKeyPattern(ctx context.Context, keyPattern string) (err error)
}

func NewCache(redisClient *redis.Client) ICache {
	return &cache{
		redisClient: redisClient,
	}
}

func (c *cache) Write(ctx context.Context, key string, data interface{}, ttl time.Duration) (err error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}
	err = c.redisClient.Set(ctx, key, jsonData, ttl).Err()
	return
}

func (c *cache) Get(ctx context.Context, key string, data interface{}) (err error) {
	cachedData, err := c.redisClient.Get(ctx, key).Result()
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(cachedData), data)
	if err != nil {
		return
	}

	return
}

func (c *cache) Delete(ctx context.Context, key string) (err error) {
	result, err := c.redisClient.Del(ctx, key).Result()
	if err != nil {
		return
	}

	if result == 0 {
		err = fmt.Errorf("key %s does not exist", key)
	} else {
		err = fmt.Errorf("cache deleted successfully")
	}

	return
}

func (c *cache) DeleteByKeyPattern(ctx context.Context, keyPattern string) (err error) {
	keys, err := c.redisClient.Keys(ctx, keyPattern).Result()
	if err != nil {
		return
	}

	if len(keys) == 0 {
		fmt.Println("No keys match the pattern")
		return
	}

	_, err = c.redisClient.Del(ctx, keys...).Result()
	if err != nil {
		return
	}

	log.Info().Msgf("Cache keys deleted successfully: %v", keys)

	return
}
