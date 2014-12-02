package main

import (
	"log"
	"master/master"
	"master/master/receiveAndSendRequestToSlave"
	"master/master/slaveMonitor"
	"master/master/webserverCommunication"
	"net/http"
)

func main() {
	slaveMap, webServerAddress := master.SetUp()

	http.HandleFunc("/receive_heartbeat", func(_ http.ResponseWriter, r *http.Request) {
		slaveMap = slaveMonitor.ReceiveSlaveHeartbeat(r, slaveMap, webServerAddress)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		receiveAndSendRequestToSlave.ReceiveRequestAndSendToSlave(w, r, slaveMap)
	})
	http.HandleFunc("/webserver_heartbeat", func(_ http.ResponseWriter, r *http.Request) {
		webServerAddress = webserverCommunication.UpdateWebserverAddress(r, webServerAddress)
	})
	http.HandleFunc("/webserver_init", func(_ http.ResponseWriter, r *http.Request) {
		webServerAddress = webserverCommunication.SendWebserverInit(r, slaveMap)
	})
	go slaveMonitor.MonitorSlaves(3, slaveMap, webServerAddress)
	log.Fatal(http.ListenAndServe(":5000", nil))
}
