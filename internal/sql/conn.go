package sql

import (
	"context"

	"github.com/leg100/otf/internal/sql/pggen"
)

type (
	Connection struct {
		pggen.Querier // generated queries
		conn
	}
)

// WaitAndLock obtains an exclusive session-level advisory lock. If another
// session holds the lock with the given id then it'll wait until the other
// session releases the lock. The given fn is called once the lock is obtained
// and when the fn finishes the lock is released.
func (db *DB) WaitAndLock(ctx context.Context, id int64, fn func() error) (err error) {
	// A dedicated connection is obtained. Using a connection pool would cause
	// problems because a lock must be released on the same connection on which
	// it was obtained.
	conn, err := db.Pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	if _, err = conn.Exec(ctx, "SELECT pg_advisory_lock($1)", id); err != nil {
		return err
	}
	defer func() {
		_, closeErr := conn.Exec(ctx, "SELECT pg_advisory_unlock($1)", id)
		if err != nil {
			db.Error(err, "unlocking session-level advisory lock")
			return
		}
		err = closeErr
	}()

	if err = fn(); err != nil {
		return err
	}
	return err
}

// Tx provides the caller with a callback in which all operations are conducted
// within a transaction.
func (conn *Connection) Tx(ctx context.Context, callback func(conn) error) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := callback(&Connection{pggen.NewQuerier(tx), conn}); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
