package db

import (
	"context"
	"database/sql"
)

type sqlxContext interface {
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
	GetContext(ctx context.Context, name any, query string, args ...any) error
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}
