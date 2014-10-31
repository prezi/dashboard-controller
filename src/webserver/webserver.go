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
	"strconv"
)

type Message struct {
	ID string
	URL string
}

type Reply struct {
	Code string
	URL string
	ID string
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
    	resp, err := client.Post("http://localhost:5000", "application/json", strings.NewReader(string(jsonMessage)))
    	if err != nil {
       		panic(err)
    	}
    	defer resp.Body.Close()
}

func reply(URL,status_code, rb_ID string) ([]byte) {

	r:=Reply{status_code,URL,rb_ID}
	jsonMessage, err := json.Marshal(r)
	if err!=nil {
		log.Fatal(err)
	}
	return jsonMessage
}

func sendInfo(w http.ResponseWriter,status_code string,URL string,rb_ID string) {
	replyMessage:=reply(URL,status_code,rb_ID)
	w.Header().Set("Content-Type", "application/json")
	w.Write(replyMessage)	
}

func formHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		if (r.URL.Path=="/"){
			r.URL.Path+="form.html"
		}
		mimetype:=mime.TypeByExtension(filepath.Ext(r.URL.Path[1:]))
		w.Header().Set("Content-type", mimetype)
		t, err := template.ParseFiles(r.URL.Path[1:]) 
		if (err!=nil){
			log.Fatal(err)
		} else {
			t.Execute(w, nil)
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
    	client := &http.Client{}
    	resp, err := client.Post("http://localhost:5000", "application/json", strings.NewReader(string(jsonMessage)))
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
		URL:=r.FormValue("url")
		rb_ID:=r.FormValue("rb-id")
		status_code:=statusCode(URL)
		sendMaster(URL,rb_ID)
		sendInfo(w,strconv.Itoa(status_code),URL,rb_ID)
	}
}

func main() {
	http.HandleFunc("/",formHandler)
	http.HandleFunc("/form-submit",submitHandler)
	http.ListenAndServe("localhost:4003", nil)
}
