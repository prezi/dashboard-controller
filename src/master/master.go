package main

import (
	"github.com/gorilla/mux"
	"log"
	"master/master"
	"master/master/proxyMonitor"
	"master/master/slaveMapHandler"
	"master/master/slaveMonitor"
	"net/http"
	"network"
	"website"
)

var (
	SLAVE_BINARY_PATH = network.GetRelativeFilePath("../../bin/slave")
	MASTER_PORT       = "5000"
)

func main() {
	slaveMap := master.GetSlaveMap()
	router := mux.NewRouter()
	website.InitiateWebsiteHandlers(slaveMap, router)
	router.HandleFunc("/receive_heartbeat", func(_ http.ResponseWriter, r *http.Request) {
		slaveMap = slaveMonitor.ReceiveSlaveHeartbeat(r, slaveMap)
	})
	router.HandleFunc("/receive_proxy_heartbeat", func(_ http.ResponseWriter, r *http.Request) {
		proxyMonitor.ReceiveProxyHeartbeat(r)
	})
	router.HandleFunc("/get_slave_binary", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, SLAVE_BINARY_PATH)
	})

	slaveMapHandler.InitiateSlaveMapHandler(router, slaveMap)

	http.Handle("/", router)
	go slaveMonitor.MonitorSlaves(3, slaveMap)
	go proxyMonitor.MonitorProxy(10)
	log.Fatal(http.ListenAndServe(":"+MASTER_PORT, nil))
}
