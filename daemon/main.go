package main

import (
	"github.com/tau-OS/xenon/daemon/auth"
	"github.com/tau-OS/xenon/daemon/storage"
)

func main() {
	if err := storage.InitLocalStorage(); err != nil {
		panic(err.Error())
	}

	defer storage.Local.Close()

	if err := auth.LogIn(); err != nil {
		panic(err.Error())
	}
}
