package main

import (
	"log"
	"master/master"
	"master/master/receiveAndSendRequestToSlave"
	"master/master/slaveMonitor"
	"net/http"
)

func main() {
	slaveMap := master.SetUp()

	http.HandleFunc("/receive_heartbeat", func(_ http.ResponseWriter, r *http.Request) {
		slaveMap = slaveMonitor.ReceiveSlaveHeartbeat(r, slaveMap)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		receiveAndSendRequestToSlave.ReceiveRequestAndSendToSlave(w, r, slaveMap)
	})
	go slaveMonitor.MonitorSlaves(3, slaveMap)

	log.Fatal(http.ListenAndServe(":5000", nil))
}
