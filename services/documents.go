package services

import (
	"github.com/andrew-d/docstore/models"
)

// DocumentsService is an abstraction over methods operating on documents.
type DocumentsService interface {
	// Get a document by ID.
	Get(id int) (*models.Document, error)

	// List all documents.
	List(opts *DocumentsListOptions) ([]*models.Document, error)

	// Create a new document and all associated files.  This function will set the
	// ID of the document to the newly-inserted ID on success.  It will return an
	// error if one of the given files already exists.
	Create(doc *models.Document) error
}

type DocumentsListOptions struct {
	ListOptions
}
