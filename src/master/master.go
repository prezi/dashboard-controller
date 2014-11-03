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

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe("localhost:5000", nil)
}

func handler(_ http.ResponseWriter, request *http.Request) {

	receiveAndMapSlaveAddress(request)


	POSTRequestBody, _ := ioutil.ReadAll(request.Body)
	defer request.Body.Close()

	slave := parseJson(POSTRequestBody)
	slaveIPMap := initializeSlaveIPs()
	destinationURL := destinationUrl(slave.ID, slaveIPMap)
	sendUrlValueMessageToSlave(destinationURL, slave.URL)
}

func receiveAndMapSlaveAddress(request *http.Request) {
	slaveIPAddress := request.PostFormValue("slaveIPAddress")
	fmt.Println("slaveIPAddress: ", slaveIPAddress)
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
