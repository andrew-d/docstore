package datastore

import (
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/andrew-d/docstore/services"
)

func rid(db *sqlx.DB, s string) string {
	if db.DriverName() != "postgres" {
		return s
	}

	if strings.HasSuffix(s, ";") {
		s = s[0 : len(s)-2]
	}

	return s + `RETURNING id;`
}

// M is a helper type that is an alias for map[string]interface{}
type M map[string]interface{}

// transact calls the given function in a DB transaction.  It will roll back
// the transaction on failure and commit on success.
func transact(db *sqlx.DB, fn func(tx *sqlx.Tx) error) error {
	var (
		tx  *sqlx.Tx
		err error
	)

	if tx, err = db.Beginx(); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = fn(tx); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// dbretry will call the given function until it either returns nil (no error),
// or the limit is reached.  It will return the final error returned if the
// limit is exceeded.  No retry will be performed if the error is a 'Not Found'
// or 'Exists' error.
func dbretry(limit int, fn func() error) error {
	var err error

	for i := 0; i < limit; i++ {
		if err = fn(); err == nil {
			return nil
		}

		if services.IsNotFound(err) || services.IsExists(err) {
			return err
		}
	}

	return err
}
