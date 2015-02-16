package main

type Collection struct {
	Id           int64      `json:"id"`
	Name         string     `json:"name" sql:"not null"`
	Documents    []Document `json:"documents"`
	CollectionId int64      `json:"collection_id"`

	// TODO: return collection vs. collection id?
}
