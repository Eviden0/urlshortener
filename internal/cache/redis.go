package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aeilang/urlshortener/config"
	"github.com/aeilang/urlshortener/internal/repo"
	"github.com/go-redis/redis/v8"
)

const (
	urlKeyPrefix = "url:"
)

type Cache interface {
	GetURL(ctx context.Context, shortCode string) (*repo.Url, error)
	SetURL(ctx context.Context, url repo.Url) error
	DeleteURL(ctx context.Context, shortCode string) error
	Close() error
}

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(cfg config.RedisConfig) (Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		DB:       cfg.DB,
		Password: cfg.Password,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &redisCache{client: client}, nil
}

func (c *redisCache) Close() error {
	return c.client.Close()
}

func (c *redisCache) GetURL(ctx context.Context, shortCode string) (*repo.Url, error) {
	key := urlKeyPrefix + shortCode
	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var url repo.Url
	if err := json.Unmarshal(data, &url); err != nil {
		return nil, err
	}

	// Check if URL has expired
	if url.ExpiresAt.Before(time.Now()) {
		c.client.Del(ctx, key)
		return nil, nil
	}

	return &url, nil
}

func (c *redisCache) SetURL(ctx context.Context, url repo.Url) error {
	key := urlKeyPrefix + url.ShortCode

	data, err := json.Marshal(url)
	if err != nil {
		return err
	}

	if url.ExpiresAt.Before(time.Now()) {
		return nil
	}

	return c.client.Set(ctx, key, data, time.Until(url.ExpiresAt)).Err()
}

func (c *redisCache) DeleteURL(ctx context.Context, shortCode string) error {
	key := urlKeyPrefix + shortCode
	return c.client.Del(ctx, key).Err()
}
