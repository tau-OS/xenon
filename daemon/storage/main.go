package storage

import (
	"errors"
	"log"
	"os"

	"git.mills.io/prologic/bitcask"
)

var l = log.New(os.Stderr, "[ db ] ", log.LstdFlags)

type LocalStorage struct {
	*bitcask.Bitcask
}

func (s *LocalStorage) GetItem(key string) string {
	data, err := s.Get([]byte(key))

	if errors.Is(err, bitcask.ErrKeyNotFound) {
		return ""
	}

	if err != nil {
		l.Fatalf("FAIL to get item `%s`: %s", key, err)
	}

	return string(data)
}

func (s *LocalStorage) SetItem(key, value string) {
	if err := s.Put([]byte(key), []byte(value)); err != nil {
		l.Fatalf("FAIL to set item `%s` -> `%s`: %s", key, value, err)
	}
}

var Local *LocalStorage
var Keyring *KeyringStorage

func InitLocalStorage() error {
	l.Println("Initializing local storage...")
	home, err := os.UserHomeDir()
	if err != nil {
		l.Fatal("FAIL to get home dir: " + err.Error())
	}

	db, err := bitcask.Open(home + "/.config/tausync/db")
	if err != nil {
		return err
	}

	Local = &LocalStorage{db}

	return nil
}
