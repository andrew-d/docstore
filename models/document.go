package models

import (
	"database/sql"
	"encoding/json"
)

// TODO: Rename CollectionId ==> Collection?

type Document struct {
	// Database fields

	Id           int64
	Name         string
	CreatedAt    int64         `db:"created_at"`
	CollectionId sql.NullInt64 `db:"collection_id"`

	// Private fields

	tagIds  []int64
	fileIds []int64
}

// Marshal this document as JSON
func (d *Document) MarshalJSON() ([]byte, error) {
	ret := map[string]interface{}{
		"id":         d.Id,
		"name":       d.Name,
		"created_at": d.CreatedAt,
	}

	if d.CollectionId.Valid {
		ret["collection_id"] = d.CollectionId.Int64
	}

	if len(d.tagIds) > 0 {
		ret["tags"] = d.tagIds
	}

	if len(d.fileIds) > 0 {
		ret["files"] = d.fileIds
	}

	return json.Marshal(ret)
}

// Helper method to attach the given tags to this document, so they render when
// marshalling to JSON.
func (d *Document) WithTags(tags []Tag) *Document {
	d.tagIds = make([]int64, len(tags))
	for i, tag := range tags {
		d.tagIds[i] = tag.Id
	}
	return d
}

// Helper method to attach the given files to this document, so they render
// when marshalling to JSON.
func (d *Document) WithFiles(files []File) *Document {
	d.fileIds = make([]int64, len(files))
	for i, file := range files {
		d.fileIds[i] = file.Id
	}
	return d
}
