package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
//	"log"
	"encoding/json"
)

type RaspberryPiIP struct {
	IP string
}

type rapsberryPiData struct {
	ID string
	URL string
}

var rapsberryPiIP = RaspberryPiIP{"http://10.0.0.231:8080"}

func handler(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(body))
//	curl -d '{"IP":"1"}' -H "Content-Type: application/json" http://localhost:5000
//returns: 2014/10/29 21:57:39 {"IP":"1"}

	var rpip RaspberryPiIP
	_ = json.Unmarshal(body, &rpip)
	fmt.Println(rpip.IP)


//	fmt.Println(req.PostFormValue("URL"))
//	fmt.Println(req.PostFormValue("URL"))
	// curl -X POST -d "URL=http://google.com" http://localhost:5000
	// returns: http://google.com
}

func main() {
	raspberryPiIP := make(map[string]string)
	raspberryPiIP["1"] = "http://10.0.0.42:8080"
	raspberryPiIP["2"] = "http://10.0.0.231:8080"

	http.HandleFunc("/", handler)
	http.ListenAndServe("localhost:5000", nil)
}
