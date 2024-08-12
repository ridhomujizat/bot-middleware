package redis

import (
	"bot-middleware/internal/pkg/util"
	"context"
	"fmt"
	"time"

	redisClient "github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisClient struct {
	client *redisClient.Client
}

func NewRedisClient(addr, password string, db int) *RedisClient {
	rdb := redisClient.NewClient(&redisClient.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		util.HandleAppError(err, "Redis", "Ping", false)
		return nil
	}

	return &RedisClient{
		client: rdb,
	}
}

func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	err := r.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return util.HandleAppError(fmt.Errorf("failed set key %s: %w", key, err), "Redis", "Set", true)
	}
	return nil
}

func (r *RedisClient) Get(key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redisClient.Nil {
			return "", redisClient.Nil
		}
		return "", util.HandleAppError(fmt.Errorf("failed get key %s: %w", key, err), "Redis", "Get", true)
	}
	return val, nil
}

func (r *RedisClient) Del(key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return util.HandleAppError(fmt.Errorf("failed del key %s: %w", key, err), "Redis", "Del", true)
	}
	return nil
}

func (r *RedisClient) Exists(key string) (bool, error) {
	val, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, util.HandleAppError(fmt.Errorf("failed get exists key %s: %w", key, err), "Redis", "Exists", true)
	}
	return val > 0, nil
}

func (r *RedisClient) Expire(key string, expiration time.Duration) error {
	err := r.client.Expire(ctx, key, expiration).Err()
	if err != nil {
		return util.HandleAppError(fmt.Errorf("failed set expiration key %s: %w", key, err), "Redis", "Expire", true)
	}
	return nil
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}
