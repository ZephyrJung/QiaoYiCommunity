package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ZephyrJung/QiaoYiCommunity/pkg/redis"
	"github.com/ZephyrJung/QiaoYiCommunity/pkg/response"
	"github.com/gin-gonic/gin"
	goRedis "github.com/redis/go-redis/v9"
)

var rateLimitScript = goRedis.NewScript(`
	local key = KEYS[1]
	local max = tonumber(ARGV[1])
	local window = tonumber(ARGV[2])
	local current = redis.call('INCR', key)
	if current == 1 then
		redis.call('EXPIRE', key, window)
	end
	return current
`)

type RateLimiter struct {
	redis       *redis.Client
	maxRequests int
	window      time.Duration
}

func NewRateLimiter(redisClient *redis.Client, maxRequests int, windowSeconds int) *RateLimiter {
	return &RateLimiter{
		redis:       redisClient,
		maxRequests: maxRequests,
		window:      time.Duration(windowSeconds) * time.Second,
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		key := fmt.Sprintf("ratelimit:%s:%s", ip, path)

		result, err := rl.redis.RunScript(context.Background(), rateLimitScript, []string{key}, rl.maxRequests, int(rl.window.Seconds()))
		if err != nil {
			c.Next()
			return
		}

		count, ok := result.(int64)
		if !ok {
			c.Next()
			return
		}

		if count > int64(rl.maxRequests) {
			c.Header("Retry-After", fmt.Sprintf("%d", int(rl.window.Seconds())))
			response.ErrorWithStatus(c, http.StatusTooManyRequests, response.CodeTooMany, "too many requests")
			c.Abort()
			return
		}

		c.Next()
	}
}
