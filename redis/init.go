package redis

import (
	"context"
	"fmt"
	redis "github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"time"
)

var rClient *redis.Client

func init() {
	rClient = redis.NewClient(&redis.Options{
		Addr:               "172.25.240.10:30379",
		Password:           "Welcome@123",
		DB:                 0,
		PoolSize:           30,
		MinIdleConns:       10,
		PoolTimeout:        3 * time.Second,
		IdleCheckFrequency: 30 * time.Second,
		IdleTimeout:        60 * time.Second,
	})
}

func GetToken(key string) (string, bool) {
	val, err := rClient.Get(context.Background(), key).Result()
	if err != nil {
		zap.L().Error("redis get error.", zap.Error(err))
		return "", false
	}
	zap.L().Debug(val)

	return val, true
}

func SetToken(token string) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "172.25.240.10:30379",
		Password: "Welcome@123",
		DB:       0,
	})

	err := rdb.Set(context.Background(), "token1", token, 0).Err()
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("set token successful")

	return nil
}
