package postgresql

import "time"

type Config struct {
	DSN               string        `env:"DSN"`
	MaxConnIdleTime   time.Duration `env:"MAX_CONN_IDLE_TIME" envDefault:"30m"`
	HealthCheckPeriod time.Duration `env:"HEALTH_CHECK_PERIOD" envDefault:"1m"`
	RequestTimeout    time.Duration `env:"REQUEST_TIMEOUT" envDefault:"5s"`
}
