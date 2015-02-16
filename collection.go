package main

import (
	"database/sql"
)

type Collection struct {
	Id           int64         `json:"id"`
	Name         string        `json:"name"`
	Documents    []Document    `json:"documents"`
	CollectionId sql.NullInt64 `json:"collection_id" db:"collection_id"`
}
