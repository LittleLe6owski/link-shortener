package postgresql

import (
	"context"
	"database/sql"
	"log"

	"github.com/georgysavva/scany/v2/dbscan"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewConnect(
	ctx context.Context,
	cfg Config,
	logger *log.Logger,
	col metrics.MetricIface,
) (*Connect, error) {
	c, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, err
	}

	c.ConnConfig.Tracer = &combinedTracer{
		logger:        logger,
		level:         cfg.LogLevel,
		col:           col,
		otel:          otelpgx.NewTracer(),
		enableLogging: cfg.EnableLogging,
		dbName:        c.ConnConfig.Database,
	}

	c.MaxConnIdleTime = cfg.MaxConnIdleTime

	conn, err := pgxpool.NewWithConfig(ctx, c)
	if err != nil {
		return nil, err
	}

	err = conn.Ping(ctx)
	if err != nil {
		return nil, err
	}

	col.StartPoolStatsMonitoring(ctx, conn, cfg.MetricsUpdateInterval, c.ConnConfig.Database)

	scanAPI, err := newScanApi()
	if err != nil {
		return nil, err
	}

	return &Connect{
		Pgx:            pgxConn,
		ScanAPI:        scanAPI,
		RequestTimeout: cfg.RequestTimeout,
	}, nil
}

func newScanApi() (*pgxscan.API, error) {
	scanner, err := pgxscan.NewDBScanAPI(
		dbscan.WithScannableTypes((*sql.Scanner)(nil)),
	)
	if err != nil {
		return nil, err
	}

	return pgxscan.NewAPI(scanner)
}
