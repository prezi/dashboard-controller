package main

import (
	"master/master"
	"net/http"
	"log"
)

func main() {
	slaveMap := master.SetUp()
	http.HandleFunc("/register_webserver", func(w http.ResponseWriter, r *http.Request) {
			master.MonitorWebserverRegistration(w, r)
		})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		master.ReceiveRequestAndSendToSlave(w, r, slaveMap)
	})
	http.HandleFunc("/receive_heartbeat", func(w http.ResponseWriter, r *http.Request) {
		master.MonitorSlaveHeartbeats(w, r, slaveMap)
		})
	http.HandleFunc("/webserver_init", func(w http.ResponseWriter, r *http.Request) {
		master.WebserverRequestSlaveIds(w, r, slaveMap)
		})
	go master.MonitorSlaves(3, slaveMap)
	log.Fatal(http.ListenAndServe(":5000", nil))
}
