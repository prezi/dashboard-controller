package main

import (
	"log"
	"net/http"
	"proxy/manageIPTables"
	"proxy/proxy"
)

var PROXY_HTTP_SERVER_PORT = "7878"

func main() {
	masterURL := proxy.SetUp()
	go proxy.Heartbeat(1, masterURL)

	http.HandleFunc("/update_iptables", func(_ http.ResponseWriter, r *http.Request) {
		manageIPTables.UpdateIPTables(r)
	})

	log.Fatal(http.ListenAndServe(":"+PROXY_HTTP_SERVER_PORT, nil))
}
