package models

import (
	"time"
)

// First pass at an API is below
// --------------------------------------------------

type Document struct {
	ID          int
	Description string
	CreatedAt   time.Time
	Files       []File
	Tags        []Tag `gorm:"many2many:document_tags"`
}

type File struct {
	ID         int
	ContentKey string `sql:"not null;unique"` // for castore
}

type Tag struct {
	ID   int
	Name string `sql:"not null;unique"`
}

// Other ideas:
//    - Aliases for tags
//    - Collections to group documents
//    - Other... ?
