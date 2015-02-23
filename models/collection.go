package models

import (
	"database/sql"
	"encoding/json"
)

type Collection struct {
	Id           int64         `json:"id"`
	Name         string        `json:"name"`
	Documents    []Document    `json:"documents"`
	CollectionId sql.NullInt64 `json:"collection_id" db:"collection_id"`
}

// Marshal this collection as JSON
func (c *Collection) MarshalJSON() ([]byte, error) {
	ret := map[string]interface{}{
		"id":   c.Id,
		"name": c.Name,
	}

	if c.CollectionId.Valid {
		ret["collection_id"] = c.CollectionId.Int64
	}

	// TODO: documents

	return json.Marshal(ret)
}
