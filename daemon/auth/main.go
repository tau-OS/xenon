package auth

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/logto-io/go/client"

	"github.com/tau-OS/xenon/daemon/storage"
)

var l = log.NewWithOptions(os.Stderr, log.Options{
	ReportCaller: true,
	Prefix:       "Auth",
	Level:        log.ParseLevel(os.Getenv("LOG_LEVEL")),
})

var appId = "xo0jronb7inwpqdf5ilf8"
var logtoConfig = &client.LogtoConfig{
	Endpoint: "https://auth.fyralabs.com",
	AppId:    appId,
	Scopes:   []string{"openid", "profile", "offline_access"},
	Resources: []string{
		"https://sync.fyralabs.com",
	},
	Prompt: "consent",
}
var LogtoClient *client.LogtoClient

func initializeLogto() {
	LogtoClient = client.NewLogtoClient(logtoConfig, storage.Keyring)
}

func startInteractiveAuth() {
	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    "localhost:9090",
		Handler: mux,
	}

	mux.Handle("/callback", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := LogtoClient.HandleSignInCallback(r); err != nil {
			l.Fatal("Failed to handle sign-in callback: " + err.Error())
		}

		if _, err := io.WriteString(w, "You are now signed in! You may close this window."); err != nil {
			l.Fatal("Failed to write response: " + err.Error())
		}

		go func() {
			if err := srv.Shutdown(context.Background()); err != nil {
				l.Fatal("Failed to shutdown server: " + err.Error())
			}
		}()
	}))

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		redirect, err := LogtoClient.SignIn("http://localhost:9090/callback")
		if err != nil {
			l.Fatal("Failed to generate sign-in link: " + err.Error())
		}

		http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
	}))

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		l.Fatal("Failed to start server: " + err.Error())
	}
}

const prompt = `
┌─────────────────────────────────────────────┐
│ You are not signed in. On a browser, go to: │
│                                             │
│       ==>  http://localhost:9090  <==       │
│                                             │
└─────────────────────────────────────────────┘`

func EnsureAuthenticated() {
	initializeLogto()

	if LogtoClient.IsAuthenticated() {
		l.Info("Already authenticated, skipping interactive auth")
		return
	}

	println(prompt)
	startInteractiveAuth()
}
