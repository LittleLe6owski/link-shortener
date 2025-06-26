package instance

import (
	"context"
	"fmt"

	"github.com/LittleLe6owski/link-shortener/config"
	"github.com/LittleLe6owski/link-shortener/internal/repository/postgresql"
	"github.com/LittleLe6owski/link-shortener/internal/repository/redis"
	redisLib "github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type infrastructure struct {
	PGConnect   postgresql.Connect
	RedisClient redisLib.UniversalClient
}

type infrastractureOption func(context.Context, config.Config, zerolog.Logger, *infrastructure) error

func initInfrastructure(
	ctx context.Context, cfg config.Config, log zerolog.Logger, opts ...infrastractureOption,
) (infrastructure, error) {
	infrastructure := new(infrastructure)

	for _, addOption := range opts {
		if err := addOption(ctx, cfg, log, infrastructure); err != nil {
			return *infrastructure, err
		}
	}

	return *infrastructure, nil
}

func WithPostgres(ctx context.Context, cfg config.Config, log zerolog.Logger, i *infrastructure) error {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s pool_max_conns=%s",
		cfg.Postgresql.User,
		cfg.Postgresql.Password,
		cfg.Postgresql.Host,
		cfg.Postgresql.Port,
		cfg.Postgresql.DBName,
		cfg.Postgresql.SslMode,
		cfg.Postgresql.PoolMaxConnections,
	)

	pgCfg := postgresql.Config{
		DSN:               dsn,
		MaxConnIdleTime:   cfg.Postgresql.MaxConnIdleTime,
		HealthCheckPeriod: cfg.Postgresql.HealthCheckPeriod,
		RequestTimeout:    cfg.Postgresql.RequestTimeout,
	}

	pgc, err := postgresql.NewConnector(ctx, pgCfg)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrWrongPostgresConn, err)
	}

	i.PGConnect = pgc

	return nil
}

func WithRedis(ctx context.Context, cfg config.Config, log zerolog.Logger, i *infrastructure) error {
	redisConfig := redis.Config{
		Hosts:        cfg.Redis.Host,
		Password:     cfg.Redis.Password,
		Database:     cfg.Redis.Database,
		ReadTimeout:  cfg.Redis.ReadTimeout,
		WriteTimeout: cfg.Redis.WriteTimeout,
	}

	redisClient, err := redis.NewClient(ctx, redisConfig)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrWrongRedisConn, err)
	}

	i.RedisClient = redisClient

	return nil
}
