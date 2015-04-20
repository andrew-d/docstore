package store

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"

	"github.com/andrew-d/docstore/store/models"
)

var (
	ErrExist = errors.New("store: path exists")
)

type Store struct {
	path string
	db   gorm.DB
}

func New(path string) (*Store, error) {
	_, err := os.Stat(path)
	if err == nil {
		return nil, ErrExist
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	// Make the directory
	err = os.Mkdir(path, 0700)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open("sqlite3", filepath.Join(path, "docstore.sqlite3"))
	if err != nil {
		return nil, err
	}

	// Create models in the database.
	db.AutoMigrate(&models.Tag{}, &models.File{}, &models.Document{})

	ret := &Store{
		path: path,
		db:   db,
	}
	return ret, nil
}

func Open(path string) (*Store, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open("sqlite3", filepath.Join(path, "docstore.sqlite3"))
	if err != nil {
		return nil, err
	}

	ret := &Store{
		path: path,
		db:   db,
	}
	return ret, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}
