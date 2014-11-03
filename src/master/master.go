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

	// TODO: need error handling if there are no slaves mapped
	slave := parseJson(POSTRequestBody)
	// slaveIPMap := initializeSlaveIPs() // this is creating the map each time a request is received
	destinationURL := destinationUrl(slave.ID, slaveIPMap)
	sendUrlValueMessageToSlave(destinationURL, slave.URL)
}

func receiveAndMapSlaveAddress(_ http.ResponseWriter, request *http.Request) {
	slaveIPAddress := request.PostFormValue("slaveIPAddress")
	fmt.Println("In receiveAndMapSlaveAddress, slaveIPAddress: ", slaveIPAddress)
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

func initializeSlaveIPs() (slaveIPMap map[string]string) {
	slaveIPs := make(map[string]string)
	slaveIPs["1"] = "http://10.0.0.42:8080"
	slaveIPs["2"] = "http://10.0.0.231:8080"

	return slaveIPs
}

func destinationUrl(slaveID string, slaveIPs map[string]string) (url string) {
	destination := url
	if slaveID == "1" {
		destination = slaveIPs["1"]
	}  else if slaveID == "2" {
		destination = slaveIPs["2"]
	}

	return destination
}
