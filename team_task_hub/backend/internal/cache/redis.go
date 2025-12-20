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

// RenewKey 为指定的缓存键续期
// key: 缓存键
// expiration: 新的过期时间
func RenewKey(key string, expiration time.Duration) error {
	// 使用 EXPIRE 命令重置键的生存时间
	err := Client.Expire(ctx, key, expiration).Err()
	if err != nil {
		log.Printf("缓存键续期失败: key=%s, error=%v", key, err)
		return err
	}
	return nil
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
