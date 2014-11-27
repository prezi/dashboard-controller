package main

import (
	"fmt"
	"net/http"
	"network"
	"os"
	"slave/slave"
)

func main() {
	ownPort, slaveName, masterURL, OS, BrowserProcess := slave.SetUp()
	go slave.Heartbeat(1, slaveName, ownPort, masterURL)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		BrowserProcess = slave.BrowserHandler(w, r, OS, BrowserProcess)
		})
	http.HandleFunc("/receive_killsignal", func(_ http.ResponseWriter, request *http.Request) {
			if "die" == request.FormValue("message") {
				fmt.Println("Master refused slave. Please restart slave with a different name.")
				os.Exit(1)
			}
		})
	err := http.ListenAndServe(":" + ownPort, nil)
	network.ErrorHandler(err, "Error starting HTTP server: %v\n")
}
