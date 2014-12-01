package master

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
)

var webServerAddress = "http://localhost:4003"

type IdList struct {
	Id []string
}

func UpdateWebserverAddress(r *http.Request) (err error) {
	newWebServerAddress, err := getWebserverAddress(r)
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

func SendWebserverInit(r *http.Request, slaveMap map[string]Slave) {
	if r.FormValue("message") == "update me!" {
		webServerAddress, _ = getWebserverAddress(r)
		fmt.Println(sendSlaveListToWebserver(webServerAddress, slaveMap))
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
