package main

import (
	"strings"
	"fmt"
	"net/http"
	"log"
	"mime"
	"path/filepath"
	"html/template"
	"encoding/json"
)

type Message struct { // this will be the json: { "ID": "1", "URL": "http://google.com"}
	ID string
	URL string
}

func statusCode(link string) (int) {
	response, err := http.Head(link)
	if (err!=nil) {
		return 0
	} else {
		return response.StatusCode
	}
}

func sendMaster(url, id string) {

		m := Message{id, url}
		jsonMessage, err := json.Marshal(m)
		if (err!=nil) {
			log.Fatal(err)
		}
		fmt.Println(string(jsonMessage))

    	client := &http.Client{}
    	resp, err := client.Post("http://localhost:4005", "application/json", strings.NewReader(string(jsonMessage)))
    	if err != nil {
       		panic(err)
    	}
    	defer resp.Body.Close()
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if (r.URL.Path=="/"){
			http.Redirect(w,r,"/form.html",301)
		} else {
			mimetype:=mime.TypeByExtension(filepath.Ext(r.URL.Path[1:]))
			w.Header().Set("Content-type", mimetype)
			t, err := template.ParseFiles(r.URL.Path[1:]) 
			if (err!=nil){
				log.Fatal(err)
			} else {
				t.Execute(w, nil)
			}	
		}
	} else {
		URL:=r.FormValue("url")
		rb_ID:=r.FormValue("rb-id")
		fmt.Println(URL,rb_ID)
		fmt.Println(statusCode(URL))

		m := Message{rb_ID, URL}
		jsonMessage, err := json.Marshal(m)
		if (err!=nil) {
			log.Fatal(err)
		}
		fmt.Println(string(jsonMessage))

		//req, err := http.NewRequest("POST", "http://10.0.0.114:5000", bytes.NewBuffer(jsonMessage))
		// req, err := http.NewRequest("POST", "http://localhost:4005", bytes.NewBuffer(jsonMessage))
		// if err != nil {
  //      		panic(err)
  //   	}
  //   	// req.Header.Set("X-Custom-Header", "myvalue")
  //   	req.Header.Set("Content-Type", "application/json")

    	client := &http.Client{}
    	resp, err := client.Post("http://localhost:4005", "application/json", strings.NewReader(string(jsonMessage)))
    	if err != nil {
       		panic(err)
    	}
    	defer resp.Body.Close()

		//http.Post("http://10.0.0.114:5000", string(b))
		http.Redirect(w,r,"/form.html",301)
	}
	
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method=="POST" {
		fmt.Println("Form submitted!!")
		URL:=r.FormValue("url")
		rb_ID:=r.FormValue("rb-id")
		fmt.Println(URL,rb_ID)
		fmt.Println(statusCode(URL))
		sendMaster(URL,rb_ID)
	}
}

func main() {

	http.HandleFunc("/",formHandler)
	http.HandleFunc("/form-submit",submitHandler)

	http.ListenAndServe("localhost:4003", nil)
}
