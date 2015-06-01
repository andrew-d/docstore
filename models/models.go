package models

import (
	"time"
)

type Document struct {
	ID          int64     `db:"id" json:"id"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	Files       []File    `db:"-" json:"files"`
	Tags        []Tag     `db:"-" json:"tags"`
}

type File struct {
	ID         int64  `db:"id" json:"id"`
	DocumentID int64  `db:"document_id" json:"document_id"`
	Name       string `db:"name" json:"name"`
	ContentKey string `db:"content_key" json:"content_key"`
}

type Tag struct {
	ID   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

// Other ideas:
//    - Aliases for tags
//    - Collections to group documents
//    - Other... ?
