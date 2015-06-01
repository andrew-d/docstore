package models

var databaseSchema = []string{
	// Create tables
	`
CREATE TABLE IF NOT EXISTS documents (
	id            INTEGER PRIMARY KEY,
	description   TEXT NOT NULL,
	created_at    INTEGER NOT NULL
);`,
	`
CREATE TABLE IF NOT EXISTS tags (
	id   INTEGER PRIMARY KEY,
	name VARCHAR(255) NOT NULL UNIQUE
)`,
	`
CREATE TABLE IF NOT EXISTS files (
	id          INTEGER PRIMARY KEY,
	document_id INTEGER NOT NULL,
	content_key VARCHAR(32) NOT NULL UNIQUE,
	name        VARCHAR(255) NOT NULL,
	FOREIGN KEY (document_id) REFERENCES documents(id)
)`,

	// TODO: create indexes
}

func Schema() []string {
	return databaseSchema
}
