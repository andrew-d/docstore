package models

type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Document struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at" db:"created_at"`
}
