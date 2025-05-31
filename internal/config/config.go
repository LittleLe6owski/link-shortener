package config

// type Config struct {
// 	Cache        cache.Config         `envPrefix:"CACHE_"`
// 	Postgres     PostgresClientConfig `envPrefix:"POSTGRES_"`
// 	DelayOnStart time.Duration        `env:"DELAY_ON_START" envDefault:"0s"`
// }

// type PostgresClientConfig struct {
// 	User                  string        `env:"USER" validate:"required"`
// 	Password              string        `env:"PASSWORD" validate:"required"`
// 	Host                  string        `env:"HOST" validate:"required"`
// 	Port                  string        `env:"PORT" validate:"required"`
// 	DBName                string        `env:"DBNAME" validate:"required"`
// 	SslMode               string        `env:"SSLMODE" validate:"required"`
// 	LogLevel              string        `env:"LOG_LEVEL" envDefault:"INFO" validate:"required,oneof=DEBUG INFO WARN ERROR"`
// 	MetricsUpdateInterval time.Duration `env:"METRICS_UPDATE_INTERVAL" validate:"required" envDefault:"5s"`
// 	PoolMaxConnections    string        `env:"POOL_MAX_CONNS" envDefault:"10"`
// 	BeforeConnect         func(ctx context.Context, config *pgx.ConnConfig) error
// 	AfterConnect          func(ctx context.Context, conn *pgx.Conn) error
// }

// type LogParams struct {
// 	TimestampLayout string `env:"LOG_TIMESTAMP_LAYOUT" envDefault:"2006-01-02T15:04:05Z07:00" validate:"required"`
// 	Prefix          string `env:"LOG_PREFIX"`
// 	EncodingType    string `env:"LOG_ENCODING_TYPE" envDefault:"JSON" validate:"required,oneof=JSON PLAIN"`
// 	LogLevel        string `env:"LOG_LEVEL" envDefault:"INFO" validate:"required,oneof=DEBUG INFO WARN ERROR"`
// 	EnableCaller    bool   `env:"LOG_ENABLE_CALLER"`
// }
