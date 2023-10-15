package main

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/tau-OS/xenon/server/conduit"
	"github.com/tau-OS/xenon/server/config"
	"github.com/tau-OS/xenon/server/database"
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

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(c.Get("Authorization"))
	})

	app.Use("/api", jwtware.New(jwtware.Config{
		JWKSetURLs: []string{"https://auth.fyralabs.com/oidc/jwks"},
	}))

	// This serves a websocket connection providing a JSON-RPC 2.0 API to the user's personal "conduit service"
	// This service is specific to the user and is used as a way for connected devices to communicate with each other
	app.Get("/api/conduit", conduit.HandleConduitRequest, websocket.New(conduit.HandleWebSocketConnection))

	// app.Get("/api/ack", func(c *fiber.Ctx) error {
	// 	// for acknoledging a client running?
	// 	return c.SendStatus(200)
	// })

	// app.Get("/api/events", func(c *fiber.Ctx) error {
	// 	return c.SendStatus(200)
	// })

	// app.Get("/api/clipboard", func(c *fiber.Ctx) error {
	// 	return c.SendStatus(200)
	// })

	// app.Post("/api/clipboard", func(c *fiber.Ctx) error {
	// 	return c.SendStatus(200)
	// })

	if err := app.Listen(":8080"); err != nil {
		panic(err.Error())
	}
}
