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
	POSTRequestBody, _ := ioutil.ReadAll(request.Body)
	defer request.Body.Close()

	slave := parseJson(POSTRequestBody)
	slaveIPMap := initializeSlaveIPs()
	destinationURL := destinationUrl(slave.ID, slaveIPMap)
	sendUrlValueMessageToServer(destinationURL, slave.URL)
}

func parseJson(input []byte) (slave Slave) {
	err := json.Unmarshal(input, &slave)
	if err != nil {
		fmt.Println("error:", err)
	}
	return
}

func sendUrlValueMessageToServer( slaveURL string, urlToDisplay string) {
	client := &http.Client{}

	form := url.Values{}
	form.Set("url", urlToDisplay)

	_,_ = client.PostForm(slaveURL, form)
}

func initializeSlaveIPs() (slaveIPMap map[string]string) {
	raspberryPiIP := make(map[string]string)
	raspberryPiIP["1"] = "http://10.0.0.42:8080"
	raspberryPiIP["2"] = "http://10.0.0.231:8080"

	return raspberryPiIP
}

func destinationUrl(slaveID string, slaveIPMap map[string]string) (url string) {
	destination := url
	if slaveID == "1" {
		destination = slaveIPMap["1"]
	}  else if slaveID == "2" {
		destination = slaveIPMap["2"]
	}

	return destination
}
