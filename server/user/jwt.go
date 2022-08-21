package user

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

type UserClaims struct {
	jwt.RegisteredClaims
}

func (c UserClaims) Valid() error {
	if err := c.RegisteredClaims.Valid(); err != nil {
		return err
	}

	if !c.VerifyIssuer("https://accounts.fyralabs.com/oidc", true) {
		return errors.New("invalid issuer")
	}

	if !c.VerifyAudience("https://xenon.tauos.co/", true) {
		return errors.New("invalid audience")
	}

	return nil
}
