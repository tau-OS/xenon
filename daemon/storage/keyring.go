package storage

import (
	"errors"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/zalando/go-keyring"
)

var kl = l.With(os.Stderr, log.Options{
	Prefix: "Keyring",
})

const SERVICE = "TauSync"

type KeyringStorage struct{}

func (s *KeyringStorage) Get(name string) (*string, error) {
	out, err := keyring.Get(SERVICE, name)
	if errors.Is(err, keyring.ErrNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get key `%s`: %w", name, err)
	}

	return &out, nil
}

func (s *KeyringStorage) Set(name, val string) error {
	err := keyring.Set(SERVICE, name, val)
	if err != nil {
		return fmt.Errorf("failed to set key `%s` to `%s`: %w", name, val, err)
	}

	return nil
}

func (s *KeyringStorage) Delete(name string) error {
	err := keyring.Delete(SERVICE, name)
	if err != nil {
		return fmt.Errorf("failed to delete key `%s`: %w", name, err)
	}

	return nil
}

// Used by Logto's client for storage

func (s *KeyringStorage) GetItem(name string) string {
	out, err := keyring.Get(SERVICE, name)
	if errors.Is(err, keyring.ErrNotFound) {
		return ""
	}
	if err != nil {
		kl.Fatalf("Failed to get key `%s`: %s\n", name, err.Error())
	}
	return out
}

func (s *KeyringStorage) SetItem(name, val string) {
	err := keyring.Set(SERVICE, name, val)
	if err != nil {
		kl.Fatalf("Failed to set key `%s` to `%s`: %s", name, val, err.Error())
	}
}
