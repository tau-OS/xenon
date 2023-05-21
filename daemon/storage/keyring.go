package storage

import (
	"errors"
	"log"

	"github.com/zalando/go-keyring"
)

const SERVICE = "TauSync"

func GetKey(name string) string {
	out, err := keyring.Get(SERVICE, name)
	if errors.Is(err, keyring.ErrNotFound) {
		return ""
	}
	if err != nil {
		log.Fatalf("[keys] FAIL to get key `%s`: %s\n", name, err.Error())
	}
	return out
}

func SetKey(name, val string) {
	err := keyring.Set(SERVICE, name, val)
	if err != nil {
		log.Fatalf("[keys] FAIL to set key `%s` to `%s`: %s", name, val, err.Error())
	}
}
