package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	_defaultMaxPoolSize  = 10
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

type Postgres struct {
	pool *pgxpool.Pool

	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration
}

func New(ctx context.Context, url string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(pg)
	}

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("postgres - New - pgxpool.ParseConfig: %w", err)
	}

	poolConfig.MaxConns = int32(pg.maxPoolSize)

	var pool *pgxpool.Pool
	var lastErr error

	for pg.connAttempts > 0 {
		ctxTimeout, cancel := context.WithTimeout(ctx, pg.connTimeout)

		pool, lastErr = pgxpool.NewWithConfig(ctxTimeout, poolConfig)
		if lastErr == nil {
			if lastErr = pool.Ping(ctxTimeout); lastErr == nil {
				cancel()
				break
			}
		}

		cancel()

		pg.connAttempts--
		if pg.connAttempts == 0 {
			break
		}

		select {
		case <-time.After(pg.connTimeout):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	if lastErr != nil {
		return nil, fmt.Errorf("postgres - New - attempts exhausted: %w", lastErr)
	}

	pg.pool = pool

	return pg, nil
}

func (p *Postgres) Pool() *pgxpool.Pool {
	return p.pool
}

func (p *Postgres) Close() {
	if p.pool != nil {
		p.pool.Close()
	}
}
