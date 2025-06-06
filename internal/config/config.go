package config

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Config struct {
	Postgresql PostgresClientConfig `envPrefix:"POSTGRESQL_"`
	Cache      RedisConfig          `envPrefix:"REDIS_"`
	Service    ServiceConfig        `envPrefix:"SERVICE_"`
	LogConfig  LogConfig            `env:"LOG_"`
}

type ServiceConfig struct {
	Name                  string        `env:"SERVICE_NAME,required" validate:"required"`
	GracefulShutdownDelay time.Duration `env:"GRACEFUL_SHUTDOWN" validate:"required" envDefault:"1m"`
	GRPC                  ServerConfig  `envPrefix:"GRPC_"`
	HTTP                  ServerConfig  `envPrefix:"HTTP_"`
	Profiling             bool          `env:"SERVICE_PROFILING_ENABLE"`
}

type ServerConfig struct {
	Host            string        `env:"HOST" envDefault:"0.0.0.0" validate:"required"`
	Port            string        `env:"PORT" validate:"required"`
	ShutdownTimeout time.Duration `env:"GRACEFUL_SHUTDOWN_TIMEOUT" envDefault:"20s" validate:"required"`
}

type PostgresClientConfig struct {
	User                  string        `env:"USER" validate:"required"`
	Password              string        `env:"PASSWORD" validate:"required"`
	Host                  string        `env:"HOST" validate:"required"`
	Port                  string        `env:"PORT" validate:"required"`
	DBName                string        `env:"DBNAME" validate:"required"`
	SslMode               string        `env:"SSLMODE" validate:"required"`
	LogLevel              string        `env:"LOG_LEVEL" envDefault:"INFO" validate:"required,oneof=DEBUG INFO WARN ERROR"`
	MetricsUpdateInterval time.Duration `env:"METRICS_UPDATE_INTERVAL" validate:"required" envDefault:"5s"`
	PoolMaxConnections    string        `env:"POOL_MAX_CONNS" envDefault:"10"`
	BeforeConnect         func(ctx context.Context, config *pgx.ConnConfig) error
	AfterConnect          func(ctx context.Context, conn *pgx.Conn) error
}

type RedisConfig struct {
}

type LogConfig struct {
	TimestampLayout string `env:"TIMESTAMP_LAYOUT" envDefault:"2006-01-02T15:04:05Z07:00" validate:"required"`
	Prefix          string `env:"PREFIX"`
	EncodingType    string `env:"ENCODING_TYPE" envDefault:"JSON" validate:"required,oneof=JSON PLAIN"`
	LogLevel        string `env:"LEVEL" envDefault:"INFO" validate:"required,oneof=DEBUG INFO WARN ERROR"`
	EnableCaller    bool   `env:"ENABLE_CALLER"`
}
