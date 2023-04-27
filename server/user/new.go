package user

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/tau-OS/xenon/server/database"
	"github.com/tau-OS/xenon/server/ent/user"
)

func isNew(id string) bool {
	ctx := context.Background()
	_, err := database.DatabaseClient.User.Query().Where(user.IDEQ(id)).Only(ctx)
	return err == nil
}

func ChkNew(c *fiber.Ctx) {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["id"].(string)
	if isNew(id) {
		initUser(claims)
	}
}

// Initialise a user record (they have an account but never used sync)
func initUser(claims jwt.MapClaims) error {
	dbctx := context.Background()
	if err := database.DatabaseClient.User.Create().SetID(claims["id"].(string)).Exec(dbctx); err != nil {
		return err
	}
	// ...
	return nil
}
