package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

type Slave struct {
	ID string
	URL string
}

func handler(responseWriter http.ResponseWriter, request *http.Request) {
	POSTrequestBody, _ := ioutil.ReadAll(request.Body)
	var slave Slave
	_ = json.Unmarshal(POSTrequestBody, &slave)
	fmt.Println(slave.ID)
	fmt.Println(slave.URL)

	raspberryPiIP := make(map[string]string)
	raspberryPiIP["1"] = "http://10.0.0.42:8080"
	raspberryPiIP["2"] = "http://10.0.0.231:8080"

	var slaveAddress string
	if slave.ID == "1" {
		slaveAddress = raspberryPiIP["1"]
	}  else if slave.ID == "2" {
		slaveAddress = raspberryPiIP["2"]
	}

	fmt.Println(slaveAddress)
//	http.PostForm(slaveAddress, request.Form)
}

func main() {


	http.HandleFunc("/", handler)
	http.ListenAndServe("localhost:5000", nil)
}
