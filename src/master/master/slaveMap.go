package master

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

var webServerAddress = "http://localhost:4003" // TODO: make dynamic webserver address

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
// I don't think this should be in the scope.. Renaming the slave is laborious, and the user could easily restart it..
func SetUp() (slaveMap map[string]Slave) {
	slaveMap = make(map[string]Slave)
	return
}

func printServerResponse(error error, slaveName string) {
	if error != nil {
		fmt.Printf("Error communicating with webserver: %v\n", error)
		fmt.Printf("%v not updated on webserver.\n", slaveName)
	} else {
		fmt.Printf("Added \"%v\" to webserver slave list.\n", slaveName)
	}
}

func sendSlaveListToWebserver(webServerAddress string, slaveMap map[string]Slave) (err error) {
	err = nil
	client := &http.Client{}
	webServerAddress = webServerAddress + "/receive_slave"
	var idList IdList
	for slaveName := range slaveMap {
		idList.Id = append(idList.Id, slaveName)
	}
	jsonMessage, err := json.Marshal(idList)
	_, err = client.Post(webServerAddress, "application/json", strings.NewReader(string(jsonMessage)))
	return err
}

func getWebserverAddress(request *http.Request) (webServerAddress string, err error) {
	slaveIP, _, err := net.SplitHostPort(request.RemoteAddr)

	if err != nil {
		fmt.Printf("Error: %v\n",err)
		return
	}
	webServerAddress = "http://" + slaveIP

	webServerPort := request.PostFormValue("webserverPort")
	if webServerPort != "" {
		webServerAddress += ":" + webServerPort
	} else {
		fmt.Println("port number not found.")
		return
	}

	return
}

func SendWebserverInit(r *http.Request, slaveMap map[string]Slave) {
	if r.FormValue("message") == "update me!" {
		webServerAddress, _ = getWebserverAddress(r)
		fmt.Println(sendSlaveListToWebserver(webServerAddress, slaveMap))
	}
}
