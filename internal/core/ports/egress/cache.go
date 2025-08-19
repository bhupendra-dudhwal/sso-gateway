package egress

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type CacheConnectionPorts interface {
	Connect(ctx context.Context) (*redis.Client, error)
	Close() error
}
