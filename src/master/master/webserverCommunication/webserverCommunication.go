package webserverCommunication

import (
	"encoding/json"
	"fmt"
	"master/master"
	"net"
	"net/http"
	"network"
	"strings"
)

type IdList struct {
	Id []string
}

// If webserver address updates, will the slave map update properly as well?
// What if we have more than one webserver pinging the master?
func UpdateWebserverAddress(r *http.Request, webServerAddress string) (newWebServerAddress string, err error) {
	newWebServerAddress, err = getWebserverAddress(r)
	if webServerAddress != newWebServerAddress {
		fmt.Printf("Webserver address has changed from %v to %v.\n", webServerAddress, newWebServerAddress)
		webServerAddress = newWebServerAddress
	}
	return
}

func getWebserverAddress(request *http.Request) (webServerAddress string, err error) {
	slaveIP, _, err := net.SplitHostPort(request.RemoteAddr)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
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

func SendWebserverInit(r *http.Request, slaveMap map[string]master.Slave) (UpdatedWebServerAddress string) {
	if r.FormValue("message") == "update me!" {
		UpdatedWebServerAddress, _ = getWebserverAddress(r)
		err := SendSlaveListToWebserver(UpdatedWebServerAddress, slaveMap)
		network.ErrorHandler(err, "Error: %v\n")
	}
	return
}

func SendSlaveListToWebserver(webServerAddress string, slaveMap map[string]master.Slave) (err error) {
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
