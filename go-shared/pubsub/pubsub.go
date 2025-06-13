package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/go-redis/redis/v8"
)

type PubSubClient struct {
	client *redis.Client
	ctx    context.Context
}

var (
	pubSubClient PubSubClient
	once         sync.Once
)

func GetPubSubClient() PubSubClient {
	once.Do(func() {
		redisHost := os.Getenv("REDIS_HOST")
		redisPort := os.Getenv("REDIS_PORT")
		redisPassword := os.Getenv("REDIS_PASSWORD")

		redisDB := 0
		if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
			if db, err := strconv.Atoi(dbStr); err == nil {
				redisDB = db
			}
		}

		rdb := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
			Password: redisPassword,
			DB:       redisDB,
		})

		pubSubClient = PubSubClient{
			client: rdb,
			ctx:    context.Background(),
		}
	})

	return pubSubClient
}

func (p *PubSubClient) Publish(channel string, message interface{}) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return p.client.Publish(p.ctx, channel, jsonData).Err()
}

func (p *PubSubClient) Subscribe(channel string) (<-chan *redis.Message, error) {
	pubsub := p.client.Subscribe(p.ctx, channel)
	_, err := pubsub.Receive(p.ctx)
	if err != nil {
		return nil, err
	}
	return pubsub.Channel(), nil
}

func (p *PubSubClient) Unsubscribe(channel string) error {
	pubsub := p.client.Subscribe(p.ctx, channel)
	return pubsub.Unsubscribe(p.ctx, channel)
}

func (p *PubSubClient) Close() error {
	return p.client.Close()
}
