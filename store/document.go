package store

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/andrew-d/docstore/store/models"
)

func (s *Store) CreateDocument(name string) (models.Document, error) {
	d := models.Document{
		Name:      name,
		CreatedAt: time.Now().UTC().UnixNano(),
	}
	if err := s.db.Create(&d).Error; err != nil {
		return models.Document{}, err
	}

	return d, nil
}

func (s *Store) GetDocumentByName(name string) (models.Document, bool, error) {
	var d models.Document

	if err := s.db.Where("name = ?", name).First(&d).Error; err != nil {
		if err == gorm.RecordNotFound {
			return models.Document{}, false, nil
		}

		return models.Document{}, false, err
	}

	return d, true, nil
}

func (s *Store) GetDocumentById(id int64) (models.Document, bool, error) {
	var d models.Document

	if err := s.db.Where("id = ?", id).First(&d).Error; err != nil {
		if err == gorm.RecordNotFound {
			return models.Document{}, false, nil
		}

		return models.Document{}, false, err
	}

	return d, true, nil
}

func (s *Store) GetTagsByDocument(id int64) ([]models.Tag, error) {
	panic("Unimplemented")
}

func (s *Store) AddDocumentToTag(document_id, tag_id int64) error {
	// TODO: take a *Document and update?
	panic("Unimplemented")
}
