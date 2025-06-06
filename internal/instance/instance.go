package instance

import (
	"context"
	"errors"
	"fmt"

	"github.com/LittleLe6owski/link-shortener/internal/config"
	"github.com/LittleLe6owski/link-shortener/internal/postgresql"
	"github.com/LittleLe6owski/link-shortener/internal/usecase"
	"github.com/LittleLe6owski/link-shortener/pkg"
)

var (
	ErrCannotParseConfig    = errors.New("unable to get messaging config")
	ErrWrongPostgresConn    = errors.New("unable to init Postgres connection")
	ErrWrongMongoConn       = errors.New("unable to init MongoDB connection")
	ErrWrongMongoCollection = errors.New("unable to get MongoDB collection")
	ErrWrongRedisConn       = errors.New("unable to init Redis connection")
	ErrWrongProducerConn    = errors.New("unable to init Kafka producer")
	ErrWrongProducerTopic   = errors.New("unable to init Kafka producer for topic")
	ErrWrongConsumerConn    = errors.New("unable to init Kafka consumer")
	ErrWrongConsumerTopic   = errors.New("unable to init Kafka consumer for topic")
)

func New(cfg config.Config) (*pkg.App, error) {
	linkShortener, err := pkg.NewApp(cfg)
	if err != nil {
		return nil, err
	}

	// imRepo, err := InitIMRepo(draftProject)
	// if err != nil {
	// 	return nil, err
	// }

	// pgRepo, err := InitPgRepo(draftProject)
	// if err != nil {
	// 	return nil, err
	// }

	// cacheMyItem, err := InitRedisCache(draftProject)
	// if err != nil {
	// 	return nil, err
	// }

	// apiManager := apimanager.NewMyItemManager(
	// 	imRepo,
	// 	pgRepo,
	// 	mongoRepoMyItem,
	// 	cacheMyItem,
	// 	producerMyItem,
	// 	draftProject.Logger(),
	// )

	// queueManager := queuemanager.NewMyItemManager(
	// 	imRepo,
	// 	pgRepo,
	// 	mongoRepoMyItem,
	// 	draftProject.Logger(),
	// )

	// messagingConsumer, err := InitMyItemConsumer(draftProject, queueManager)
	// if err != nil {
	// 	return nil, fmt.Errorf("%w: %v", ErrWrongConsumerConn, err)
	// }

	// draftProject.AddServer(messagingConsumer.(server.Server))

	// grpcServer := app.GRPCProvider(draftProject)
	// grpcServer.RegisterService(
	// 	&apiv1.DraftProjectService_ServiceDesc,
	// 	api.NewDraftProjectServiceServer(apiManager),
	// )
	// draftProject.AddServer(grpcServer)

	// httpServer := app.HttpServerProvider(draftProject)
	// gatewayHandler := app.GRPCGatewayHandlerProvider(draftProject)
	// httpServer.Router().Handle("/api/", gatewayHandler)
	// draftProject.AddPreHook(RegisterGatewayHook(gatewayHandler))
	// draftProject.AddPreHook(app.RegisterTracerHook())
	// draftProject.AddServer(httpServer)

	// metricsServer := app.MetricsProvider(draftProject)
	// draftProject.AddServer(metricsServer)

	// draftProject.AddPostHook(func(ctx context.Context, a *app.App) error {
	// 	a.SetReadyState(true)
	// 	return nil
	// })

	return linkShortener, nil
}

// func InitRedisCache(app *app.App) (usecase.MyItemCacheManager, error) {
// 	cfg, ok := app.Config().Service.ExtraConfig.(*config.Config)
// 	if !ok {
// 		return nil, ErrCannotParseExtraConfig
// 	}

// 	redisCache, err := cache.NewRedisCache(context.Background(), &cfg.Cache.Config, app.Logger(), collectorName)
// 	if err != nil {
// 		return nil, fmt.Errorf("%w: %v", ErrWrongRedisConn, err)
// 	}

// 	return redisCache.NewMyItemRepository(cfg.Cache, cfg.MaxPageSize), nil
// }

// func InitIMRepo(app *app.App) (usecase.MyItemStorageManager, error) {
// 	cfg, ok := app.Config().Service.ExtraConfig.(*config.Config)
// 	if !ok {
// 		return nil, ErrCannotParseExtraConfig
// 	}

// 	return inmemory.NewMyItemRepository(cfg.InMemory, app.Logger(), cfg.MaxPageSize), nil
// }

func InitPgRepo(ctx context.Context, app *pkg.App) (usecase.MyItemStorageManager, error) {

	pgr, err := postgresql.NewRepository(ctx, postgresql.Config{})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrWrongPostgresConn, err)
	}

	return postgres.NewMyItemPostgresRepository(pgr, cfg.MaxPageSize), nil
}
