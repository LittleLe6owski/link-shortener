package postgresql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/tracelog"
)

type CustomTracer struct{}

func (t *CustomTracer) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	switch level {
	case tracelog.LogLevelTrace:
		log.Printf("TRACE: %s %v\n", msg, data)
	case tracelog.LogLevelDebug:
		log.Printf("DEBUG: %s %v\n", msg, data)
	case tracelog.LogLevelInfo:
		log.Printf("INFO: %s %v\n", msg, data)
	case tracelog.LogLevelWarn:
		log.Printf("WARN: %s %v\n", msg, data)
	case tracelog.LogLevelError:
		log.Printf("ERROR: %s %v\n", msg, data)
	default:
		log.Printf("UNKNOWN: %s %v\n", msg, data)
	}
}
