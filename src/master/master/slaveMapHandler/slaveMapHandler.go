package slaveMapHandler

import (
	"master/master"
	"net/http"
	"fmt"
	"encoding/json"
	"website"
)

type SlaveMap struct {
	SlaveMap map[string]master.Slave
}

func InitiateSlaveMapHandler(slaveMap map[string]master.Slave) {
	http.HandleFunc("/slavemap", func(responseWriter http.ResponseWriter, request *http.Request) {
		slavemapHandler(responseWriter, request, slaveMap)
	})
}

func slavemapHandler(responseWriter http.ResponseWriter, _ *http.Request, slaveMap map[string]master.Slave) {
	slaveNames := website.GetSlaveNamesFromMap(slaveMap)
	responseWriter.Header().Set("Content-Type", "application/json")
	jsonMessage, err := json.Marshal(slaveNames)
	if err != nil {
		fmt.Println(err)
	}
	responseWriter.Write(jsonMessage)
}
