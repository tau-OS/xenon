package conduit

import (
	"bytes"
	"context"
	"filippo.io/age"
	"github.com/creachadair/jrpc2"
	"github.com/samber/lo"
	"github.com/tau-OS/xenon/daemon/crypt"
	conduitServer "github.com/tau-OS/xenon/server/conduit"
	"io"
	"os"
)

func BroadcastMessage(ctx context.Context, service *jrpc2.Client) error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	devicesRes, err := service.Call(ctx, "ListConnectedDevices", nil)
	if err != nil {
		return err
	}

	var connectedDevices []conduitServer.DeviceInfo
	if err := devicesRes.UnmarshalResult(&connectedDevices); err != nil {
		return err
	}

	recipients := lo.Map(connectedDevices, func(device conduitServer.DeviceInfo, index int) age.Recipient {
		return lo.Must(age.ParseX25519Recipient(device.PublicKey))
	})

	var encrypted bytes.Buffer

	writer, err := crypt.Encrypt(&encrypted, recipients...)
	if err != nil {
		return err
	}

	if _, err := io.WriteString(writer, "Hello from "+hostname); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}

	_, err = service.Call(context.Background(), "BroadcastMessage", conduitServer.BroadcastMessageParams{
		Message: encrypted.Bytes(),
	})
	if err != nil {
		return err
	}

	return nil
}
