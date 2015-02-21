package models

type Document struct {
	Id           int64 `json:"id"`
	CreatedAt    int64 `json:"created_at" db:"created_at"`
	CollectionId int64 `json:"collection_id" db:"collection_id"`
}

// Helper method to return a JSON-serializable struct of the document with the
// given tags.
func (d *Document) WithTags(tags []Tag) interface{} {
	// Get IDs
	tagIds := make([]int64, len(tags))
	for i, tag := range tags {
		tagIds[i] = tag.Id
	}

	// Make document
	type returnType struct {
		Id           int64   `json:"id"`
		CreatedAt    int64   `json:"created_at"`
		CollectionId int64   `json:"collection_id"`
		Tags         []int64 `json:"tags"`
	}

	ret := &returnType{
		Id:           d.Id,
		CreatedAt:    d.CreatedAt,
		CollectionId: d.CollectionId,
		Tags:         tagIds,
	}
	return ret
}
