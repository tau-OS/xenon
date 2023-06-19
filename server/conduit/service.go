package conduit

import (
	"context"
	"sync"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/handler"
	"github.com/samber/lo"
)

type DeviceInfoContextKey struct{}

type DeviceInfo struct {
	Name      string
	PublicKey string
}

type ConnectedDevice struct {
	*DeviceInfo
	RPCServer *jrpc2.Server
}

type ConduitService struct {
	connectedDevices []ConnectedDevice
	mutex            sync.RWMutex
}

func (c *ConduitService) ListConnectedDevices(ctx context.Context) []DeviceInfo {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return lo.Map(c.connectedDevices, func(device ConnectedDevice, index int) DeviceInfo {
		return *device.DeviceInfo
	})
}

type BroadcastMessageParams struct {
	Message []byte
}

type BroadcastMessageNotification struct {
	Message []byte
	Sender  DeviceInfo
}

func (c *ConduitService) BroadcastMessage(ctx context.Context, params BroadcastMessageParams) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	currentDevice := ctx.Value(DeviceInfoContextKey{}).(*DeviceInfo)

	recipients := lo.Filter(c.connectedDevices, func(device ConnectedDevice, index int) bool {
		return currentDevice.PublicKey != device.PublicKey
	})

	for _, recipient := range recipients {
		err := recipient.RPCServer.Notify(ctx, "ReceiveBroadcastMessage", BroadcastMessageNotification{
			Message: params.Message,
			// TODO: This simplifies things for now. We should probably refer to the device by its public key instead.
			Sender: *currentDevice,
		})
		if err != nil {
			// TODO: Do something with this error
			println(err.Error())
		}
	}

	return nil
}

func (c *ConduitService) NewRPCServer(name, publicKey string) *jrpc2.Server {
	deviceInfo := &DeviceInfo{
		Name:      name,
		PublicKey: publicKey,
	}

	server := jrpc2.NewServer(handler.Map{
		"ListConnectedDevices": handler.New(c.ListConnectedDevices),
		"BroadcastMessage":     handler.New(c.BroadcastMessage),
	}, &jrpc2.ServerOptions{
		AllowPush: true,
		NewContext: func() context.Context {
			return context.WithValue(context.Background(), DeviceInfoContextKey{}, deviceInfo)
		},
	})

	c.mutex.Lock()
	defer c.mutex.Unlock()

	// TODO: If device already exists, disconnect it or deny connection

	c.connectedDevices = append(c.connectedDevices, ConnectedDevice{
		DeviceInfo: deviceInfo,
		RPCServer:  server,
	})

	return server
}

func (c *ConduitService) RemoveDevice(publicKey string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.connectedDevices = lo.Filter(c.connectedDevices, func(device ConnectedDevice, index int) bool {
		return device.PublicKey != publicKey
	})
}
