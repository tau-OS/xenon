package auth

import (
	"net/http"

	"github.com/logto-io/go/client"
	"github.com/tau-OS/xenon/daemon/storage"
)

var Client *client.LogtoClient = client.NewLogtoClient(&client.LogtoConfig{
	Endpoint:           "https://logto.fyralabs.com",
	AppId:              "MKpzEzmCr8Mmov9Sz7OEE",
	AppSecret:          "",
	Scopes:             []string{"openid", "profile", "offline_access"},
	Resources:          []string{},
	Prompt:             "consent",
}, storage.Local)


func LogIn() error {
	url, err := Client.SignIn("http://localhost:6969/callback")
	if err != nil {
		return err
	}

	errorCh := make(chan error)
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, url, http.StatusFound)
	})

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		if err := Client.HandleSignInCallback(r); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error while signing in"))
			errorCh <- err
		}

		errorCh <- nil
	})

	server := &http.Server{Addr: ":6969", Handler: mux}
	go func() {
		println("Listening on http://localhost:6969")
		if err := server.ListenAndServe(); err != nil {
			errorCh <- err
		}
	}()

	return <-errorCh
}
