package redis

import (
	"time"
)

type Config struct {
	Enabled      bool          `env:"ENABLED" envDefault:"true"`
	Hosts        string        `env:"HOSTS"`
	Password     string        `env:"PASSWORD"`
	Database     int           `env:"DATABASE"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT"`
	LogLevel     string        `env:"LOG_LEVEL" envDefault:"INFO" validate:"required,oneof=DEBUG INFO WARN ERROR"`
	TTL          time.Duration `env:"TTL"`
}
