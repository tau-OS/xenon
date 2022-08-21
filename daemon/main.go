package main

import (
	"context"

	"golang.design/x/clipboard"
)

func main() {
	channel := clipboard.Watch(context.Background(), clipboard.FmtText)

	for {
		data := <-channel

		println(string(data))
	}
}
