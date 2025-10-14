package infra

import (
	"context"
	"fmt"

	"go-api-demo/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
)

func NewPostgresPool(ctx context.Context, cfg config.Postgres, debug bool) (*pgxpool.Pool, error) {
	databaseUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?search_path=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Schema,
	)

	config, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		return nil, err
	}

	if debug {
		config.ConnConfig.Tracer = &tracelog.TraceLog{
			Logger: tracelog.LoggerFunc(func(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
				fmt.Printf("[PGX] %s: %s %+v\n", level, msg, data)
			}),
			LogLevel: tracelog.LogLevelDebug,
		}
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
