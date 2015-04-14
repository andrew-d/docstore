package models

import (
	"database/sql"
)

// TODO: Rename CollectionId ==> Collection?

type Document struct {
	// Database fields

	Id           int64         `json:"id"`
	Name         string        `json:"name"`
	CreatedAt    int64         `json:"created_at" db:"created_at"`
	CollectionId sql.NullInt64 `json:"collection_id" db:"collection_id"`

	// Non-database fields

	Tags  []int64 `json:"tags" db:"-"`
	Files []int64 `json:"files" db:"-"`
}
