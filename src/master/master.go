package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"net/url"
)

type Slave struct {
	ID string
	URL string
}

var slaveIPMap = make(map[string]string)

func main() {
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

	sendUrlValueMessageToSlave(destinationSlaveAddress, slave.URL)
}

func receiveAndMapSlaveAddress(_ http.ResponseWriter, request *http.Request) {
	slaveIPAddress := request.PostFormValue("slaveIPAddress")
	fmt.Println("Slave IP address received: ", slaveIPAddress)
	slaveIPMap["1"] = slaveIPAddress
	fmt.Println("MAPPED: ", slaveIPAddress)
}

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
	fmt.Println("slaveIPAddress: ", slaveIPAddress)

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
	fmt.Println("slaveAddress in destinationSlaveAddress: ", slaveAddress)
	return slaveAddress
}
