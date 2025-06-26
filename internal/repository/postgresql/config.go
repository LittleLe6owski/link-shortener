package postgresql

import "time"

type Config struct {
	DSN               string        
	MaxConnIdleTime   time.Duration 
	HealthCheckPeriod time.Duration 
	RequestTimeout    time.Duration 
}
