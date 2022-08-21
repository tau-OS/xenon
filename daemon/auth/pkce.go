package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"math/big"

	"golang.org/x/oauth2/authhandler"
)

var pkceVerifierChars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~")

func getPKCEVerifier() (string, error) {
	buf := make([]rune, 64)

	for i := range buf {
		randInt, err := rand.Int(rand.Reader, big.NewInt(int64(len(pkceVerifierChars))))
		if err != nil {
			return "", err
		}

		buf[i] = pkceVerifierChars[randInt.Int64()]
	}

	return string(buf), nil
}

func getPKCEChallenge(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func getNewPKCEParams() (*authhandler.PKCEParams, error) {
	verifier, err := getPKCEVerifier()
	if err != nil {
		return nil, err
	}

	challenge := getPKCEChallenge(verifier)
	return &authhandler.PKCEParams{
		Verifier:        verifier,
		Challenge:       challenge,
		ChallengeMethod: "S256",
	}, nil
}
