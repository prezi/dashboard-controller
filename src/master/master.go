package main

import (
	"net/http"
	"fmt"
//	"encoding/json"
//	"io/ioutil"
//	"os"
)

//var protocol string
var host string
var port string

type initVars struct {
	Protocol string
	Host string
	Port string
}

func formatUrl (rasPiNum string) (string) {
//	finalUrl := protocol + rasPiNum + port
	return "http://10.0.0.42:4000"
}

func handler(w http.ResponseWriter, r *http.Request) {
	rasPiNum := r.PostFormValue("RPID")
//	urlToPOSTToSlave := r.PostFormValue("url")
	formattedUrl := formatUrl(rasPiNum)
	fmt.Println(formattedUrl, "Sending", r.PostFormValue("command"), "request.") //send requests by
	// curl localhost:4001 -X POST -d "RPID=32423&url=http://9gag.com/"
	http.PostForm(formattedUrl, r.Form)
}


func main() {
//	jsonInit()
	http.HandleFunc("/", handler) // redirect all commands to the handler function
	http.ListenAndServe("localhost:4001", nil) // listen for connections at port 9999 on the local machine
}
