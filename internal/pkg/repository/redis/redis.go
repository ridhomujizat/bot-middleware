package redis

import (
	"context"
	"fmt"
	"time"

	"bot-middleware/internal/pkg/util"

	"github.com/pterm/pterm"
	redisClient "github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisClient struct {
	client *redisClient.Client
}

func NewRedisClient(addr, password string, db int) *RedisClient {
	var rdb *redisClient.Client
	var err error

	for retries := 0; retries < 5; retries++ {
		rdb = redisClient.NewClient(&redisClient.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		})

		err = rdb.Ping(ctx).Err()
		if err == nil {
			break
		}

		pterm.Error.Printfln("retry %d: failed to connect to Redis: %+v", retries+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		util.HandleAppError(fmt.Errorf("failed to connect to Redis after retries: %w", err), "Redis", "Connect", true)
		return nil
	}

	return &RedisClient{
		client: rdb,
	}
}

func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	err := r.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return util.HandleAppError(fmt.Errorf("failed to set key %s: %w", key, err), "Redis", "Set", true)
	}
	return nil
}

func (r *RedisClient) Get(key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redisClient.Nil {
			util.HandleAppError(fmt.Errorf("key %s not found", key), "Redis", "Get", false)
			return "", redisClient.Nil
		}
		return "", util.HandleAppError(fmt.Errorf("failed to get key %s: %w", key, err), "Redis", "Get", true)
	}
	return val, nil
}

func (r *RedisClient) Del(key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return util.HandleAppError(fmt.Errorf("failed to delete key %s: %w", key, err), "Redis", "Del", true)
	}
	return nil
}

func (r *RedisClient) Exists(key string) (bool, error) {
	val, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, util.HandleAppError(fmt.Errorf("failed to check existence of key %s: %w", key, err), "Redis", "Exists", true)
	}
	return val > 0, nil
}

func (r *RedisClient) Expire(key string, expiration time.Duration) error {
	err := r.client.Expire(ctx, key, expiration).Err()
	if err != nil {
		return util.HandleAppError(fmt.Errorf("failed to set expiration for key %s: %w", key, err), "Redis", "Expire", true)
	}
	return nil
}

func (r *RedisClient) Close() error {
	err := r.client.Close()
	if err != nil {
		return util.HandleAppError(fmt.Errorf("failed to close Redis connection: %w", err), "Redis", "Close", true)
	}
	return nil
}
