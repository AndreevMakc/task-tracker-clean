package postgres

import "time"

type Option func(*Postgres)

func MaxPoolSize(size int) Option {
	return func(pg *Postgres) {
		if size > 0 {
			pg.maxPoolSize = size
		}
	}
}

func ConnAttempts(attempts int) Option {
	return func(pg *Postgres) {
		if attempts > 0 {
			pg.connAttempts = attempts
		}
	}
}

func ConnTimeout(timeout time.Duration) Option {
	return func(pg *Postgres) {
		if timeout > 0 {
			pg.connTimeout = timeout
		}
	}
}
