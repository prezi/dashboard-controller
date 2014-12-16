package slaveMapHandler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"master/master"
	"net/http"
	"sort"
)

type SlaveMap struct {
	SlaveMap map[string]master.Slave
}

func InitiateSlaveMapHandler(slaveMap map[string]master.Slave, router *mux.Router) {
	router.HandleFunc("/slavemap", func(responseWriter http.ResponseWriter, request *http.Request) {
		slavemapHandler(responseWriter, request, slaveMap)
	})
}

func slavemapHandler(responseWriter http.ResponseWriter, _ *http.Request, slaveMap map[string]master.Slave) {
	slaveNames := GetSlaveNamesFromMap(slaveMap)
	responseWriter.Header().Set("Content-Type", "application/json")

	slaveNameJson, error := json.Marshal(slaveNames)
	if error != nil {
		fmt.Println(error)
	}
	responseWriter.Write(slaveNameJson)
}

func GetSlaveMap() (slaveMap map[string]master.Slave) {
	slaveMap = make(map[string]master.Slave)
	return
}

func GetSlaveNamesFromMap(slaveMap map[string]master.Slave) (slaveNames []string) {
	for index := range slaveMap {
		slaveNames = append(slaveNames, index)
	}
	sort.Strings(slaveNames)
	return
}
