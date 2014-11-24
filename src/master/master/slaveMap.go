package master

import (
	"net/http"
	"fmt"
	"strings"
	"time"
	"encoding/json"
)

var webserverAddress = "http://localhost:4003"// TODO: make dynamic webserver address

// var slaveIPMap = make(map[string]string)
var slaveIPMap = initializeSlaveIPs()
var slaveHeartbeatMap = make(map[string]time.Time) 
// TODO: Create a single map with name as key and IP, heartbeat time as values.
// Should values be lists or hashmaps or objects?

type IdList struct {
	Id []string
}

func SetUp() (slaveMap map[string]string) {
	return slaveIPMap
}

func initializeSlaveIPs() (slaveIPMap map[string]string) {
	slaveIPs := make(map[string]string)
	slaveIPs["1"] = "http://10.0.0.122:8080"
	slaveIPs["2"] = "http://10.0.1.11:8080"

	return slaveIPs
}

func ReceiveAndMapSlaveAddress(_ http.ResponseWriter, request *http.Request) {
	slaveName := request.PostFormValue("slaveName")
	slaveURL := request.PostFormValue("slaveURL")
	fmt.Printf("\nNEW SLAVE RECEIVED.\n")
	fmt.Println("Slave Name: ", slaveName)
	fmt.Println("Slave URL: ", slaveURL)

	if returnedIPAddress, existsInMap := slaveIPMap[slaveName]; existsInMap {
		fmt.Printf("WARNING: Slave with name \"%v\" already exists with the IP address: %v. \nUpdating %v's IP address to %v.\n", slaveName, returnedIPAddress, slaveName, slaveURL)
	}
	slaveIPMap[slaveName] = slaveURL
	err := sendSlaveToWebserver(webserverAddress, slaveIPMap)
	printServerResponse(err, slaveName)

	slaveHeartbeatMap[slaveName] = time.Now()
	fmt.Printf("Mapped \"%v\" to %v.\n", slaveName, slaveURL)
	fmt.Println("Valid slave IDs are: ", slaveIPMap)
}

func printServerResponse(error error, slaveName string) {
	if error != nil {
		fmt.Printf("Error communicating with webserver: %v\n", error)
		fmt.Printf("%v not updated on webserver.\n", slaveName)
	} else {
		fmt.Printf("Added \"%v\" to webserver slave list.\n", slaveName)
	}
}

func MonitorSlaveHeartbeats(_ http.ResponseWriter, request *http.Request) {
	slaveName := request.PostFormValue("slaveName")
	heartbeatTimestamp := time.Now()
	slaveHeartbeatMap[slaveName] = heartbeatTimestamp
}

func MonitorSlaves(timeInterval int) {
	timer := time.Tick(time.Duration(timeInterval) * time.Second)   
    for _ = range timer {
		removeDeadSlaves(timeInterval)
    }
}

func removeDeadSlaves(deadTime int) {
	for slaveName, lastHeartbeatTime := range slaveHeartbeatMap {
		if time.Now().Sub(lastHeartbeatTime) > time.Duration(deadTime) * time.Second {
			fmt.Printf("\nREMOVING DEAD SLAVE: %v\n", slaveName)
			delete(slaveHeartbeatMap, slaveName)
			delete(slaveIPMap, slaveName)
			fmt.Println("Updated Slave Map: ", slaveIPMap)
			fmt.Printf("\n\n")
			sendSlaveToWebserver(webserverAddress, slaveIPMap)
		}
	}
}

func sendSlaveToWebserver(webserverAddress string, slaveIPs map[string]string) (err error) {
	err = nil
	client := &http.Client{}
	webserverAddress = webserverAddress + "/receive_slave"
	var idList IdList
	for slaveName := range slaveIPs {
        idList.Id = append(idList.Id, slaveName)
    }
	jsonMessage, err := json.Marshal(idList)
	_,err = client.Post(webserverAddress, "application/json", strings.NewReader(string(jsonMessage)))
	return err
}
