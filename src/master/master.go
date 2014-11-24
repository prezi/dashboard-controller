package main

import (
	"master/master"
	"net/http"
)

func main() {
	slaveMap := master.SetUp()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		master.ReceiveRequestAndSendToSlave(w, r, slaveMap)
	})
	http.HandleFunc("/receive_slave", func(w http.ResponseWriter, r *http.Request) {
		master.ReceiveAndMapSlaveAddress(w, r, slaveMap)
	})
	http.HandleFunc("/receive_heartbeat", func(w http.ResponseWriter, r *http.Request) {
		master.MonitorSlaveHeartbeats(w, r, slaveMap)
		})
	http.HandleFunc("/webserver-init", func(w http.ResponseWriter, r *http.Request) {
		master.WebserverRequestSlaveIds(w, r, slaveMap)
		})
	// http.HandleFunc("/webserver-init", master.WebserverRequestSlaveIds)
	go master.MonitorSlaves(3, slaveMap) 
	http.ListenAndServe(":5000", nil)
}
