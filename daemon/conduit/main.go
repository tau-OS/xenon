package conduit

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"filippo.io/age"
	"github.com/charmbracelet/log"
	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
	"github.com/samber/lo"
	"github.com/tau-OS/xenon/daemon/auth"
	"github.com/tau-OS/xenon/daemon/crypt"
	conduitServer "github.com/tau-OS/xenon/server/conduit"

	"golang.org/x/net/websocket"
)

// This handles the client-side portion of the conduit service
// Here, we open and maintain the conduit connection, and expose a simple API for interacting with it

var l = log.NewWithOptions(os.Stderr, log.Options{
	ReportCaller: true,
	Prefix:       "Conduit",
})

func Run() {
	hostname, err := os.Hostname()
	if err != nil {
		l.Fatal(err)
	}

	conduitURL, err := url.Parse("ws://192.168.64.1:8080/api/conduit")
	if err != nil {
		l.Fatal(err)
	}

	query := conduitURL.Query()
	query.Add("deviceName", hostname)
	query.Add("publicKey", crypt.PublicKey())
	conduitURL.RawQuery = query.Encode()

	originURL, err := url.Parse("http://192.168.64.1:8080")
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

	service := jrpc2.NewClient(channel.RawJSON(ws, ws), &jrpc2.ClientOptions{
		OnNotify: func(req *jrpc2.Request) {
			switch req.Method() {
			case "ReceiveBroadcastMessage":
				var notification conduitServer.BroadcastMessageNotification
				if err := req.UnmarshalParams(&notification); err != nil {
					l.Fatal(err)
				}

				reader, err := crypt.Decrypt(bytes.NewBuffer([]byte(notification.Message)))
				if err != nil {
					l.Errorf("failed to decrypt broadcast message from %s (%s): %s", notification.Sender.PublicKey, notification.Sender.Name, err)
					return
				}

				bytes, err := io.ReadAll(reader)
				if err != nil {
					l.Errorf("failed to read decrypted broadcast message from %s (%s): %s", notification.Sender.PublicKey, notification.Sender.Name, err)
					return
				}

				l.Info("owo: %s", string(bytes))
			}
		},
	})

	res, err := service.Call(context.Background(), "rpc.serverInfo", nil)
	if err != nil {
		l.Fatal(err)
	}

	l.Infof("Connected to conduit server: %s", res.ResultString())

	for {
		devicesRes, err := service.Call(context.Background(), "ListConnectedDevices", nil)
		if err != nil {
			l.Fatal(err)
		}

		var connectedDevices []conduitServer.DeviceInfo
		if err := devicesRes.UnmarshalResult(&connectedDevices); err != nil {
			l.Fatal(err)
		}

		recipients := lo.Map(connectedDevices, func(device conduitServer.DeviceInfo, index int) age.Recipient {
			return lo.Must(age.ParseX25519Recipient(device.PublicKey))
		})

		encrypted := &bytes.Buffer{}

		writer, err := crypt.Encrypt(encrypted, recipients...)
		if err != nil {
			l.Fatal(err)
		}

		if _, err := io.WriteString(writer, "Hello from "+hostname); err != nil {
			l.Fatal(err)
		}
		if err := writer.Close(); err != nil {
			l.Fatal(err)
		}

		lo.Must(io.ReadAll(lo.Must(crypt.Decrypt(bytes.NewBuffer([]byte(encrypted.String()))))))

		_, err = service.Call(context.Background(), "BroadcastMessage", conduitServer.BroadcastMessageParams{
			Message: encrypted.String(),
		})
		if err != nil {
			l.Fatal(err)
		}
		time.Sleep(5 * time.Second)
	}
}
