package main

import (
	"net/http"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"os"
)

var protocol string
var host string
var port string

type initVars struct {
	Protocol string
	Host string
	Port string
}


func jsonInit(){
	// read the config file, report any errors
	file, e := ioutil.ReadFile("./serverConfig.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	var jsontype initVars

	// Decode the JSON
	err := json.Unmarshal(file, &jsontype)
	if err != nil {
		fmt.Printf("Unmarshal error: %v\n", e)
		os.Exit(1)
	}

	// set variables from JSON file
	protocol = jsontype.Protocol
	host = jsontype.Host
	port = jsontype.Port
}

// create url that sends command to the designated Raspberry Pi
func formatUrl (rasPiNum string) (string) {
	finalUrl := protocol + rasPiNum + port

	return finalUrl
}

func handler(w http.ResponseWriter, r *http.Request) {
	rasPiNum := r.PostFormValue("RPID")

	formattedUrl := formatUrl(rasPiNum)

	fmt.Println(formattedUrl, "Sending", r.PostFormValue("command"), "request.")

	http.PostForm(formattedUrl, r.Form)
}


func main() {
	jsonInit()
	http.HandleFunc("/", handler) // redirect all commands to the handler function
	http.ListenAndServe(host + port, nil) // listen for connections at port 9999 on the local machine
}
