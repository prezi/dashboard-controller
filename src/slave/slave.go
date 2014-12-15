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
	ownPort, slaveName, masterURL, OS, browserProcess := slave.SetUp()
	go slave.Heartbeat(1, slaveName, ownPort, masterURL)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch OS {
			case "Linux":
				browserProcess = LinuxBrowserHandler.BrowserHandler(w, r, browserProcess)

			case "OS X":
				browserProcess = OSXBrowserHandler.BrowserHandler(w, r, browserProcess)
			}
//		browserProcess = slave.BrowserHandler(w, r, OS, browserProcess)
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
