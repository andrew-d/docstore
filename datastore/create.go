package datastore

import (
	"github.com/andrew-d/docstore/models"

	"github.com/jmoiron/sqlx"
)

// Create will attempt to create the database schema in the given database.
func Create(db *sqlx.DB) error {
	return transact(db, func(tx *sqlx.Tx) error {
		for _, stmt := range models.Schema() {
			if _, err := tx.Exec(stmt); err != nil {
				return err
			}
		}

		return nil
	})
}

// MustCreate is the same as Create, but panics on error.
func MustCreate(db *sqlx.DB) {
	err := Create(db)
	if err != nil {
		panic(err)
	}
}
