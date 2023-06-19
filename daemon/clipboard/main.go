package clipboard

import (
	"context"
	"os"

	"github.com/charmbracelet/log"
	"github.com/tau-OS/xenon/daemon/conduit"
	"golang.design/x/clipboard"
)

var l = log.NewWithOptions(os.Stderr, log.Options{
	ReportCaller: true,
	Prefix:       "Clipboard",
	Level:        log.DebugLevel,
})

type clipboardPayload struct {
	Type clipboard.Format
	Data []byte
}

func Run() {
	ctx := context.Background()

	if err := clipboard.Init(); err != nil {
		l.Fatalf("Failed to initialize clipboard: %s", err.Error())
	}

	textClipboard := clipboard.Watch(ctx, clipboard.FmtText)
	imageClipboard := clipboard.Watch(ctx, clipboard.FmtImage)

	conduit.RegisterBroadcastHandler("clipboard", func(payload clipboardPayload) {
		l.Debugf("Received clipboard payload of type %d with size %d", payload.Type, len(payload.Data))
		clipboard.Write(payload.Type, payload.Data)
	})

	for {
		select {
		case text := <-textClipboard:
			l.Debugf("Text clipboard changed: %s", string(text))
			conduit.Broadcast(ctx, "clipboard", clipboardPayload{
				Type: clipboard.FmtText,
				Data: text,
			})
		case image := <-imageClipboard:
			l.Debugf("Image clipboard changed, size: %d", len(image))
			conduit.Broadcast(ctx, "clipboard", clipboardPayload{
				Type: clipboard.FmtImage,
				Data: image,
			})
		}
	}
}
