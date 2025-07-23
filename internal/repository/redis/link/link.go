package redislink

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
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

	return l.client.Set(ctx, l.formKeyForSave(link.ID, link.ShortURI), b, l.ttl).Err()
}

func (l Link) GetByID(ctx context.Context, id uuid.UUID) (domain.Link, error) {
	key, err := l.formKeyForGet(id, 0)
	if err != nil {
		return domain.Link{}, err
	}

	iter := l.client.Scan(ctx, 0, key, 1).Iterator()
	if !iter.Next(ctx) {
		return domain.Link{}, nil
	}

	val, err := l.client.Get(ctx, iter.Val()).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return domain.Link{}, nil
		}

		return domain.Link{}, err
	}

	if err := l.client.Expire(ctx, iter.Val(), l.ttl).Err(); err != nil {
		return domain.Link{}, fmt.Errorf("failed to update TTL: %v", err)
	}

	link := domain.Link{}
	err = json.Unmarshal(val, &link)
	if err != nil {
		return domain.Link{}, err
	}

	return link, nil
}

func (l Link) GetByShortURI(ctx context.Context, shortURI int64) (domain.Link, error) {
	key, err := l.formKeyForGet(uuid.Nil, shortURI)
	if err != nil {
		return domain.Link{}, err
	}

	iter := l.client.Scan(ctx, 0, key, 1).Iterator()
	if !iter.Next(ctx) {
		return domain.Link{}, nil
	}

	val, err := l.client.Get(ctx, iter.Val()).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return domain.Link{}, nil
		}

		return domain.Link{}, err
	}

	if err := l.client.Expire(ctx, iter.Val(), l.ttl).Err(); err != nil {
		return domain.Link{}, fmt.Errorf("failed to update TTL: %v", err)
	}

	link := domain.Link{}
	err = json.Unmarshal(val, &link)
	if err != nil {
		return domain.Link{}, err
	}

	return link, nil
}

func (l Link) DeleteByID(ctx context.Context, id uuid.UUID) error {
	key, err := l.formKeyForGet(id, 0)
	if err != nil {
		return err
	}

	iter := l.client.Scan(ctx, 0, key, 1).Iterator()
	if !iter.Next(ctx) {
		return redis.Nil
	}

	return l.client.Del(ctx, iter.Val()).Err()
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

func (l Link) formKeyForSave(id uuid.UUID, shortURI int64) string {
	return strings.Join(
		[]string{l.globalPrefix, id.String(), strconv.Itoa(int(shortURI))}, "_",
	)
}

func (l Link) formKeyForGet(id uuid.UUID, shortURI int64) (string, error) {
	switch {
	case id == uuid.Nil:
		return strings.
			Join([]string{l.globalPrefix, "*", strconv.Itoa(int(shortURI))}, "_"), nil
	case shortURI == 0:
		return strings.
			Join([]string{l.globalPrefix, id.String(), "*"}, "_"), nil
	default:
		return "", errors.New("id and short uri is empty")
	}
}
