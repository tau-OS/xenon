package conduit

import (
	"context"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/handler"
)

type ConduitService struct {
}

func (c *ConduitService) Identify(ctx context.Context, name string, publicKey string) {

}

func (c *ConduitService) ListConnectedDevices(ctx context.Context) {

}

func (c *ConduitService) BroadcastMessage(ctx context.Context, message string) {

}

func (c *ConduitService) NewServer() *jrpc2.Server {
	return jrpc2.NewServer(handler.Map{
		"Identify":             handler.New(c.Identify),
		"ListConnectedDevices": handler.New(c.ListConnectedDevices),
		"BroadcastMessage":     handler.New(c.BroadcastMessage),
	}, &jrpc2.ServerOptions{
		NewContext: func() context.Context {
			return context.Background()
		},
	})
}
