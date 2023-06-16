package main

import (
	"io/ioutil"
	"net/http"

	"github.com/tau-OS/xenon/daemon/auth"
	"github.com/tau-OS/xenon/daemon/crypt"
	"github.com/tau-OS/xenon/daemon/storage"
)

func main() {
	if err := storage.InitLocalStorage(); err != nil {
		panic(err.Error())
	}

	defer storage.Local.Close()

	auth.EnsureAuthenticated()

	crypt.InitializeMachineIdentity()

	// clipboard.Run()

	accessToken, err := auth.LogtoClient.GetAccessToken("https://sync.fyralabs.com")
	if err != nil {
		panic(err.Error())
	}

	req, err := http.NewRequest("GET", "http://localhost:8080/", nil)
	if err != nil {
		panic(err.Error())
	}
	req.Header.Set("Authorization", "Bearer "+accessToken.Token)

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}

	println(res.StatusCode)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	println(string(body))

	// ioutil.readAll(res.Body)
}
