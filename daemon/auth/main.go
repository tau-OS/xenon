package auth

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/logto-io/go/client"

	"github.com/tau-OS/xenon/daemon/storage"
)

var l = log.New(os.Stderr, "[auth] ", log.LstdFlags)

var appId = "xo0jronb7inwpqdf5ilf8"
var logtoConfig = &client.LogtoConfig{
	Endpoint:  "https://auth.fyralabs.com",
	AppId:     appId,
	Scopes:    []string{"openid", "profile", "offline_access"},
	Resources: []string{},
	Prompt:    "consent",
}
var logtoClient *client.LogtoClient

func initializeLogto() {
	logtoClient = client.NewLogtoClient(logtoConfig, storage.Keyring)
}

func startInteractiveAuth() {
	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    "localhost:9090",
		Handler: mux,
	}

	mux.Handle("/callback", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := logtoClient.HandleSignInCallback(r); err != nil {
			l.Fatalln("Failed to handle sign-in callback: " + err.Error())
		}

		if _, err := io.WriteString(w, "You are now signed in! You may close this window."); err != nil {
			l.Fatalln("Failed to write response: " + err.Error())
		}

		go func() {
			if err := srv.Shutdown(context.Background()); err != nil {
				l.Fatalln("Failed to shutdown server: " + err.Error())
			}
		}()
	}))

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		redirect, err := logtoClient.SignIn("http://localhost:9090/callback")
		if err != nil {
			l.Fatalln("Failed to generate sign-in link: " + err.Error())
		}

		http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
	}))

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		l.Fatalln("Failed to start server: " + err.Error())
	}
}

const prompt = `
┌─────────────────────────────────────────────┐
│ You are not signed in. On a browser, go to: │
│                                             │
│       ==>  http://localhost:9090  <==       │
│                                             │
│ Done? Paste the authentication token below. │
└─────────────────────────────────────────────┘`

func EnsureAuthenticated() {
	initializeLogto()

	if logtoClient.IsAuthenticated() {
		l.Println("Already authenticated, skipping interactive auth")
		return
	}

	l.Println(prompt)
	startInteractiveAuth()
}
