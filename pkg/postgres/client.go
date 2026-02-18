package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	Pool *pgxpool.Pool
}

type Option func(*pgxpool.Config)


func NewClient(ctx context.Context, connString string, opts ...Option) (*Client, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	for _, opt := range opts {
		opt(config)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Client{Pool: pool}, nil
}

func WithMaxConns(maxConns int32) Option {
	return func(c *pgxpool.Config) {
		c.MaxConns = maxConns
	}
}

func WithMinConns(minConns int32) Option {
	return func(c *pgxpool.Config) {
		c.MinConns = minConns
	}
}

func WithMaxConnLifetime(d time.Duration) Option {
	return func(c *pgxpool.Config) {
		c.MaxConnLifetime = d
	}
}

func WithMaxConnIdleTime(d time.Duration) Option {
	return func(c *pgxpool.Config) {
		c.MaxConnIdleTime = d
	}
}

func (c *Client) Close() {
	c.Pool.Close()
}
