package main

import (
	"context"

	"github.com/tau-OS/xenon/daemon/auth"
)

func main() {
	if err := auth.RunAuthenticationFlow(context.Background()); err != nil {
		panic(err.Error())
	}
}
