package website

import (
	"github.com/gorilla/mux"
	"master/master"
	"master/master/website/session"
	"net/http"
)

func InitiateWebsiteHandlers(slaveMap map[string]master.Slave) {
	http.Handle("/assets/images/", http.StripPrefix("/assets/images/", http.FileServer(http.Dir(IMAGES_PATH))))
	http.Handle("/assets/javascripts/", http.StripPrefix("/assets/javascripts/", http.FileServer(http.Dir(JAVASCRIPTS_PATH))))
	http.Handle("/assets/stylesheets/", http.StripPrefix("/assets/stylesheets/", http.FileServer(http.Dir(STYLESHEETS_PATH))))

	router := mux.NewRouter()
	router.HandleFunc("/", IndexPageHandler)
	router.HandleFunc("/login", session.LoginHandler).Methods("POST")
	router.HandleFunc("/logout", session.LogoutHandler).Methods("POST")
	http.Handle("/", router)

	router.HandleFunc("/internal", func(w http.ResponseWriter, r *http.Request) {
		slaveNames := getSlaveNamesFromMap(slaveMap)
		FormHandler(w, r, slaveNames)
	})
	http.HandleFunc("/form-submit", func(w http.ResponseWriter, r *http.Request) {
		SubmitHandler(w, r, slaveMap)
	})
}

func getSlaveNamesFromMap(slaveMap map[string]master.Slave) (slaveNames []string) {
	for k := range slaveMap {
		slaveNames = append(slaveNames, k)
	}
	return
}
