package main

import (
	"net/http"
	"fmt"
)

type RaspberryPiIP struct {
	IP string
}

var rapsberryPiIP = RaspberryPiIP{"http://10.0.0.231:8080"}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Sending POST request to", rapsberryPiIP.IP, "with url", r.Form)
	// curl localhost:5000 -X POST -d "RPID=32423&url=http://9gag.com/"

	fmt.Println(r)
	http.PostForm(rapsberryPiIP.IP, r.Form)
}

func main() {
	raspberryPiIP := make(map[string]string)
	raspberryPiIP["1"] = "http://10.0.0.42:8080"
	raspberryPiIP["2"] = "http://10.0.0.231:8080"

	http.HandleFunc("/", handler)
	http.ListenAndServe("localhost:5000", nil)
}
