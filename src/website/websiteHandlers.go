package website

import (
	"github.com/gorilla/mux"
	"master/master"
	"net/http"
	"website/session"
)

func InitiateWebsiteHandlers(slaveMap map[string]master.Slave, router *mux.Router, filePathToUserAuthenticationData string) {
	http.Handle("/assets/images/", http.StripPrefix("/assets/images/", http.FileServer(http.Dir(IMAGES_PATH))))
	http.Handle("/assets/javascripts/", http.StripPrefix("/assets/javascripts/", http.FileServer(http.Dir(JAVASCRIPTS_PATH))))
	http.Handle("/assets/stylesheets/", http.StripPrefix("/assets/stylesheets/", http.FileServer(http.Dir(STYLESHEETS_PATH))))

	router.HandleFunc("/", IndexPageHandler)
	router.HandleFunc("/login", func(responseWriter http.ResponseWriter, request *http.Request) {
			session.LoginHandler(responseWriter, request, filePathToUserAuthenticationData)
		}).Methods("POST")

	router.HandleFunc("/logout", session.LogoutHandler).Methods("POST")

	router.HandleFunc("/internal", func(w http.ResponseWriter, r *http.Request) {
		slaveNames := master.GetSlaveNamesFromMap(slaveMap)
		FormHandler(w, r, slaveNames)
	})
	router.HandleFunc("/form-submit", func(w http.ResponseWriter, r *http.Request) {
		SubmitHandler(w, r, slaveMap)
	})
}
