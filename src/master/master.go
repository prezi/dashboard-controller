package main

import (
	"log"
	"master/master"
	"net/http"
)

func main() {
	slaveMap := master.SetUp()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		master.ReceiveRequestAndSendToSlave(w, r, slaveMap)
	})
	http.HandleFunc("/receive_heartbeat", func(_ http.ResponseWriter, r *http.Request) {
		master.ReceiveSlaveHeartbeat(r, slaveMap)
	})
	http.HandleFunc("/webserver_heartbeat", func(_ http.ResponseWriter, r *http.Request) {
		master.UpdateWebserverAddress(r)
	})
	http.HandleFunc("/webserver_init", func(_ http.ResponseWriter, r *http.Request) {

		master.SendWebserverInit(r, slaveMap)
	})
	go master.MonitorSlaves(3, slaveMap)
	log.Fatal(http.ListenAndServe(":5000", nil))
}
