package redislink

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/LittleLe6owski/link-shortener/internal/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Link struct {
	ttl          time.Duration
	client       redis.UniversalClient
	globalPrefix string
}

func NewRepository(client redis.UniversalClient, defaultTTL time.Duration, prefix string) Link {
	return Link{client: client, ttl: defaultTTL, globalPrefix: prefix}
}

func (l Link) Set(ctx context.Context, link domain.Link) error {
	b, err := json.Marshal(link)
	if err != nil {
		return err
	}

	return l.client.Set(ctx, l.formKey(link.ID), b, l.ttl).Err()
}

func (l Link) GetByID(ctx context.Context, id uuid.UUID) (domain.Link, error) {
	val, err := l.client.Get(ctx, l.formKey(id)).Bytes()
	if err != nil {
		return domain.Link{}, err
	}

	link := domain.Link{}
	err = json.Unmarshal(val, &link)
	if err != nil {
		return domain.Link{}, err
	}

	return link, nil
}

func (l Link) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return l.client.Del(ctx, l.formKey(id)).Err()
}

func (r Link) Clear(ctx context.Context) error {
	switch client := r.client.(type) {
	case *redis.Client:
		iter := client.Scan(ctx, 0, fmt.Sprintf("%s:*", r.globalPrefix), 0).Iterator()
		for iter.Next(ctx) {
			if err := client.Del(ctx, iter.Val()).Err(); err != nil {
				return err
			}
		}
		return iter.Err()
	case *redis.ClusterClient:
		return client.ForEachMaster(ctx, func(ctx context.Context, client *redis.Client) error {
			iter := client.Scan(ctx, 0, fmt.Sprintf("%s:*", r.globalPrefix), 0).Iterator()
			for iter.Next(ctx) {
				if err := client.Del(ctx, iter.Val()).Err(); err != nil {
					return err
				}
			}
			return iter.Err()
		})
	default:
		return fmt.Errorf("unknown client type")
	}
}

func (l Link) formKey(id uuid.UUID) string {
	return strings.Join([]string{l.globalPrefix, id.String()}, "_")
}
