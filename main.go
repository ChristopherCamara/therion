package main

import (
	"fmt"
	"os"

	webview "github.com/webview/webview_go"
)

func main() {
	devMode := false
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "--dev" {
			devMode = true
		}
	}

	w := webview.New(devMode)
	defer w.Destroy()

	server := newServer(devMode)
	go server.start()
	if devMode {
		reload := make(chan bool)
		go startDevMode(reload)
		go func() {
			for range reload {
				server.loadTemplates()
				w.Navigate(fmt.Sprintf("http://localhost:%d", SERVER_PORT))
			}
		}()
	}

	w.SetTitle("Therion")
	w.SetSize(1710, 1107, webview.HintNone)
	w.Navigate(fmt.Sprintf("http://localhost:%d", SERVER_PORT))
	w.Run()
}
