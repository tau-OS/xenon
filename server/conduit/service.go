package conduit

import (
	"context"
	"sync"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/handler"
)

type ConduitService struct {
	connectedDevices []string
	mutex            sync.Mutex
}

func (c *ConduitService) ListConnectedDevices(ctx context.Context) error {
	return nil
}

func (c *ConduitService) BroadcastMessage(ctx context.Context, message string) error {
	return nil
}

func (c *ConduitService) NewServer() *jrpc2.Server {
	return jrpc2.NewServer(handler.Map{
		"ListConnectedDevices": handler.New(c.ListConnectedDevices),
		"BroadcastMessage":     handler.New(c.BroadcastMessage),
	}, &jrpc2.ServerOptions{
		AllowPush: true,
		NewContext: func() context.Context {
			return context.Background()
		},
	})
}
