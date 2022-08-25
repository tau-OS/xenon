package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/authhandler"
)

const clientID = "MKpzEzmCr8Mmov9Sz7OEE"

func generateState() (string, error) {
	buf := make([]byte, 16)

	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf), nil
}

func authHandler(authCodeURL string) (code string, state string, err error) {
	// Thanks to Xiao Yijun in the Logto Discord server
	authCodeURL += "&prompt=consent"
	type AuthResponse struct {
		Code  string
		State string
	}

	errorCh := make(chan error)
	codeCh := make(chan AuthResponse)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, authCodeURL, http.StatusFound)
	})

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code = r.URL.Query().Get("code")
		state = r.URL.Query().Get("state")

		w.WriteHeader(200)
		w.Write([]byte("You may close the window now."))

		codeCh <- AuthResponse{
			Code:  code,
			State: state,
		}
	})

	server := &http.Server{Addr: ":6969", Handler: mux}
	go func() {
		println("Listening on http://localhost:6969")
		if err := server.ListenAndServe(); err != nil {
			errorCh <- err
		}
	}()

	select {
	case err := <-errorCh:
		return "", "", err
	case authResponse := <-codeCh:
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			return "", "", err
		}

		return authResponse.Code, authResponse.State, nil
	}
}

func RunAuthenticationFlow(ctx context.Context) error {
	provider, err := oidc.NewProvider(ctx, "https://accounts.fyralabs.com/oidc")
	if err != nil {
		return err
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: clientID,
	})

	// Configure an OpenID Connect aware OAuth2 client.
	oauth2Config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: "",
		RedirectURL:  "http://localhost:6969/callback",
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{"offline_access", oidc.ScopeOpenID, "profile"},
	}

	state, err := generateState()
	if err != nil {
		return err
	}

	pkceParams, err := getNewPKCEParams()
	if err != nil {
		return err
	}

	tokenSource := authhandler.TokenSourceWithPKCE(ctx, oauth2Config, state, authHandler, pkceParams)

	token, err := tokenSource.Token()
	if err != nil {
		return err
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return errors.New("id_token was not a string")
	}

	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return err
	}

	var claims map[string]interface{}
	if err := idToken.Claims(&claims); err != nil {
		return err
	}

	return nil
}
