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
	ctxWithName := context.WithValue(ctx, "name", name)
	return d.db.SelectContext(ctxWithName, dest, query, args...)
}

func (d *Database) Get(ctx context.Context, name string, dest any, query string, args ...any) error {
	ctxWithName := context.WithValue(ctx, "name", name)
	return d.db.GetContext(ctxWithName, dest, query, args...)
}

func (d *Database) Exec(ctx context.Context, name string, query string, args ...any) (sql.Result, error) {
	ctxWithName := context.WithValue(ctx, "name", name)
	return d.db.ExecContext(ctxWithName, query, args...)
}
