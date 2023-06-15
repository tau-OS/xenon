package storage

import (
	"os"

	"git.mills.io/prologic/bitcask"
	"github.com/charmbracelet/log"
)

var l = log.NewWithOptions(os.Stderr, log.Options{
	ReportCaller: true,
	Prefix:       "Storage",
})

type LocalStorage struct {
	*bitcask.Bitcask
}

var Local *LocalStorage
var Keyring *KeyringStorage

func InitLocalStorage() error {
	l.Info("Initializing local storage...")
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
