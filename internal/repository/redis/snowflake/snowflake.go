package snowlake

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Snowflake struct {
	client       redis.UniversalClient
	globalPrefix string
}

func NewRepository(client redis.UniversalClient, prefix string) Snowflake {
	return Snowflake{client: client, globalPrefix: prefix}
}

func (sl Snowflake) CaptureNodeID(ctx context.Context) (int64, error) {
	return 0, nil
}
