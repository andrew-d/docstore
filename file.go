package main

type File struct {
	Id         int64  `json:"id"`
	Name       string `json:"name" sql:"not null"`
	Hash       string `json:"hash" sql:"not null;unique"`
	DocumentId int64  `json:"document_id"` // foreign key to Document
}
