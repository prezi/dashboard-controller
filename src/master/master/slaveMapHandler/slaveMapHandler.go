package slaveMapHandler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"master/master"
	"net/http"
	"website"
)

type SlaveMap struct {
	SlaveMap map[string]master.Slave
}

func InitiateSlaveMapHandler(router *mux.Router, slaveMap map[string]master.Slave) {
	router.HandleFunc("/slavemap", func(responseWriter http.ResponseWriter, request *http.Request) {
		slavemapHandler(responseWriter, request, slaveMap)
	})
}

func slavemapHandler(responseWriter http.ResponseWriter, _ *http.Request, slaveMap map[string]master.Slave) {
	slaveNames := website.GetSlaveNamesFromMap(slaveMap)
	responseWriter.Header().Set("Content-Type", "application/json")

	slaveNameJson, error := json.Marshal(slaveNames)
	if error != nil {
		fmt.Println(error)
	}
	responseWriter.Write(slaveNameJson)
}
