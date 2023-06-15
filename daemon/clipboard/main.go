package clipboard

import (
	"context"
	"os"

	"github.com/charmbracelet/log"
	"golang.design/x/clipboard"
)

var l = log.NewWithOptions(os.Stderr, log.Options{
	ReportCaller: true,
	Prefix:       "Clipboard",
})

func Run() {
	ctx := context.Background()

	if err := clipboard.Init(); err != nil {
		l.Fatalf("Failed to initialize clipboard: %s", err.Error())
	}

	textClipboard := clipboard.Watch(ctx, clipboard.FmtText)
	imageClipboard := clipboard.Watch(ctx, clipboard.FmtImage)

	for {
		select {
		case text := <-textClipboard:
			l.Infof("Text clipboard changed: %s", string(text))
		case image := <-imageClipboard:
			l.Infof("Image clipboard changed, size: %d", len(image))
		}
	}
}
