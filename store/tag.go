package store

import (
	"github.com/jinzhu/gorm"

	"github.com/andrew-d/docstore/store/models"
)

func (s *Store) CreateTag(name string) (models.Tag, error) {
	t := models.Tag{Name: name}
	if err := s.db.Create(&t).Error; err != nil {
		return models.Tag{}, err
	}

	return t, nil
}

func (s *Store) GetTagByName(name string) (models.Tag, bool, error) {
	var t models.Tag

	if err := s.db.Where("name = ?", name).First(&t).Error; err != nil {
		if err == gorm.RecordNotFound {
			return models.Tag{}, false, nil
		}

		return models.Tag{}, false, err
	}

	return t, true, nil
}

func (s *Store) GetTagById(id int64) (models.Tag, bool, error) {
	var t models.Tag

	if err := s.db.Where("id = ?", id).First(&t).Error; err != nil {
		if err == gorm.RecordNotFound {
			return models.Tag{}, false, nil
		}

		return models.Tag{}, false, err
	}

	return t, true, nil
}

func (s *Store) GetDocumentsByTag(id int64) ([]models.Document, error) {
	panic("Unimplemented")
}

func (s *Store) AddTagToDocument(tag_id, document_id int64) error {
	// TODO: take a *Document and update?
	panic("Unimplemented")
}
