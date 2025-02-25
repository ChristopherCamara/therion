package main

import (
	"fmt"
	"os"

	webview "github.com/webview/webview_go"
)

func main() {
	devMode := false
	reload := make(chan bool)
	args := os.Args[1:]
	w := webview.New(false)
	defer w.Destroy()
	for _, arg := range args {
		if arg == "--dev" {
			devMode = true
		}
	}
	server := newServer(devMode)
	go server.start()
	if devMode {
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
