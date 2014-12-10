package main

import (
	"log"
	"master/master"
	"master/master/slaveMonitor"
	"master/master/slaveMapHandler"
	"net/http"
	"website"
)

func main() {
	slaveMap := master.SetUp()
	website.InitiateWebsiteHandlers(slaveMap)
	http.HandleFunc("/receive_heartbeat", func(_ http.ResponseWriter, r *http.Request) {
		slaveMap = slaveMonitor.ReceiveSlaveHeartbeat(r, slaveMap)
	})
	slaveMapHandler.InitiateSlaveMapHandler(slaveMap)

	go slaveMonitor.MonitorSlaves(3, slaveMap)
	log.Fatal(http.ListenAndServe(":5000", nil))
}
