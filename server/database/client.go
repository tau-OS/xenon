package database

import (
	"context"

	"github.com/tau-OS/xenon/server/config"

	_ "github.com/lib/pq"
	"github.com/tau-OS/xenon/server/ent"
)

var DatabaseClient *ent.Client

func InitializeDatabase() error {
	client, err := ent.Open("postgres", config.Environment.DatabaseOptions)
	if err != nil {
		return err
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		panic(err.Error())
	}

	DatabaseClient = client

	return nil
}
