package store

import (
	"path/filepath"

	"github.com/andrew-d/castore"

	"github.com/andrew-d/docstore/store/models"
)

var _ = models.Document{}

type Store struct {
	content *castore.CAStore
}

func New(root string) (*Store, error) {
	content, err := castore.New(castore.Options{
		BasePath: filepath.Join(root, "content"),
	})
	if err != nil {
		return nil, err
	}

	ret := &Store{
		content: content,
	}
	return ret, nil
}

// First pass at an API is below
// --------------------------------------------------

// Documents

func (s *Store) CreateDocument( /* TODO */ ) error {
	panic("unimplemented")
}

func (s *Store) AddFileToDocument( /* TODO */ ) error {
	panic("unimplemented")
}

func (s *Store) RemoveFileFromDocument( /* TODO */ ) error {
	panic("unimplemented")
}

func (s *Store) GetDocumentFiles( /* TODO */ ) error {
	panic("unimplemented")
}

func (s *Store) AddTagToDocument( /* TODO */ ) error {
	panic("unimplemented")
}

func (s *Store) RemoveTagFromDocument( /* TODO */ ) error {
	panic("unimplemented")
}

func (s *Store) FindDocumentById( /* TODO */ ) error {
	panic("unimplemented")
}

// Tags

func (s *Store) CreateTag( /* TODO */ ) error {
	panic("unimplemented")
}

func (s *Store) FindDocumentsByTag( /* TODO */ ) error {
	panic("unimplemented")
}
