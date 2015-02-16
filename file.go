package main

type File struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Hash       string `json:"hash"`
	DocumentId int64  `json:"document_id" db:"document_id"`
}
