package main

import (
	"slave/slaveModules"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func main() {
	port, slaveName, masterIP, OS := slaveModule.SetUp()
	go slaveModule.Heartbeat(1, slaveName, masterIP)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		slaveModule.BrowserHandler(w, r, OS)
		})
	err := http.ListenAndServe(":" + strconv.Itoa(port), nil)
	if err != nil {
		fmt.Printf("Error starting HTTP server: %v\n", err)
		fmt.Println("ERROR: Failed to start HTTP server.")
		fmt.Println("Aborting program.")
		os.Exit(1)
	}
}
