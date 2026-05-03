package redis

import (
	"context"
	"fmt"
	"time"
	"github.com/redis/go-redis/v9"
	"github.com/ZephyrJung/QiaoYiCommunity/config"
)

type Client struct {
	client *redis.Client
}

func NewClient(cfg config.RedisConfig) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}
	return &Client{client: client}, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *Client) Del(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

func (c *Client) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	return c.client.SetNX(ctx, key, value, ttl).Result()
}

func (c *Client) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.client.TTL(ctx, key).Result()
}

func SMSCodeKey(phone string) string {
	return fmt.Sprintf("sms:code:%s", phone)
}

func RefreshTokenKey(userID int64, tokenID string) string {
	return fmt.Sprintf("refresh:%d:%s", userID, tokenID)
}

func RateLimitKey(ip, path string) string {
	return fmt.Sprintf("ratelimit:%s:%s", ip, path)
}
