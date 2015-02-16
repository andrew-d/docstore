package main

// This is a list of strings, instead of one big string, because SQLite will
// not execute multiple statements in a single call.
var databaseSchema = []string{
	`
CREATE TABLE IF NOT EXISTS collections (
	id            INTEGER PRIMARY KEY,
	name          VARCHAR(255) NOT NULL,
	collection_id INTEGER,

	FOREIGN KEY (collection_id) REFERENCES collections(id)
)`,
	`
CREATE TABLE IF NOT EXISTS documents (
	id            INTEGER PRIMARY KEY,
	created_at    INTEGER NOT NULL,
	collection_id INTEGER,

	FOREIGN KEY (collection_id) REFERENCES collections(id)
)`,
	`
CREATE TABLE IF NOT EXISTS tags (
	id   INTEGER PRIMARY KEY,
	name VARCHAR(255) NOT NULL UNIQUE
)`,
	`
CREATE TABLE IF NOT EXISTS files (
	id   INTEGER PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	hash VARCHAR(32) NOT NULL UNIQUE
)`,
	`
CREATE TABLE IF NOT EXISTS document_tags (
	document_id INTEGER NOT NULL,
	tag_id      INTEGER NOT NULL,

	FOREIGN KEY (document_id) REFERENCES documents(id),
	FOREIGN KEY (tag_id) REFERENCES tags(id)
)`,
	`
CREATE TABLE IF NOT EXISTS collection_documents (
	collection_id INTEGER NOT NULL,
	document_id   INTEGER NOT NULL,

	FOREIGN KEY (collection_id) REFERENCES collections(id),
	FOREIGN KEY (document_id) REFERENCES documents(id)
)`,
}

// TODO:
//	- Indexes
//	- Primary key for document_tags
