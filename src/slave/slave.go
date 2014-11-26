package main

import (
	"slave/slave"
	"network"
	"net/http"
	"os"
	"fmt"
)

func main() {
	ownPort, slaveName, masterURL, OS := slave.SetUp()
	go slave.Heartbeat(1, slaveName, ownPort, masterURL)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		slave.BrowserHandler(w, r, OS)
		})
	http.HandleFunc("/receive_killsignal", func(_ http.ResponseWriter, request *http.Request) {
			if "die" == request.FormValue("message") {
				fmt.Println("Master refused slave. Please restart slave with a different name.")
				os.Exit(1)
			}
		})

	http.HandleFunc("/new_name", func(w http.ResponseWriter, r *http.Request) {
			slave.BrowserHandler(w, r, OS)
		})
	err := http.ListenAndServe(":" + ownPort, nil)
	network.ErrorHandler(err, "Error starting HTTP server: %v\n")
}
