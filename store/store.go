package store

import (
	"errors"
	"os"

	"github.com/jinzhu/gorm"

	"github.com/andrew-d/docstore/store/models"
)

var (
	ErrExist = errors.New("store: path exists")
)

type Store struct {
	path string
	db   *gorm.DB
}

func New(path string) (*Store, error) {
	_, err := os.Stat(path)
	if err == nil {
		return nil, ErrExist
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	// Make the directory
	err = os.Mkdir(path, 0600)
	if err != nil {
		return nil, err
	}

	// TODO: create things in dir
	_ = models.Tag{}
	ret := &Store{
		path: path,
	}
	return ret, nil
}

func Open(path string) (*Store, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	// TODO: open things in dir
	ret := &Store{
		path: path,
	}
	return ret, nil
}
