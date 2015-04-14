package daos

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

func inTransaction(db *sqlx.DB, cb func(tx *sqlx.Tx) error) error {
	finished := false
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if !finished {
			tx.Rollback()
		}
	}()

	err = cb(tx)
	if err != nil {
		return err
	}

	finished = true
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func iQuery(db *sqlx.DB, s string) string {
	if db.DriverName() == "postgres" {
		s = strings.TrimRight(s, "; ") + " RETURNING id"
	}

	return s
}
