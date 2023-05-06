package main

import (
	"crypto/tls"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/logto-io/go/client"
)

var AuthSecret string // defined during compile time
const AuthId = "milfl97u099w4az99rp3l"

var authcfg = &client.LogtoConfig{
	Endpoint:  "https://auth.fyralabs.com",
	AppId:     AuthId,
	AppSecret: AuthSecret,
	Scopes:    []string{"openid", "profile", "offline_access"},
	Resources: []string{},
	Prompt:    "consent",
}

func signin(c *fiber.Ctx) bool {
	sess, err := store.Get(c)
	if err != nil {
		println("can't get store sess: " + err.Error())
		return true
	}
	logtoClient := client.NewLogtoClient(authcfg, &SessionStorage{session: sess})
	uri, err := logtoClient.SignIn("https://sync.fyralabs.com/sign-in-callback")
	if err != nil {
		println("can't sign in: " + err.Error())
		return true
	}
	c.Redirect(uri, http.StatusTemporaryRedirect)
	return false
}

type SessionStorage struct {
	session *session.Session
}

func authcallback(c *fiber.Ctx) *client.LogtoClient {
	sess, err := store.Get(c)
	if err != nil {
		println("can't get store sess: " + err.Error())
		return nil
	}
	session := &SessionStorage{session: sess}
	logtoClient := client.NewLogtoClient(authcfg, session)
	// convert to http.Request as required by HandleSignInCallback()
	req := c.Request()
	httpreq := http.Request{}
	if string(req.Header.Protocol()) == "https" {
		httpreq.TLS = &tls.ConnectionState{Version: 1} // anything but nil
		httpreq.Header.Add("X-Forwarded-Proto", c.Get("X-Forwarded-Proto"))
		httpreq.Host = string(req.Host())
		httpreq.RequestURI = string(req.RequestURI())
	}
	if e := logtoClient.HandleSignInCallback(&httpreq); e != nil {
		println("cannot handle sign in callback: " + e.Error())
		return nil
	}
	return logtoClient
}

func (storage *SessionStorage) GetItem(key string) string {
	value := storage.session.Get(key)
	if value == nil {
		return ""
	}
	return value.(string)
}

func (storage *SessionStorage) SetItem(key, value string) {
	storage.session.Set(key, value)
	storage.session.Save()
}
