package config

import (
	"time"

	"gitlab.bronevik.space/bronevik/libs/bro/pkg/server"
)

type BaseConfig struct {
	Name                  string        `env:"SERVICE_NAME,required" validate:"required"`
	GracefulShutdownDelay time.Duration `env:"SERVICE_GRACEFUL_SHUTDOWN_DELAY" envDefault:"5s" validate:"required"`
	GRPC                  server.Config `envPrefix:"GRPC_"`
	HTTP                  server.Config `envPrefix:"HTTP_"`
	Profiling             bool          `env:"SERVICE_PROFILING_ENABLE"`
}
