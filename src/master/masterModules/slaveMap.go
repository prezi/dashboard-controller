package masterModule

import (
	"net/http"
	"fmt"
	"strings"
	"net/url"
)


// var slaveIPMap = make(map[string]string)
var slaveIPMap = initializeSlaveIPs()
var slaveHeartbeatMap = make(map[string]string) // TODO: make these map to time values

func SetUp() (slaveMap map[string]string) {
	return slaveIPMap
}

func ReceiveAndMapSlaveAddress(_ http.ResponseWriter, request *http.Request) {
	slaveName := request.PostFormValue("slaveName")
	slaveIPAddress := request.PostFormValue("slaveIPAddress")
	fmt.Printf("\nNEW SLAVE RECEIVED.\n")
	fmt.Println("Slave Name: ", slaveName)
	fmt.Println("Slave IP address: ", slaveIPAddress)

	if returnedIPAddress, existsInMap := slaveIPMap[slaveName]; existsInMap == false {
		webserverIPAddressAndExtentionArray := []string{"http://localhost:4003", "/receive_slave"}

		err := sendSlaveToWebserver(webserverIPAddressAndExtentionArray, slaveName)
		printServerResponse(err, slaveName)
	} else {
		fmt.Printf("WARNING: Slave with name \"%v\" already exists with the IP address: %v. \nUpdating %v's IP address to %v.\n", slaveName, returnedIPAddress, slaveName, slaveIPAddress)
	}
	slaveIPMap[slaveName] = slaveIPAddress
	fmt.Printf("Mapped \"%v\" to %v.\n", slaveName, slaveIPAddress)
	fmt.Println("Valid slave IDs are: ", slaveIPMap)
}

// func monitorSlaveHeartbeats(_ http.ResponseWriter, request *http.Request) {
// 	slaveName := request.PostFormValue("slaveName")
// 	heartbeatTimestamp := request.PostFormValue("heartbeatTimestamp")
// }

func removeDeadSlaves() {

}



func sendSlaveToWebserver(webserverIPAddressAndExtentionArray []string, slaveName string) (err error) {
	client := &http.Client{}
	webserverReceiveSlaveAddress := strings.Join(webserverIPAddressAndExtentionArray, "")

	form := url.Values{}
	form.Set("slaveName", slaveName)
	_, err = client.PostForm(webserverReceiveSlaveAddress, form)

	printServerResponse(err, slaveName)

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

func initializeSlaveIPs() (slaveIPMap map[string]string) {
	slaveIPs := make(map[string]string)
	slaveIPs["1"] = "http://10.0.0.122:8080"
	slaveIPs["2"] = "http://10.0.1.11:8080"

	return slaveIPs
}

