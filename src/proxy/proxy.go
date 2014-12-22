package main

import (
	"net/http"
	"network"
	"proxy/proxy"
)

// run the proxy in a separate terminal window
func main() {
	masterURL := proxy.Start()
	go proxy.Heartbeat(1, masterURL)
	// TODO: Is mitmproxy still alive even after the Go program ends? Do we need to have a function here to keep the Go process alive?
	err := http.ListenAndServe(":6980", nil)
	network.ErrorHandler(err, "Error starting HTTP server: %v\n")
}
