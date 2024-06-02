package sqli

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type pgResult struct {
	cmdTag pgconn.CommandTag
}

func (pr *pgResult) LastInsertId() (int64, error) {
	return 0, ErrNotSupported
}

func (pr *pgResult) RowsAffected() (int64, error) {
	return pr.cmdTag.RowsAffected(), nil
}

func NewPostgreConn(ctx context.Context, connUrl string) (Conn, error) {
	conn, err := pgx.Connect(ctx, connUrl)
	if err != nil {
		return nil, err
	}

	return newPgxWrapper(conn), nil
}

type pgxWrapper struct {
	conn *pgx.Conn
}

func newPgxWrapper(conn *pgx.Conn) *pgxWrapper {
	return &pgxWrapper{conn: conn}
}

func (p *pgxWrapper) Query(ctx context.Context, sql string, args ...any) (Rows, error) {
	r, err := p.conn.Query(ctx, sql, args...)

	return r, err
}

func (p *pgxWrapper) QueryRow(ctx context.Context, sql string, args ...any) Row {
	r := p.conn.QueryRow(ctx, sql, args...)

	return r
}

func (p *pgxWrapper) Exec(ctx context.Context, sql string, args ...any) (Result, error) {
	r, err := p.conn.Exec(ctx, sql, args...)

	return &pgResult{cmdTag: r}, err
}

func (p *pgxWrapper) Close() error {
	return p.conn.Close(context.Background())
}

func (p *pgxWrapper) Ping() error {
	return p.conn.Ping(context.Background())
}

func (p *pgxWrapper) Begin(ctx context.Context) (Tx, error) {
	tx, err := p.conn.Begin(ctx)
	if err != nil {
		return nil, err
	}

	txWrap := pgxTxWrapper{tx: tx}

	return &txWrap, nil
}

type pgxTxWrapper struct {
	tx pgx.Tx
}

func (px *pgxTxWrapper) Query(ctx context.Context, sql string, args ...any) (Rows, error) {
	return px.tx.Query(ctx, sql, args...)
}

func (px *pgxTxWrapper) QueryRow(ctx context.Context, sql string, args ...any) Row {
	return px.tx.QueryRow(ctx, sql, args...)
}

func (px *pgxTxWrapper) Exec(ctx context.Context, sql string, args ...any) (Result, error) {
	r, err := px.tx.Exec(ctx, sql, args...)

	return &pgResult{cmdTag: r}, err
}

func (px *pgxTxWrapper) Commit(ctx context.Context) error {
	return px.tx.Commit(ctx)
}

func (px *pgxTxWrapper) Rollback(ctx context.Context) error {
	return px.tx.Rollback(ctx)
}
