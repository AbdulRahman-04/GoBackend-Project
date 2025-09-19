package utils

import (
	"context"
	"fmt"

	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/config"
	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
)

func ConnectRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.Redis.Host,
		Password: config.AppConfig.Redis.Password,
		DB:       config.AppConfig.Redis.DB,
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("❌ Redis connection failed: %v", err))
	}

	fmt.Println("✅ Connected to Redis successfully")
}

func RedisSet(key string, value string) error {
	return RedisClient.Set(Ctx, key, value, 0).Err()
}

func RedisGet(key string) (string, error) {
	return RedisClient.Get(Ctx, key).Result()
}

func RedisDel(key string) error {
	return RedisClient.Del(Ctx, key).Err()
}
