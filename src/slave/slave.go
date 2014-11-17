package main

import (
	"slave/slaveModules"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

var OS string

func main() {
	port, slaveName, masterIP := slaveModule.SetUp()
	OS = slaveModule.GetOS()
	slaveModule.Heartbeat(slaveName, masterIP)
	http.HandleFunc("/", handleRequest)

	// start HTTP server with given address and handler
	// handler=nil will default handler to DefaultServeMux
	err := http.ListenAndServe(":" + strconv.Itoa(port), nil)
	if err != nil {
		fmt.Printf("Error starting HTTP server: %v\n", err)
		fmt.Println("ERROR: Failed to start HTTP server.")
		fmt.Println("Aborting program.")
		os.Exit(1)
	}
}

func handleRequest(writer http.ResponseWriter, request *http.Request) {
	url := request.PostFormValue("url")
	slaveModule.KillBrowser(OS)
	slaveModule.OpenBrowser(OS, url)
	fmt.Fprintf(writer, "SUCCESS. \"%v\" has been posted.\n", url)
}
