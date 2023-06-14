package main

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/icza/gog"
	"github.com/tau-OS/xenon/server/config"
	"github.com/tau-OS/xenon/server/database"
	"github.com/tau-OS/xenon/server/user"
)

var store = session.New()

func main() {
	if AuthSecret == "" {
		panic("Oops! `AuthSecret` is not set! Rebuild using this: `go build -ldflags '-X auth.AuthSecret=...' server`")
	}
	if err := config.InitializeEnv(); err != nil {
		panic(err.Error())
	}

	if err := database.InitializeDatabase(); err != nil {
		panic(err.Error())
	}

	app := fiber.New(fiber.Config{
		AppName: "Xenon",
	})

	app.Use("/api", jwtware.New(jwtware.Config{
		KeySetURL:            "https://logto.fyralabs.com/oidc/auth", // or `/oidc/token`?
		Claims:               &user.UserClaims{},
		KeyRefreshUnknownKID: gog.Ptr(false),
	}))

	app.Get("/sign-in", func(c *fiber.Ctx) error {
		if signin(c) {
			return c.SendStatus(http.StatusInternalServerError)
		}
		return nil
	})

	app.Get("/sign-in-callback", func(c *fiber.Ctx) error {
		logtoClient := authcallback(c)
		if logtoClient == nil {
			return c.SendStatus(500)
		}
		user.ChkNew(c)
		token, err := logtoClient.GetAccessToken("")
		if err != nil {
			println("! Can't get access token: " + err.Error())
			return c.SendStatus(500)
		}
		out := "Paste the following into the application: <br>\n"
		out += token.Token + "." + token.Scope + "." + string(rune(token.ExpiresAt))
		return c.Send([]byte(out))
	})

	app.Get("/api/ack", func(c *fiber.Ctx) error {
		// for acknoledging a client running?
		return c.SendStatus(200)
	})

	app.Get("/api/events", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Get("/api/clipboard", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Post("/api/clipboard", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	if err := app.Listen(":8080"); err != nil {
		panic(err.Error())
	}
}
