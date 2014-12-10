package main

import (
	"log"
	"master/master"
	"master/master/slaveMonitor"
	"master/master/slaveMapHandler"
	"net/http"
	"website"
	"github.com/gorilla/mux"
)

func main() {
	slaveMap := master.SetUp()
	router := mux.NewRouter()
	website.InitiateWebsiteHandlers(slaveMap, router)
	router.HandleFunc("/receive_heartbeat", func(_ http.ResponseWriter, r *http.Request) {
		slaveMap = slaveMonitor.ReceiveSlaveHeartbeat(r, slaveMap)
	})
	slaveMapHandler.InitiateSlaveMapHandler(router, slaveMap)

	http.Handle("/", router)
	go slaveMonitor.MonitorSlaves(3, slaveMap)
	log.Fatal(http.ListenAndServe(":5000", nil))
}
