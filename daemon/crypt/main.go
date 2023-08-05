package crypt

import (
	"io"
	"os"

	"filippo.io/age"
	"github.com/charmbracelet/log"
	"github.com/tau-OS/xenon/daemon/storage"
)

var l = log.NewWithOptions(os.Stderr, log.Options{
	ReportCaller: true,
	Prefix:       "Crypt",
	Level:        log.ParseLevel(os.Getenv("LOG_LEVEL")),
})

var machineIdentity *age.X25519Identity

func PublicKey() string {
	return machineIdentity.Recipient().String()
}

func Decrypt(data io.Reader) (io.Reader, error) {
	// TODO: We should check the sender's public key against the list of known devices
	return age.Decrypt(data, machineIdentity)
}

func Encrypt(destination io.Writer, recipients ...age.Recipient) (io.WriteCloser, error) {
	// TODO: Verify that our recipients is a trusted device
	return age.Encrypt(destination, recipients...)
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
