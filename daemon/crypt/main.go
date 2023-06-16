package crypt

import (
	"os"

	"filippo.io/age"
	"github.com/charmbracelet/log"
	"github.com/tau-OS/xenon/daemon/storage"
)

var l = log.NewWithOptions(os.Stderr, log.Options{
	ReportCaller: true,
	Prefix:       "Crypt",
})

var machineIdentity *age.X25519Identity

func PublicKey() string {
	return machineIdentity.Recipient().String()
}

func InitializeMachineIdentity() {
	l.Info("Initializing machine identity...")

	key, err := storage.Keyring.Get("machine_key")
	if err != nil {
		l.Fatal("Failed to get machine key: " + err.Error())
	}

	if key != nil {
		identity, err := age.ParseX25519Identity(*key)
		if err != nil {
			l.Fatal("Failed to parse machine key: " + err.Error())
		}

		machineIdentity = identity
		l.Info("Machine key already exists, initialized from keyring")

		return
	}

	identity, err := age.GenerateX25519Identity()
	if err != nil {
		l.Fatal("Failed to generate machine key: " + err.Error())
	}

	err = storage.Keyring.Set("machine_key", identity.String())
	if err != nil {
		l.Fatal("Failed to save machine key: " + err.Error())
	}

	machineIdentity = identity

	l.Info("Machine identity initialized")
}
