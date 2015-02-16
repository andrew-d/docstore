package models

type Document struct {
	Id           int64 `json:"id"`
	CreatedAt    int64 `json:"created_at" db:"created_at"`
	CollectionId int64 `json:"collection_id" db:"collection_id"`
}
