package db

import (
	"context"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/snlansky/coral/pkg/logging"
)

var logger = logging.MustGetLogger("cache")

var _client *redis.Client

func InitRedis(uri string, db int, password string) error {
	if uri == "" {
		logger.Warn("redis not setup, we will not have cache")
		return nil
	}
	_client = redis.NewClient(&redis.Options{
		Network:      "",
		Addr:         uri,
		DB:           db,
		OnConnect:    nil,
		Username:     "",
		Password:     password,
		DialTimeout:  time.Second * 5,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		PoolSize:     16,
		MinIdleConns: 4,
	})

	_, err := _client.Ping(context.Background()).Result()
	if err != nil {
		return err
	}

	return nil
}

func GetClient() *redis.Client {
	return _client
}
