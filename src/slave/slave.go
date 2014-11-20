package main

import (
	"slave/slave"
	"network"
	"net/http"
	"strconv"
)

func main() {
	port, slaveName, masterURL, OS := slave.SetUp()
	go slave.Heartbeat(1, slaveName, masterURL)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		slave.BrowserHandler(w, r, OS)
		})
	err := http.ListenAndServe(":" + strconv.Itoa(port), nil)
	network.ErrorHandler(err, "Error starting HTTP server: %v\n")
}
