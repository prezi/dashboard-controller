package main

import (
	"slave/slave"
	"network"
	"net/http"
)

func main() {
	ownPort, slaveName, masterURL, OS := slave.SetUp()
	go slave.Heartbeat(1, slaveName, ownPort, masterURL)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		slave.BrowserHandler(w, r, OS)
		})
	err := http.ListenAndServe(":" + ownPort, nil)
	network.ErrorHandler(err, "Error starting HTTP server: %v\n")
}
