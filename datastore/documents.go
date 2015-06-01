package datastore

import (
	"time"

	"github.com/andrew-d/castore"
	"github.com/jmoiron/sqlx"
	"github.com/lann/squirrel"

	"github.com/andrew-d/docstore/models"
	"github.com/andrew-d/docstore/services"
)

type documentsStore struct {
	db      *sqlx.DB
	content *castore.CAStore
}

func newDocumentsStore(db *sqlx.DB, path string) (*documentsStore, error) {
	content, err := castore.New(castore.Options{
		BasePath: path,
	})
	if err != nil {
		return nil, err
	}

	ret := &documentsStore{
		db:      db,
		content: content,
	}
	return ret, nil
}

func (s *documentsStore) Get(id int) (*models.Document, error) {
	var (
		documents []*models.Document
		stmt      *sqlx.NamedStmt
		err       error
	)

	if stmt, err = s.db.PrepareNamed(`SELECT * FROM documents WHERE id = :id;`); err != nil {
		return nil, err
	}
	if err = stmt.Select(&documents, M{"id": id}); err != nil {
		return nil, err
	}

	if len(documents) == 0 {
		// TODO: not found error?
		return nil, nil
	}

	return documents[0], nil
}

func (s *documentsStore) List(opts *services.DocumentsListOptions) ([]*models.Document, error) {
	if opts == nil {
		opts = &services.DocumentsListOptions{}
	}

	query := squirrel.Select("*").From("documents")
	// TODO: conditions go here

	query = query.
		OrderBy("created_at DESC").
		Limit(uint64(opts.PerPageOrDefault())).
		Offset(uint64(opts.Offset()))

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var documents []*models.Document
	if err = s.db.Select(&documents, sql, args...); err != nil {
		return nil, err
	}

	return documents, nil
}

func (s *documentsStore) Create(doc *models.Document) error {
	if doc == nil {
		panic("cannot create a nil document")
	}

	return dbretry(3, func() error {
		return transact(s.db, func(tx *sqlx.Tx) error {
			// Insert the document
			created_at := time.Now().UTC()
			res, err := tx.NamedExec(rid(s.db, `INSERT INTO documents VALUES (:desc, :created_at);`),
				map[string]interface{}{
					"desc":       doc.Description,
					"created_at": created_at.Unix(),
				})
			if err != nil {
				return err
			}

			document_id, err := res.LastInsertId()
			if err != nil {
				return err
			}

			// Check for and create each file
			newFiles := []models.File{}
			for _, file := range doc.Files {
				var (
					existing []*models.File
					stmt     *sqlx.NamedStmt
				)
				if stmt, err = s.db.PrepareNamed(`SELECT * FROM files WHERE content_key = :key;`); err != nil {
					return err
				}
				if err = stmt.Select(&existing, M{"key": file.ContentKey}); err != nil {
					return err
				}

				if len(existing) > 0 {
					return services.ErrFileExists
				}

				// Create a new file for this document.
				fres, err := tx.NamedExec(rid(s.db, `INSERT INTO files VALUES (:did, :name, :key);`),
					M{
						"did":  document_id,
						"name": file.Name,
						"key":  file.ContentKey,
					})
				if err != nil {
					return err
				}

				file_id, err := fres.LastInsertId()
				if err != nil {
					return err
				}

				newFiles = append(newFiles, models.File{
					ID:         file_id,
					DocumentID: document_id,
					Name:       file.Name,
					ContentKey: file.ContentKey,
				})
			}

			// TODO: tags

			// Update document and return.
			doc.ID = document_id
			doc.CreatedAt = created_at
			doc.Files = newFiles
			return nil
		})
	})
}

var _ services.DocumentsService = &documentsStore{}
