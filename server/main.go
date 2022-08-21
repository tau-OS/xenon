package main

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/icza/gog"
	"github.com/tau-OS/xenon/server/config"
	"github.com/tau-OS/xenon/server/database"
	"github.com/tau-OS/xenon/server/user"
)

func main() {
	if err := config.InitializeEnv(); err != nil {
		panic(err.Error())
	}

	if err := database.InitializeDatabase(); err != nil {
		panic(err.Error())
	}

	app := fiber.New(fiber.Config{
		AppName: "Xenon",
	})

	app.Use(jwtware.New(jwtware.Config{
		KeySetURL:            "https://accounts.fyralabs.com/oidc/jwks",
		Claims:               &user.UserClaims{},
		KeyRefreshUnknownKID: gog.Ptr(false),
	}))

	app.Get("/events", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Get("/clipboard", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Post("/clipboard", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	if err := app.Listen(":8080"); err != nil {
		panic(err.Error())
	}
}
