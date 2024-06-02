package sqli

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

var ErrNotSupported = errors.New("driver does not support this operation")
var ErrNoRows = pgx.ErrNoRows

type Row interface {
	Scan(dest ...any) error
}

type Rows interface {
	Row
	Close()
	Err() error
	Next() bool
}

type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

type Querier interface {
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) Row
	Exec(ctx context.Context, sql string, args ...any) (Result, error)
}

type Tx interface {
	Querier
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type Conn interface {
	Querier
	Close() error
	Ping() error
	Begin(ctx context.Context) (Tx, error)
}

type DB struct {
	conn Conn
}

func NewDB(conn Conn) *DB {
	return &DB{conn: conn}
}

func (db *DB) Query(ctx context.Context, sql string, args ...any) (Rows, error) {
	return db.conn.Query(ctx, sql, args...)
}

func (db *DB) QueryRow(ctx context.Context, sql string, args ...any) Row {
	return db.conn.QueryRow(ctx, sql, args...)
}

func (db *DB) Exec(ctx context.Context, sql string, args ...any) (Result, error) {
	return db.conn.Exec(ctx, sql, args...)
}

func (db *DB) Begin(ctx context.Context) (Tx, error) {
	return db.conn.Begin(ctx)
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) Ping() error {
	return db.conn.Ping()
}
