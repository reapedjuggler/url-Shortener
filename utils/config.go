package utils

import (
	"github.com/go-redis/redis"
)

var client *redis.Client = nil

func InitRedis() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func GetClient() *redis.Client {

	if client == nil {
		InitRedis()
	}

	return client
}
