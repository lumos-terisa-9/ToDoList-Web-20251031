package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"team_task_hub/backend/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	Client *redis.Client
	ctx    = context.Background()
)

// InitRedis 初始化Redis连接
func InitRedis(cfg *config.Config) error {
	Client = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
		PoolSize: 20,
	})

	// 测试连接
	if err := Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis连接失败: %v", err)
	}

	log.Println("Redis连接成功")
	return nil
}

// 通用缓存操作
func setJson(key string, value any, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return Client.Set(ctx, key, data, expiration).Err()
}

func getJson(key string, dest any) error {
	data, err := Client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func setString(key, value string, expiration time.Duration) error {
	return Client.Set(ctx, key, value, expiration).Err()
}

func exists(key string) (bool, error) {
	result, err := Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

func deleteKey(key string) error {
	return Client.Del(ctx, key).Err()
}

// Close 关闭Redis连接
func Close() {
	if Client != nil {
		Client.Close()
	}
}
