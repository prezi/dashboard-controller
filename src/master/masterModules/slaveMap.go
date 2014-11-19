package masterModule

import (
	"net/http"
	"fmt"
	"strings"
	// "net/url"
	"time"
	"encoding/json"
)

// var slaveIPMap = make(map[string]string)
var slaveIPMap = initializeSlaveIPs()
var slaveHeartbeatMap = make(map[string]time.Time) 
// TODO: Create a single map with name as key and IP, heartbeat time as values.
// Should values be tuples or hashmaps?

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
	slaveIPAddress := request.PostFormValue("slaveIPAddress")
	fmt.Printf("\nNEW SLAVE RECEIVED.\n")
	fmt.Println("Slave Name: ", slaveName)
	fmt.Println("Slave IP address: ", slaveIPAddress)




	if returnedIPAddress, existsInMap := slaveIPMap[slaveName]; existsInMap {

		fmt.Printf("WARNING: Slave with name \"%v\" already exists with the IP address: %v. \nUpdating %v's IP address to %v.\n", slaveName, returnedIPAddress, slaveName, slaveIPAddress)
		
	}
	slaveIPMap[slaveName] = slaveIPAddress
		
	webserverIPAddressAndExtentionArray := []string{"http://localhost:4003", "/receive_slave"} // TODO: make dynamic webserver address
	err := sendSlaveToWebserver(webserverIPAddressAndExtentionArray, slaveIPMap)
	printServerResponse(err, slaveName)
	slaveHeartbeatMap[slaveName] = time.Now()
	fmt.Printf("Mapped \"%v\" to %v.\n", slaveName, slaveIPAddress)
	fmt.Println("Valid slave IDs are: ", slaveIPMap)
}

func MonitorSlaveHeartbeats(_ http.ResponseWriter, request *http.Request) {
	slaveName := request.PostFormValue("slaveName")
	heartbeatTimestamp := request.PostFormValue("heartbeatTimestamp")

	timeFormat := "2006-01-02 15:04:05.999999999 -0700 MST"
	heartbeatTime, err := time.Parse(timeFormat, heartbeatTimestamp)

	if err != nil {
		fmt.Println("Error encountered when parsing heartbeat timestamp from slave.")
		fmt.Println("ERROR: ", err)
	}

	slaveHeartbeatMap[slaveName] = heartbeatTime

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
		}
	}
}

func sendSlaveToWebserver(webserverIPAddressAndExtentionArray []string, slaveIPs map[string]string) (err error) {
	err = nil
	client := &http.Client{}
	webserverReceiveSlaveAddress := strings.Join(webserverIPAddressAndExtentionArray, "")
	var idList IdList
	for slaveName := range slaveIPs {
        idList.Id = append(idList.Id, slaveName)
    }
	jsonMessage, err := json.Marshal(idList)
	_,err = client.Post(webserverReceiveSlaveAddress, "application/json", strings.NewReader(string(jsonMessage)))
	return err

}

func printServerResponse(error error, slaveName string) {
	if error != nil {
		fmt.Printf("Error communicating with webserver: %v\n", error)
		fmt.Printf("%v not updated on webserver.\n", slaveName)
	} else {
		fmt.Printf("Added \"%v\" to webserver slave list.\n", slaveName)
	}
}
