package main

import (
	"log"
	"network"
	"master/master/slaveMonitor"
	"master/master/slaveMapHandler"
	"net/http"
	"website"
	"github.com/gorilla/mux"
)

func main() {
	slaveMap := slaveMapHandler.GetSlaveMap()
	router := mux.NewRouter()

	website.InitiateWebsiteHandlers(slaveMap, router)

	router.HandleFunc("/receive_heartbeat", func(_ http.ResponseWriter, r *http.Request) {
		slaveMap = slaveMonitor.ReceiveSlaveHeartbeat(r, slaveMap)
	})
	router.HandleFunc("/get_slave_binary", func(responseWriter http.ResponseWriter, request *http.Request) {
		http.ServeFile(responseWriter, request, network.GetRelativeFilePath("../../bin/slave"))
	})

	slaveMapHandler.InitiateSlaveMapHandler(slaveMap, router)

	http.Handle("/", router)
	go slaveMonitor.MonitorSlaves(3, slaveMap)
	log.Fatal(http.ListenAndServe(":5000", nil))
}
