package config

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func NewRedisClient(cfg *viper.Viper) (redisClient *redis.Client, err error) {
	redisHost := cfg.GetString("REDIS_HOST")
	redisPort := cfg.GetString("REDIS_PORT")

	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: cfg.GetString("REDIS_PASSWORD"),
		DB:       0,
	})

	err = redisClient.Ping(context.Background()).Err()
	if err != nil {
		return
	}

	return
}
