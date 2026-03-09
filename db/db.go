package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Settings struct {
	Conn            string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

type Database struct {
	db *sqlx.DB
}

func NewDatabase(settings Settings) *Database {
	db, err := sqlx.Connect("pgx", settings.Conn)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to db: %v", err))
	}

	db.SetMaxOpenConns(settings.MaxOpenConns)
	db.SetMaxIdleConns(settings.MaxIdleConns)
	db.SetConnMaxLifetime(settings.ConnMaxLifetime)

	return &Database{db: db}
}

func (d *Database) Select(ctx context.Context, name string, dest any, query string, args ...any) error {
	return d.wrapError(name, d.db.SelectContext(ctx, dest, query, args...))
}

func (d *Database) Get(ctx context.Context, name string, dest any, query string, args ...any) error {
	return d.wrapError(name, d.db.GetContext(ctx, dest, query, args...))
}

func (d *Database) Exec(ctx context.Context, name string, query string, args ...any) (sql.Result, error) {
	res, err := d.db.ExecContext(ctx, query, args...)
	return res, d.wrapError(name, err)
}

func (d *Database) wrapError(name string, err error) error {
	if err != nil {
		return fmt.Errorf("failed execute query '%s': %w", name, err)
	}
	return nil
}
