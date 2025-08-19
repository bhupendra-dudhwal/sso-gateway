package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/bhupendra-dudhwal/go-hexagonal/internal/constants"
	"github.com/bhupendra-dudhwal/go-hexagonal/internal/core/models"
	"github.com/bhupendra-dudhwal/go-hexagonal/internal/utils"
	"github.com/redis/go-redis/v9"
)

type cache struct {
	config *models.Cache
	client *redis.Client
}

func NewCacheRepository(config *models.Cache, client *redis.Client) *cache {
	return &cache{
		config: config,
		client: client,
	}
}

// Get fetches a value from cache and optionally unmarshals into response.
func (c *cache) Get(ctx context.Context, key string, response any) (string, error) {

	// ctx, cancel := context.WithTimeout(ctx, c.config.ReadTimeout)
	// defer cancel()
	result, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", utils.ErrInvalidCacheKey
		}
		return "", fmt.Errorf("failed to get cache key %q: %w", key, err)
	}

	if response != nil {
		if err := json.Unmarshal([]byte(result), response); err != nil {
			return result, fmt.Errorf("failed to unmarshal cache key %q: %w", key, err)
		}
	}

	return result, nil
}

// Add stores or updates a value based on the provided strategy.
//   - CacheAdd: Only set if the key does not exist.
//   - CacheUpdate: Always set (overwrite if exists).
func (c *cache) Add(ctx context.Context, key string, value any, ttl time.Duration, strategy constants.CacheStrategy) error {
	// ctx, cancel := context.WithTimeout(ctx, c.config.WriteTimeout)
	// defer cancel()

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache value for key %q: %w", key, err)
	}

	switch strategy {
	case constants.CacheAdd:
		err = c.client.SetNX(ctx, key, data, ttl).Err()
	case constants.CacheUpdate:
		err = c.client.Set(ctx, key, data, ttl).Err()
	default:
		return fmt.Errorf("invalid cache strategy: %v", strategy)
	}

	if err != nil {
		return fmt.Errorf("failed to set cache key %q: %w", key, err)
	}

	return nil
}
