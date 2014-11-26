package master

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"net"
	"errors"
)

var webserverAddress = "http://localhost:4003"// TODO: make dynamic webserver address

type Slave struct {
	URL          string
	heartbeat    time.Time
	displayedURL string // TODO: store currently displayed URL for each slave
}

type IdList struct {
	Id []string
}

// TODO: Allow user to input names of currently running slaves at master startup.
// Alternatively, allow to manually add names of currently running slave while the master is running.
func SetUp() (slaveMap map[string]Slave) {
	return initializeSlaveMap()
}

func initializeSlaveMap() (slaveMap map[string]Slave) {
	slaveMap = make(map[string]Slave)
	return slaveMap
}

func printServerResponse(error error, slaveName string) {
	if error != nil {
		fmt.Printf("Error communicating with webserver: %v\n", error)
		fmt.Printf("%v not updated on webserver.\n", slaveName)
	} else {
		fmt.Printf("Added \"%v\" to webserver slave list.\n", slaveName)
	}
}

func sendSlaveListToWebserver(webserverAddress string, slaveMap map[string]Slave) (err error) {
	err = nil
	client := &http.Client{}
	webserverAddress = webserverAddress + "/receive_slave"
	var idList IdList
	for slaveName := range slaveMap {
		idList.Id = append(idList.Id, slaveName)
	}
	jsonMessage, err := json.Marshal(idList)
	_, err = client.Post(webserverAddress, "application/json", strings.NewReader(string(jsonMessage)))
	return err
}

func WebserverRequestSlaveListAtStartUp(writer http.ResponseWriter, request *http.Request, slaveMap map[string]Slave) {
	newWebserverAddress, err := getWebserverAddress(request)
	if err != nil {
		writer.WriteHeader(500)
	} else {
		writer.WriteHeader(200)
		if newWebserverAddress != webserverAddress {
			webserverAddress = newWebserverAddress
			fmt.Printf(`############## WebserverURL: %v`, webserverAddress)
			sendSlaveListToWebserver(webserverAddress, slaveMap)
		}
	}
}

func getWebserverAddress(request *http.Request) (webserverAddress string,err error) {
	err = nil
	slaveIP,_,_ := net.SplitHostPort(request.RemoteAddr)
	webserverPort := request.PostFormValue("webserverPort")
	if webserverPort == "" 	{
		err = errors.New("Cannot find sender port.")
		return
	}
	webserverAddress = "http://" + slaveIP + ":" + webserverPort
	return
}
