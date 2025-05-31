package postgresql

import "time"

type Config struct {
	Enabled               bool          `env:"ENABLED" envDefault:"true"`
	DSN                   string        `env:"DSN"`
	LogLevel              string        `env:"LOG_LEVEL" envDefault:"INFO" validate:"required,oneof=DEBUG INFO WARN ERROR"`
	MetricsUpdateInterval time.Duration `env:"METRICS_UPDATE_INTERVAL" envDefault:"5s" validate:"required"`
	EnableLogging         bool          `env:"ENABLE_LOGGING" envDefault:"false"`
	MaxConnIdleTime       time.Duration `env:"MAX_CONN_IDLE_TIME" envDefault:"30m"`
	Migration             Migration     `envPrefix:"MIGRATION_"`
	RequestTimeout        time.Duration `env:"REQUEST_TIMEOUT" envDefault:"5s"`
}

type Migration struct {
	Enabled bool   `env:"ENABLED" envDefault:"false"`
	Path    string `env:"PATH"`
}
