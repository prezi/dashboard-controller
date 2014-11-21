package main

import (
	"slave/slave"
	"network"
	"net/http"
)

func main() {
	port, slaveName, masterURL, OS := slave.SetUp()
	go slave.Heartbeat(1, slaveName, masterURL)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		slave.BrowserHandler(w, r, OS)
		})
	err := http.ListenAndServe(":" + port, nil)
	network.ErrorHandler(err, "Error starting HTTP server: %v\n")
}
