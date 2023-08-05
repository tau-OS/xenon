package conduit

import (
	"context"
	"net/http"
	"net/url"
	"os"

	"github.com/charmbracelet/log"
	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
	"github.com/tau-OS/xenon/daemon/auth"
	"github.com/tau-OS/xenon/daemon/crypt"

	"golang.org/x/net/websocket"
)

// This handles the client-side portion of the conduit service
// Here, we open and maintain the conduit connection, and expose a simple API for interacting with it

var l = log.NewWithOptions(os.Stderr, log.Options{
	ReportCaller: true,
	Prefix:       "Conduit",
	Level:        log.ParseLevel(os.Getenv("LOG_LEVEL")),
})

var service *jrpc2.Client

func Run() {
	hostname, err := os.Hostname()
	if err != nil {
		l.Fatal(err)
	}

	conduitURL, err := url.Parse("ws://192.168.122.1:8080/api/conduit")
	if err != nil {
		l.Fatal(err)
	}

	query := conduitURL.Query()
	query.Add("deviceName", hostname)
	query.Add("publicKey", crypt.PublicKey())
	conduitURL.RawQuery = query.Encode()

	originURL, err := url.Parse("http://192.168.122.1:8080")
	if err != nil {
		l.Fatal(err)
	}

	token, err := auth.LogtoClient.GetAccessToken("https://sync.fyralabs.com")
	if err != nil {
		l.Fatal(err)
	}

	ws, err := websocket.DialConfig(&websocket.Config{
		Location: conduitURL,
		Origin:   originURL,
		Version:  websocket.ProtocolVersionHybi,
		Header: http.Header{
			"Authorization": []string{"Bearer " + token.Token},
		},
	})
	if err != nil {
		l.Fatal(err)
	}

	service = jrpc2.NewClient(channel.RawJSON(ws, ws), &jrpc2.ClientOptions{
		OnNotify: handleNotify,
	})

	res, err := service.Call(context.Background(), "rpc.serverInfo", nil)
	if err != nil {
		l.Fatal(err)
	}

	l.Infof("Connected to conduit server: %s", res.ResultString())
}
