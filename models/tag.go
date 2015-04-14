package models

type Tag struct {
	// Database fields

	Id   int64  `json:"id"`
	Name string `json:"name"`

	// Non-database fields

	Documents []int64 `json:"documents" db:"-"`
}
