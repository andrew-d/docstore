package models

type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type File struct {
	ID         int64  `json:"id"`
	Hash       string `json:"hash"`
	DocumentID int64  `json:"document_id"`
}

type Document struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at" db:"created_at"`
	Files     []File `json:"files" db:"-"`
}
