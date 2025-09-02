package sqlc

import (
	"context"
	"database/sql"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type Queries struct {
	*ReadQueries
	*WriteQueries
	ReadDB  *sql.DB
	WriteDB *sql.DB
}

type ReadQueries struct {
	db DBTX
}

type WriteQueries struct {
	db DBTX
}

func New(readDB, writeDB *sql.DB) *Queries {
	return &Queries{
		ReadQueries:  &ReadQueries{db: readDB},
		WriteQueries: &WriteQueries{db: writeDB},
		ReadDB:       readDB,
		WriteDB:      writeDB,
	}
}
