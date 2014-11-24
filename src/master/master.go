package main

import (
	"master/master"
	"net/http"
)

var slaveIPMap = make(map[string]string)

func main() {
	slaveIPMap = master.SetUp()
	http.HandleFunc("/", master.ReceiveRequestAndSendToSlave)
	http.HandleFunc("/receive_slave", master.ReceiveAndMapSlaveAddress)
	http.HandleFunc("/receive_heartbeat", master.MonitorSlaveHeartbeats)
	go master.MonitorSlaves(3)
	http.ListenAndServe(":5000", nil)
}
