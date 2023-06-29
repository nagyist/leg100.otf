package sql

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-logr/logr"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/leg100/otf/internal/sql/pggen"
)

const defaultMaxConnections = 10 // max conns avail in a pgx pool

type (
	// DB provides access to the postgres db as well as queries generated from
	// SQL
	DB struct {
		*pgxpool.Pool // db connection pool

		logr.Logger
	}

	// Options for constructing a DB
	Options struct {
		Logger     logr.Logger
		ConnString string
	}

	// conn abstracts a postgres connection, which could be a single connection,
	// a pool of connections, or a transaction.
	conn interface {
		Begin(ctx context.Context) (pgx.Tx, error)
		Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
		QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
		Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
		SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	}
)

// New constructs a new DB connection pool, and migrates the schema to the
// latest version.
func New(ctx context.Context, opts Options) (*DB, error) {
	// Bump max number of connections in a pool. By default pgx sets it to the
	// greater of 4 or the num of CPUs. However, otfd acquires several dedicated
	// connections for session-level advisory locks and can easily exhaust this.
	connString, err := setDefaultMaxConnections(opts.ConnString, defaultMaxConnections)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}
	opts.Logger.Info("connected to database", "connstr", connString)

	// goose gets upset with max_pool_conns parameter so pass it the unaltered
	// connection string
	if err := migrate(opts.Logger, opts.ConnString); err != nil {
		return nil, err
	}

	return &DB{
		Pool:   pool,
		Logger: opts.Logger,
	}, nil
}

const txKey = 1

func (db *DB) Conn(ctx context.Context) *Connection {
	if tx := ctx.Value(txKey); tx != nil {
		return tx.(*Connection)
	}
	return &Connection{Querier: pggen.NewQuerier(db.Pool), conn: db.Pool}
}

func setDefaultMaxConnections(connString string, max int) (string, error) {
	// pg connection string can be either a URL or a DSN
	if strings.HasPrefix(connString, "postgres://") || strings.HasPrefix(connString, "postgresql://") {
		u, err := url.Parse(connString)
		if err != nil {
			return "", fmt.Errorf("parsing connection string url: %w", err)
		}
		q := u.Query()
		q.Add("pool_max_conns", strconv.Itoa(max))
		u.RawQuery = q.Encode()
		return url.PathUnescape(u.String())
	} else if connString == "" {
		// presume empty DSN
		return fmt.Sprintf("pool_max_conns=%d", max), nil
	} else {
		// presume non-empty DSN
		return fmt.Sprintf("%s pool_max_conns=%d", connString, max), nil
	}
}
