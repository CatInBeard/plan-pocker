package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"shared/logger"

	"github.com/go-redis/redis/v8"
)

type CacheClient struct {
	client *redis.Client
	ctx    context.Context
}

var (
	redisClient CacheClient
	once        sync.Once
)

func GetCacheClient() CacheClient {
	once.Do(func() {
		redisHost := os.Getenv("REDIS_HOST")
		redisPort := os.Getenv("REDIS_PORT")
		redisPassword := os.Getenv("REDIS_PASSWORD")

		redisDB := 0
		if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
			if db, err := strconv.Atoi(dbStr); err == nil {
				redisDB = db
			} else {
				logger.Log(logger.WARNING, "Invalid db from env, used default", err.Error())
			}
		}

		rdb := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
			Password: redisPassword,
			DB:       redisDB,
		})

		redisClient = CacheClient{
			client: rdb,
			ctx:    context.Background(),
		}
	})

	return redisClient
}

func (c *CacheClient) SetValue(key string, value string, expiration time.Duration) error {
	err := c.client.Set(c.ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *CacheClient) GetValue(key string) (string, error) {
	val, err := c.client.Get(c.ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (c *CacheClient) GetKeysByPattern(pattern string) ([]string, error) {
	keys, err := c.client.Keys(c.ctx, pattern).Result()
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func (c CacheClient) SetStructValue(key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.SetValue(key, string(jsonData), expiration)
}

func (c *CacheClient) GetStructValue(key string, result interface{}) error {
	val, err := c.GetValue(key)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), result)
}

func (c *CacheClient) DeleteKey(key string) error {
	err := c.client.Del(c.ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *CacheClient) UpdateExpiration(key string, expiration time.Duration) error {
	_, err := c.client.Expire(c.ctx, key, expiration).Result()
	return err
}
