package postgres

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresOptionTestCase struct {
	Name     string
	Option   Option
	Validate func(*pgxpool.Config) bool
}

func TestOptions(t *testing.T) {
	defaultConfig, _ := pgxpool.ParseConfig("postgres://user:password@localhost:5432/dbname")

	tests := []PostgresOptionTestCase{
		{
			Name:   "WithMaxConns",
			Option: WithMaxConns(10),
			Validate: func(c *pgxpool.Config) bool {
				return c.MaxConns == 10
			},
		},
		{
			Name:   "WithMinConns",
			Option: WithMinConns(5),
			Validate: func(c *pgxpool.Config) bool {
				return c.MinConns == 5
			},
		},
		{
			Name:   "WithMaxConnLifetime",
			Option: WithMaxConnLifetime(time.Hour),
			Validate: func(c *pgxpool.Config) bool {
				return c.MaxConnLifetime == time.Hour
			},
		},
		{
			Name:   "WithMaxConnIdleTime",
			Option: WithMaxConnIdleTime(time.Minute),
			Validate: func(c *pgxpool.Config) bool {
				return c.MaxConnIdleTime == time.Minute
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// Create a copy of default config for each test
			config := defaultConfig.Copy()
			tt.Option(config)
			
			if !tt.Validate(config) {
				t.Errorf("Option %s failed to set configuration correctly", tt.Name)
			}
		})
	}
}
