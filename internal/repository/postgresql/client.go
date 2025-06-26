package postgresql

import (
	"context"
	"database/sql"
	"time"

	"github.com/georgysavva/scany/v2/dbscan"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Connect struct {
	Pgx            *pgxpool.Pool
	ScanAPI        *pgxscan.API
	RequestTimeout time.Duration
}

func NewConnector(ctx context.Context, cfg Config) (Connect, error) {
	c, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return Connect{}, err
	}

	c.MaxConnIdleTime = cfg.MaxConnIdleTime
	c.HealthCheckPeriod = cfg.HealthCheckPeriod

	conn, err := pgxpool.NewWithConfig(ctx, c)
	if err != nil {
		return Connect{}, err
	}

	if err = conn.Ping(ctx); err != nil {
		return Connect{}, err
	}

	scanner, err := pgxscan.NewDBScanAPI(dbscan.WithScannableTypes((*sql.Scanner)(nil)))
	if err != nil {
		return Connect{}, err
	}

	scanAPI, err := pgxscan.NewAPI(scanner)
	if err != nil {
		return Connect{}, err
	}

	return Connect{Pgx: conn, ScanAPI: scanAPI, RequestTimeout: cfg.RequestTimeout}, nil
}
