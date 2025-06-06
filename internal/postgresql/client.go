package postgresql

import (
	"context"
	"database/sql"
	"time"

	"github.com/georgysavva/scany/v2/dbscan"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
)

type Connect struct {
	Pgx            *pgxpool.Pool
	ScanAPI        *pgxscan.API
	RequestTimeout time.Duration
}

func NewRepository(ctx context.Context, cfg Config) (*Connect, error) {
	c, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, err
	}

	c.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   &CustomTracer{},
		LogLevel: tracelog.LogLevelDebug,
	}

	c.MaxConnIdleTime = cfg.MaxConnIdleTime
	c.HealthCheckPeriod = cfg.HealthCheckPeriod

	conn, err := pgxpool.NewWithConfig(ctx, c)
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(ctx); err != nil {
		return nil, err
	}

	scanner, err := pgxscan.NewDBScanAPI(dbscan.WithScannableTypes((*sql.Scanner)(nil)))
	if err != nil {
		return nil, err
	}

	scanAPI, err := pgxscan.NewAPI(scanner)
	if err != nil {
		return nil, err
	}

	return &Connect{Pgx: conn, ScanAPI: scanAPI, RequestTimeout: cfg.RequestTimeout}, nil
}
