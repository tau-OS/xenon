package conduit

import (
	"bytes"
	"context"
	"encoding/json"

	"filippo.io/age"
	"github.com/puzpuzpuz/xsync/v2"
	"github.com/samber/lo"
	"github.com/tau-OS/xenon/daemon/crypt"
	conduitServer "github.com/tau-OS/xenon/server/conduit"
)

var broadcastHandlers = xsync.NewMap()

type broadcastMessage struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

// This function dispatches the correct handler to handle the raw, decrypted broadcast message
func handleBroadcastMessage(message []byte) error {
	var data json.RawMessage
	wrapper := broadcastMessage{
		Payload: &data,
	}

	if err := json.Unmarshal(message, &wrapper); err != nil {
		return err
	}

	l.Debugf("received broadcast message of type %s", wrapper.Type)

	handler, found := broadcastHandlers.Load(wrapper.Type)
	if !found {
		l.Debugf("no handler found for broadcast message type %s", wrapper.Type)
		return nil
	}

	handler.(func(json.RawMessage))(data)

	return nil
}

func Broadcast(ctx context.Context, payloadType string, payload any) error {
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

	encoded, err := json.Marshal(broadcastMessage{
		Type:    payloadType,
		Payload: payload,
	})
	if err != nil {
		return err
	}

	if _, err := writer.Write(encoded); err != nil {
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

	l.Debugf("broadcasted message of type %s", payloadType)

	return nil
}

func RegisterBroadcastHandler[T any](payloadType string, handler func(T)) {
	broadcastHandlers.Store(payloadType, func(data json.RawMessage) {
		var payload T
		if err := json.Unmarshal(data, &payload); err != nil {
			l.Errorf("failed to unmarshal broadcast message payload for type %s: %s", payloadType, err)
			return
		}

		handler(payload)
	})
}
