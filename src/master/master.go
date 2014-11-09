package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type Slave struct {
	ID string
	URL string
}

var slaveIPMap = make(map[string]string)

func main() {
	slaveIPMap = initializeSlaveIPs()
	http.HandleFunc("/", handler)
	http.HandleFunc("/receive_slave", receiveAndMapSlaveAddress)
	http.ListenAndServe("localhost:5000", nil)
}

func handler(_ http.ResponseWriter, request *http.Request) {

	POSTRequestBody, _ := ioutil.ReadAll(request.Body)
	defer request.Body.Close()

	slave := parseJson(POSTRequestBody)
	destinationSlaveAddress := destinationSlaveAddress(slave.ID)
	if destinationSlaveAddress == "" {
		fmt.Println("Abandoning request.")
		return
	}

	fmt.Printf("\nSending %v to %v at %v", slave.URL, slave.ID, destinationSlaveAddress)
	sendUrlValueMessageToSlave(destinationSlaveAddress, slave.URL)
}

func initializeSlaveIPs() (slaveIPMap map[string]string) {
	slaveIPs := make(map[string]string)
	slaveIPs["1"] = "http://10.0.0.42:8080"
	slaveIPs["2"] = "http://10.0.0.231:8080"

	return slaveIPs
}

func receiveAndMapSlaveAddress(_ http.ResponseWriter, request *http.Request) {
	slaveName := request.PostFormValue("slaveName")
	slaveIPAddress := request.PostFormValue("slaveIPAddress")
	fmt.Printf("\nNEW SLAVE RECEIVED.\n")
	fmt.Println("Slave Name: ", slaveName)
	fmt.Println("Slave IP address: ", slaveIPAddress)


	if returnedIPAddress, existsInMap := slaveIPMap[slaveName]; existsInMap == false {
		// send new slave name to webserver
		client := &http.Client{}
		webserverIPAddressAndExtentionArray := []string{"http://localhost:4003", "/receive_slave"} 
		webserverReceiveSlaveAddress := strings.Join(webserverIPAddressAndExtentionArray, "")
	
		form := url.Values{}
		form.Set("slaveName", slaveName)
		_, err := client.PostForm(webserverReceiveSlaveAddress, form)
	
		if err != nil {
			fmt.Printf("Error communicating with webserver: %v\n", err)
			fmt.Printf("%v not updated on webserver.\n", slaveName)
		} else {
			fmt.Printf("Added \"%v\" to webserver slave list.\n", slaveName)
		}
	} else {
		fmt.Printf("WARNING: Slave with name \"%v\" already exists with the IP address: %v. \nUpdating %v's IP address to %v.\n", slaveName, returnedIPAddress, slaveName, slaveIPAddress)
	}
	slaveIPMap[slaveName] = slaveIPAddress
	fmt.Printf("Mapped \"%v\" to %v.\n", slaveName, slaveIPAddress)
	fmt.Println("Valid slave IDs are: ", slaveIPMap)
}

// TODO: this doesn't return anything...?
func parseJson(input []byte) (slave Slave) {
	err := json.Unmarshal(input, &slave)
	if err != nil {
		fmt.Println("error:", err)
	}
	return
}

func sendUrlValueMessageToSlave(slaveIPAddress string, urlToDisplay string) {
	client := &http.Client{}

	form := url.Values{}
	form.Set("url", urlToDisplay)

	_,_ = client.PostForm(slaveIPAddress, form)
}

func destinationSlaveAddress(slaveID string) (slaveAddress string) {
	if len(slaveIPMap) == 0 {
		fmt.Println("ERROR: No slaves available.")
		return 
	}

	slaveAddress = slaveIPMap[slaveID]
	if slaveAddress ==  "" {
		fmt.Printf("ERROR: \"%v\" is not a valid slave ID.\n", slaveID)
		fmt.Println("Valid slave IDs are: ", slaveIPMap)
		return 
	}
	return slaveAddress
}
