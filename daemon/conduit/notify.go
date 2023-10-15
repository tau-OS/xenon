package conduit

import (
	"bytes"
	"io"

	"github.com/creachadair/jrpc2"
	"github.com/tau-OS/xenon/daemon/crypt"
	conduitServer "github.com/tau-OS/xenon/server/conduit"
)

func handleNotify(req *jrpc2.Request) {
	switch req.Method() {
	case "ReceiveBroadcastMessage":
		var notification conduitServer.BroadcastMessageNotification
		if err := req.UnmarshalParams(&notification); err != nil {
			l.Fatal(err)
		}

		reader, err := crypt.Decrypt(bytes.NewReader(notification.Message))
		if err != nil {
			l.Errorf("failed to decrypt broadcast message from %s (%s): %s", notification.Sender.PublicKey, notification.Sender.Name, err)
			return
		}

		bytes, err := io.ReadAll(reader)
		if err != nil {
			l.Errorf("failed to read decrypted broadcast message from %s (%s): %s", notification.Sender.PublicKey, notification.Sender.Name, err)
			return
		}

		if err := handleBroadcastMessage(bytes); err != nil {
			l.Errorf("failed to handle broadcast message from %s (%s): %s", notification.Sender.PublicKey, notification.Sender.Name, err)
			return
		}
	}
}
