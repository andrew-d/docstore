package models

import (
	"encoding/json"
)

type Tag struct {
	// Public fields

	Id   int64
	Name string

	// Private fields

	documentIds []int64
}

// Marshal this tag as JSON
func (t *Tag) MarshalJSON() ([]byte, error) {
	ret := map[string]interface{}{
		"id":   t.Id,
		"name": t.Name,
	}

	if len(t.documentIds) > 0 {
		ret["documents"] = t.documentIds
	}

	return json.Marshal(ret)
}

// Attach the given documents to this tag so they render when marshalling to
// JSON.
func (t *Tag) WithDocuments(docs []Document) *Tag {
	t.documentIds = make([]int64, len(docs))
	for i, doc := range docs {
		t.documentIds[i] = doc.Id
	}
	return t
}
