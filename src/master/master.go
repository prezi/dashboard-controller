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
	http.HandleFunc("/receive_heartbeat", func(w http.ResponseWriter, r *http.Request) {
		master.MonitorSlaveHeartbeats(w, r, slaveMap)
	})
	http.HandleFunc("/webserver_heartbeat", func(w http.ResponseWriter, r *http.Request) {
			master.UpdateWebserverAddress(w, r)
		})
	http.HandleFunc("/webserver_init", func(w http.ResponseWriter, r *http.Request) {

		master.SendWebserverInit(w, r, slaveMap)
	})
	go master.MonitorSlaves(3, slaveMap)
	log.Fatal(http.ListenAndServe(":5000", nil))
}
