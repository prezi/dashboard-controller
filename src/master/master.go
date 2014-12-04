package main

import (
	"github.com/gorilla/mux"
	"log"
	"master/master"
	"master/master/slaveMonitor"
	"master/master/website"
	"master/master/website/session"
	"net/http"
)

func main() {
	slaveMap := master.SetUp()

	http.Handle("/assets/images/", http.StripPrefix("/assets/images/", http.FileServer(http.Dir(website.IMAGES_PATH))))
	http.Handle("/assets/javascripts/", http.StripPrefix("/assets/javascripts/", http.FileServer(http.Dir(website.JAVASCRIPTS_PATH))))
	http.Handle("/assets/stylesheets/", http.StripPrefix("/assets/stylesheets/", http.FileServer(http.Dir(website.STYLESHEETS_PATH))))

	router := mux.NewRouter()
	router.HandleFunc("/", website.IndexPageHandler)
	router.HandleFunc("/login", session.LoginHandler).Methods("POST")
	router.HandleFunc("/logout", session.LogoutHandler).Methods("POST")
	http.Handle("/", router)

	router.HandleFunc("/internal", func(w http.ResponseWriter, r *http.Request) {
		slaveNames := getSlaveNamesFromMap(slaveMap)
		website.FormHandler(w, r, slaveNames)
	})

	http.HandleFunc("/form-submit", func(w http.ResponseWriter, r *http.Request) {
		website.SubmitHandler(w, r, slaveMap)
	})

	http.HandleFunc("/receive_heartbeat", func(_ http.ResponseWriter, r *http.Request) {
		slaveMap = slaveMonitor.ReceiveSlaveHeartbeat(r, slaveMap)
	})

	go slaveMonitor.MonitorSlaves(3, slaveMap)
	log.Fatal(http.ListenAndServe(":5000", nil))
}

func getSlaveNamesFromMap(slaveMap map[string]master.Slave) (slaveNames []string) {
	for k := range slaveMap {
		slaveNames = append(slaveNames, k)
	}
	return
}
