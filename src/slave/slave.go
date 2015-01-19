package main

import (
	"fmt"
	"net/http"
	"network"
	"os"
	"slave/slave"
	"slave/slave/LinuxBrowserHandler"
	"slave/slave/OSXBrowserHandler"
)

func main() {
	ownPort, slaveName, masterURL, proxyURL, OS := slave.SetUp()
	go slave.Heartbeat(1, slaveName, ownPort, masterURL)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch OS {
		case "Linux":
			LinuxBrowserHandler.BrowserHandler(w, r, proxyURL)

		case "OS X":
			OSXBrowserHandler.BrowserHandler(w, r)
		}
	})
	http.HandleFunc("/receive_killsignal", func(_ http.ResponseWriter, request *http.Request) {
		if "die" == request.FormValue("message") {
			fmt.Println("Slave with this name already exists. Please restart slave with a different name.")
			os.Exit(1)
		}
	})
	err := http.ListenAndServe(":"+ownPort, nil)
	network.ErrorHandler(err, "Error starting HTTP server: %v\n")
}
