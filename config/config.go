package config

import (
	"time"
)

type Config struct {
	Service    ServiceConfig
	Postgresql PostgresClientConfig `envPrefix:"POSTGRESQL_"`
	Redis      RedisConfig          `envPrefix:"REDIS_"`
	GrpcConfig GRPCConfig           `envPrefix:"GRPC_"`
	HTTPConfig HTTPConfig           `envPrefix:"HTTP_"`
}

type GRPCConfig struct {
	Host            string        `env:"HOST" envDefault:"0.0.0.0" validate:"required"`
	Port            string        `env:"PORT" validate:"required"`
	ShutdownTimeout time.Duration `env:"GRACEFUL_SHUTDOWN_TIMEOUT" envDefault:"20s" validate:"required"`
}

type HTTPConfig struct {
	Host            string        `env:"HOST" envDefault:"0.0.0.0" validate:"required"`
	Port            string        `env:"PORT" validate:"required"`
	ShutdownTimeout time.Duration `env:"GRACEFUL_SHUTDOWN_TIMEOUT" envDefault:"20s" validate:"required"`
}

type ServiceConfig struct {
	Name                  string        `env:"SERVICE_NAME" validate:"required"`
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
	User               string        `env:"USER" validate:"required"`
	Password           string        `env:"PASSWORD" validate:"required"`
	Host               string        `env:"HOST" validate:"required"`
	Port               string        `env:"PORT" validate:"required"`
	DBName             string        `env:"DBNAME" validate:"required"`
	SslMode            string        `env:"SSLMODE" validate:"required"`
	PoolMaxConnections string        `env:"POOL_MAX_CONNS" envDefault:"10"`
	MaxConnIdleTime    time.Duration `env:"MAX_CONN_IDLE_TIME" validate:"required" envDefault:"10m"`
	HealthCheckPeriod  time.Duration `env:"HEALT_CHECK_PERIOD" validate:"required" envDefault:"10m"`
	RequestTimeout     time.Duration `env:"REQUEST_TIMEOUT" validate:"required" envDefault:"4m"`
}

type RedisConfig struct {
	Hosts        string        `env:"HOSTS" validate:"required"`
	Password     string        `env:"PASSWORD" validate:"required"`
	Database     int           `env:"DB" validate:"omitempty" envDefault:"0"`
	PrefixKey    string        `env:"PREFIX" validate:"required" envDefault:"link_shortener"`
	DefaultTTL   time.Duration `env:"DEFAULT_TTL" validate:"required" envDefault:"6m"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" validate:"required" envDefault:"1m"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" validate:"required" envDefault:"1m"`
}
