package daos

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lann/squirrel"

	"github.com/andrew-d/docstore/models"
)

type DocumentDAO struct {
	DB      *sqlx.DB
	Builder squirrel.StatementBuilderType
}

// LoadDocuments will, given a query that returns some documents, run that
// query to get all the documents, and then fetch all related models (i.e.
// tags, collections, etc.) and return them all.
func (d *DocumentDAO) LoadDocuments(sql string, args ...interface{}) (
	documents []models.Document, tags []models.Tag, err error) {

	// Initialize arrays
	documents = []models.Document{}
	tags = []models.Tag{}

	// Run the query
	err = d.DB.Select(&documents, sql, args...)
	if err != nil {
		return
	}

	if len(documents) == 0 {
		return
	}

	// Get a list of document IDs
	documentIds := make([]int64, len(documents))
	for i, doc := range documents {
		documentIds[i] = doc.Id
	}

	// Fetch all tags for this set of documents
	var tagSpecs []struct {
		TagId      int64  `db:"tag_id"`
		Name       string `db:"name"`
		DocumentId int64  `db:"document_id"`
	}
	sql, args, _ = (d.Builder.
		Select("tags.id AS tag_id", "tags.name AS name", "documents.id AS document_id").
		From("tags").
		Join("document_tags ON document_tags.tag_id == tags.id").
		Join("documents ON document_tags.document_id == documents.id").
		Where(squirrel.Eq{"documents.id": documentIds}).
		ToSql())
	err = d.DB.Select(&tagSpecs, sql, args...)
	if err != nil {
		return
	}

	// Create both a list of all tags (for rendering), and a mapping of document
	// ID to tag IDs for that document (for fast tag setup, below).
	tags = make([]models.Tag, 0, len(tagSpecs))
	tagMap := make(map[int64][]int64)
	for _, spec := range tagSpecs {
		currTag := models.Tag{
			Id:   spec.TagId,
			Name: spec.Name,
		}
		tagMap[spec.DocumentId] = append(tagMap[spec.DocumentId], currTag.Id)
		tags = append(tags, currTag)
	}

	// For each document in our 'documents' array, we add the tags for it.
	empty := []int64{}
	for i := 0; i < len(documents); i++ {
		tagsForDoc := tagMap[documents[i].Id]
		if len(tagsForDoc) > 0 {
			documents[i].Tags = tagsForDoc
		} else {
			documents[i].Tags = empty
		}
	}

	// TODO: load collections for documents

	// All done!
	return
}

// LoadDocuments will, given a query that returns a document, run that query
// to get the document, and then fetch all related models (i.e. tags,
// collection, etc.) and return them all.  It will return a document with a
// zero ID if the document was not found.
func (d *DocumentDAO) LoadDocument(sql string, args ...interface{}) (
	document models.Document, tags []models.Tag, err error) {

	// Fetch this document
	err = d.DB.Get(&document, sql, args...)
	if document.Id == 0 {
		return
	}

	// Initialize arrays
	tags = []models.Tag{}

	// Load all tags for this document.
	sql, args, _ = (d.Builder.
		Select("tags.id AS id", "tags.name AS name").
		From("tags").
		Join("document_tags ON document_tags.tag_id == tags.id").
		Join("documents ON document_tags.document_id == documents.id").
		Where(squirrel.Eq{"documents.id": document.Id}).
		ToSql())
	err = d.DB.Select(&tags, sql, args...)
	if err != nil {
		return
	}

	// Set the tags on this document.
	document.Tags = make([]int64, len(tags))
	for i, tag := range tags {
		document.Tags[i] = tag.Id
	}

	// TODO: load this document's collection

	// All done!
	return
}

// CreateDocument will create a new document with the given parameters.  It
// will return an error if there is a problem creating the document.
//
// Note: pass a collectionId of 0 to indicate 'do not add to a collection'.
// Note: all operations are performed in a transaction - i.e. everything will
//       succeed, or nothing will.
func (d *DocumentDAO) CreateDocument(name string, tagIds []int64, collectionId int64) (
	document models.Document, tags []models.Tag, err error) {
	err = inTransaction(d.DB, func(tx *sqlx.Tx) error {

		// Create the new document model
		document = models.Document{
			Name:      name,
			CreatedAt: time.Now().UTC().UnixNano(),
		}

		// If the collection ID is non-zero, we add it to that collection
		if collectionId > 0 {
			var coll models.Collection

			sql, args, _ := (d.Builder.
				Select("*").
				From("collections").
				Where(squirrel.Eq{"id": collectionId}).
				ToSql())
			err := d.DB.Get(&coll, sql, args...)
			if err != nil {
				return err
			}
			if coll.Id == 0 {
				return err
			}

			document.CollectionId.Int64 = coll.Id
			document.CollectionId.Valid = true
		}

		// Insert the document
		sql, args, _ := (d.Builder.
			Insert("documents").
			Columns("name", "created_at", "collection_id").
			Values(document.Name, document.CreatedAt, document.CollectionId).
			ToSql())
		sqlRes, err := tx.Exec(iQuery(d.DB, sql), args...)
		if err != nil {
			return err
		}

		// Save the returned ID
		id, _ := sqlRes.LastInsertId()
		document.Id = id

		// Insert all tags
		for _, tag_id := range tagIds {
			tag := models.Tag{}

			// Check the tag exists
			sql, args, _ := (d.Builder.
				Select("*").
				From("tags").
				Where(squirrel.Eq{"id": id}).
				ToSql())
			err = tx.Get(&tag, sql, args...)
			if err != nil {
				return err
			}
			if tag.Id == 0 {
				return err
			}

			// Save for rendering
			tags = append(tags, tag)

			// Insert it
			sql, args, _ = (d.Builder.
				Insert("document_tags").
				Columns("document_id", "tag_id").
				Values(id, tag_id).
				ToSql())
			_, err := tx.Exec(sql, args...)
			if err != nil {
				return err
			}
		}

		// Update the returned document with the tag IDs
		document.Tags = make([]int64, len(tags))
		for i, tag := range tags {
			document.Tags[i] = tag.Id
		}

		// All good!
		return nil
	})

	// TODO: return collection

	return
}
