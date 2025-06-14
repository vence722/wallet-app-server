package redis

import (
	"context"
	"os"
	"wallet-app-server/app/config"
	"wallet-app-server/app/logger"

	goredis "github.com/redis/go-redis/v9"
)

// A global Redis client
var Client *RedisClient

// Init Redis, must be called before using the client
func Init() {
	redisConf := config.Cfg.Redis
	Client = &RedisClient{
		rdb: goredis.NewClient(&goredis.Options{
			Addr:     redisConf.Addr,
			Password: redisConf.Password,
			DB:       redisConf.DB,
		}),
	}
	if _, err := Client.rdb.Ping(context.Background()).Result(); err != nil {
		logger.Error("Redis init error: ", err.Error())
		os.Exit(-1)
	}
	logger.Info("Redis init sucess")
}

// A wrapper for the underlying Redis client library
type RedisClient struct {
	rdb *goredis.Client
}

// Get value by key
func (rc *RedisClient) Get(key string) (string, error) {
	res, err := rc.rdb.Get(context.Background(), key).Result()
	if err != nil && err != goredis.Nil {
		return res, err
	}
	return res, nil
}

// Set key and value
func (rc *RedisClient) Set(key string, value any) error {
	_, err := rc.rdb.SetArgs(context.Background(), key, value, goredis.SetArgs{KeepTTL: true}).Result()
	return err
}
