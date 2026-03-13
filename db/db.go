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
	pool *sqlx.DB
}

func NewDatabase(settings Settings) *Database {
	db, err := sqlx.Connect("pgx", settings.Conn)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to db: %v", err))
	}

	db.SetMaxOpenConns(settings.MaxOpenConns)
	db.SetMaxIdleConns(settings.MaxIdleConns)
	db.SetConnMaxLifetime(settings.ConnMaxLifetime)

	return &Database{pool: db}
}

func (d *Database) db(ctx context.Context) sqlxContext {
	if tx, ok := ctx.Value(txKey{}).(*sqlx.Tx); ok {
		return tx
	}
	return d.pool
}

func (d *Database) Select(ctx context.Context, name string, dest any, query string, args ...any) error {
	return d.wrapError(name, d.db(ctx).SelectContext(ctx, dest, query, args...))
}

func (d *Database) Get(ctx context.Context, name string, dest any, query string, args ...any) error {
	return d.wrapError(name, d.db(ctx).GetContext(ctx, dest, query, args...))
}

func (d *Database) Exec(ctx context.Context, name string, query string, args ...any) (sql.Result, error) {
	res, err := d.db(ctx).ExecContext(ctx, query, args...)
	return res, d.wrapError(name, err)
}

func (d *Database) wrapError(name string, err error) error {
	if err != nil {
		return fmt.Errorf("failed execute query '%s': %w", name, err)
	}
	return nil
}

type txKey struct{}

func (d *Database) TransactionV2(ctx context.Context, fn func(ctx context.Context) error) error {
	if _, ok := ctx.Value(txKey{}).(*sqlx.Tx); ok {
		return fn(ctx)
	}

	tx, err := d.pool.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	txCtx := context.WithValue(ctx, txKey{}, tx)

	err = fn(txCtx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("error on rollback: %v (original error: %w)", rbErr, err)
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
