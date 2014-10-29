package main

import (
	"net/http"
	"fmt"
)

func formatUrl () (string) {
	return "http://10.0.0.42:4000"
}

func handler(w http.ResponseWriter, r *http.Request) {
	formattedUrl := formatUrl()
	fmt.Println(formattedUrl, "Sending", r.PostFormValue("command"), "request.")
	// curl localhost:5000 -X POST -d "RPID=32423&url=http://9gag.com/"
	http.PostForm(formattedUrl, r.Form)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe("localhost:5000", nil)
}
