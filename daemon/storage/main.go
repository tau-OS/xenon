package storage

import (
	"errors"

	"git.mills.io/prologic/bitcask"
)

type LocalStorage struct {
	*bitcask.Bitcask
}

func (s *LocalStorage) GetItem(key string) string {
	data, err := s.Get([]byte(key))

	if errors.Is(err, bitcask.ErrKeyNotFound) {
		return ""
	}

	if err != nil {
		panic(err.Error())
	}

	return string(data)
}

func (s *LocalStorage) SetItem(key, value string) {
	if err := s.Put([]byte(key), []byte(value)); err != nil {
		panic(err.Error())
	}
}

var Local *LocalStorage

func InitLocalStorage() error {
	db, err := bitcask.Open("/tmp/db")
	if err != nil {
		return err
	}

	Local = &LocalStorage{db}

	return nil
}
