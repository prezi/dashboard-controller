package main

import (
	"log"
	"master/master"
	"master/master/receiveAndSendRequestToSlave"
	"master/master/slaveMonitor"
	"master/master/website"
	"net/http"
)

func main() {
	slaveMap := master.SetUp()

	http.Handle("/assets/images/", http.StripPrefix("/assets/images/", http.FileServer(http.Dir(website.IMAGES_PATH))))
	http.Handle("/assets/javascripts/", http.StripPrefix("/assets/javascripts/", http.FileServer(http.Dir(website.JAVASCRIPTS_PATH))))
	http.Handle("/assets/stylesheets/", http.StripPrefix("/assets/stylesheets/", http.FileServer(http.Dir(website.STYLESHEETS_PATH))))
	http.HandleFunc("/",  func(w http.ResponseWriter, r *http.Request) {
		slaveNames := make([]string, 0, len(slaveMap))
		for k := range slaveMap {
			slaveNames = append(slaveNames, k)
		}
			website.FormHandler(w, r, slaveNames)
		})
	http.HandleFunc("/form-submit", func(w http.ResponseWriter, r *http.Request) {
			website.SubmitHandler(w, r)
		})

	http.HandleFunc("/receive_heartbeat", func(_ http.ResponseWriter, r *http.Request) {
		slaveMap = slaveMonitor.ReceiveSlaveHeartbeat(r, slaveMap)
	})
	http.HandleFunc("/send_url_to_slave", func(w http.ResponseWriter, r *http.Request) {
		receiveAndSendRequestToSlave.ReceiveRequestAndSendToSlave(w, r, slaveMap)
	})
	go slaveMonitor.MonitorSlaves(3, slaveMap)
	log.Fatal(http.ListenAndServe(":5000", nil))
}
