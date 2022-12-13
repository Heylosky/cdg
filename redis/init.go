package redis

import (
	"context"
	"fmt"
	redis "github.com/go-redis/redis/v8"
)

func GetToken(key string) (string, bool) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "172.25.240.10:30379",
		Password: "Welcome@123",
		DB:       0,
	})

	val, err := rdb.Get(context.Background(), key).Result()
	if err != nil {
		fmt.Println(err)
		return "", false
	}
	fmt.Println(val)

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
