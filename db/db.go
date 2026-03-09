package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Settings struct {
	Conn            string
	MaxIdleConns    int // В pgxpool это MinConns
	MaxOpenConns    int // В pgxpool это MaxConns
	ConnMaxLifetime time.Duration
}

type Database struct {
	Pool *pgxpool.Pool
}

func NewDatabase(settings Settings) *Database {
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(settings.Conn)
	if err != nil {
		panic(fmt.Sprintf("failed to parse connection string: %v", err))
	}

	config.MaxConns = int32(settings.MaxOpenConns)
	config.MinConns = int32(settings.MaxIdleConns)
	config.MaxConnLifetime = settings.ConnMaxLifetime

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to db: %v", err))
	}

	return &Database{Pool: pool}
}

func Select[T any](ctx context.Context, d *Database, name string, query string, args ...any) ([]T, error) {
	ctxWithName := context.WithValue(ctx, "name", name)
	rows, err := d.Pool.Query(ctxWithName, query, args...)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[T])
}

func Get[T any](ctx context.Context, d *Database, name string, query string, args ...any) (*T, error) {
	ctxWithName := context.WithValue(ctx, "name", name)
	rows, err := d.Pool.Query(ctxWithName, query, args...)
	if err != nil {
		return nil, err
	}

	r, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[T])
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (d *Database) Exec(ctx context.Context, name string, query string, args ...any) (pgconn.CommandTag, error) {
	ctxWithName := context.WithValue(ctx, "name", name)
	return d.Pool.Exec(ctxWithName, query, args...)
}

func (d *Database) Close() {
	d.Pool.Close()
}
