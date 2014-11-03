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
	//	defer request.Body.Close()

	var slave Slave
	slave = parseJson(POSTRequestBody)

	raspberryPiIP := make(map[string]string)
	raspberryPiIP["1"] = "http://10.0.0.42:8080"
	raspberryPiIP["2"] = "http://10.0.0.231:8080"

	var destinationURL string
	if slave.ID == "1" {
		destinationURL = raspberryPiIP["1"]
	}  else if slave.ID == "2" {
		destinationURL = raspberryPiIP["2"]
	}

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
