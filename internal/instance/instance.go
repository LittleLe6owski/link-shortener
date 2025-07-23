package instance

import (
	"context"

	apiv1 "github.com/LittleLe6owski/link-shortener/api/v1"
	"github.com/LittleLe6owski/link-shortener/config"
	"github.com/LittleLe6owski/link-shortener/internal/controller/api"
	pglink "github.com/LittleLe6owski/link-shortener/internal/repository/postgresql/link"
	redislink "github.com/LittleLe6owski/link-shortener/internal/repository/redis/link"
	redissnowflake "github.com/LittleLe6owski/link-shortener/internal/repository/redis/snowflake"
	"github.com/LittleLe6owski/link-shortener/internal/usecase/link"
	"github.com/LittleLe6owski/link-shortener/pkg"
	"github.com/LittleLe6owski/link-shortener/pkg/server"
)

func New(cfg config.Config) (*server.App, error) {
	ctx := context.Background()

	linkShortener, err := server.NewApp(cfg)
	if err != nil {
		return nil, err
	}

	infra, err := initInfrastructure(ctx, linkShortener.Config(), linkShortener.Logger(), WithPostgres, WithRedis)
	if err != nil {
		return nil, err
	}

	pgLinkRepo := pglink.NewRepository(infra.PGConnect)
	redisLinkStorage := redislink.NewRepository(infra.RedisClient, cfg.Redis.DefaultTTL, cfg.Redis.PrefixKey)
	SnowflakeManager := redissnowflake.NewRepository(infra.RedisClient, cfg.Redis.PrefixKey)

	linkManager := link.NewLinkManger(
		pgLinkRepo,
		redisLinkStorage,
		SnowflakeManager,
		linkShortener.Logger(),
	)

	linkManagerController := api.NewLinkShortenerServiceServer(linkManager)

	grpcServer := pkg.GrpcProvider(linkShortener)
	grpcServer.RegisterService(
		&apiv1.LinkShortenerService_ServiceDesc, linkManagerController,
	)
	linkShortener.AddServer(grpcServer)

	httpServer := pkg.HttpServerProvider(linkShortener)
	gatewayHandler := pkg.GRPCGatewayHandlerProvider(linkShortener)

	httpServer.Router().Handle("/api/", gatewayHandler)
	pkg.AddSwagger(httpServer.Router())

	linkShortener.AddPostHook(pkg.RegisterGatewayHook(linkShortener, gatewayHandler))

	linkShortener.AddServer(httpServer)

	linkShortener.AddPostHook(
		func(ctx context.Context, a *server.App) error {
			a.SetReadyState(true)
			return nil
		},
	)

	return linkShortener, nil
}
