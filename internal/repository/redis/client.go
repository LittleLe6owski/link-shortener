package redis

import (
	"context"
	"strings"

	"github.com/redis/go-redis/v9"
)

const (
	ErrNil = redis.Nil
)

func NewClient(ctx context.Context, cfg Config) (redis.UniversalClient, error) {
	addr := strings.Split(cfg.Hosts, ",")

	cl := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:        addr,
		Password:     cfg.Password,
		DB:           cfg.Database,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	if c, ok := cl.(*redis.ClusterClient); ok {
		err := c.ForEachShard(ctx, func(ctx context.Context, shard *redis.Client) error {
			return shard.Ping(ctx).Err()
		})
		if err != nil {
			return nil, err
		}
	}

	return cl, nil
}
