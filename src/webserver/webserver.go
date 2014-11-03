package main

import (
	"strings"
	"net/http"
	"log"
	"mime"
	"path/filepath"
	"html/template"
	"encoding/json"
	"strconv"
)

type Message struct {
	ID  string
	URL string
}

type Reply struct {
	Code string
	URL  string
	ID   string
}

func statusCode(link string) (int) {
	response, err := http.Head(link)
	if (err != nil) {
		return 0
	} else {
		return response.StatusCode
	}
}

func sendMaster(url, id string) {
	m := Message{id, url}
	json_message, err := json.Marshal(m)
	if (err != nil) {
		log.Fatal(err)
	}

	client := &http.Client{}
	response, err := client.Post("http://localhost:5000", "application/json", strings.NewReader(string(json_message)))
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
}

func reply(URL, status_code, slave_ID string) ([]byte) {
	reply := Reply{status_code, URL, slave_ID}
	json_message, err := json.Marshal(reply)
	if err != nil {
		log.Fatal(err)
	}
	return json_message
}

func sendInfo(response_writer http.ResponseWriter, status_code string, URL string, slave_ID string) {
	reply_Message := reply(URL, status_code, slave_ID)
	response_writer.Header().Set("Content-Type", "application/json")
	response_writer.Write(reply_Message)
}

func formHandler(response_writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		if (request.URL.Path == "/") {
			request.URL.Path+="form.html"
		}
		mime_type := mime.TypeByExtension(filepath.Ext(request.URL.Path[1:]))
		response_writer.Header().Set("Content-type", mime_type)
		template, err := template.ParseFiles(request.URL.Path[1:])
		if (err != nil) {
			log.Fatal(err)
		} else {
			template.Execute(response_writer, nil)
		}
	}
}

func submitHandler(response_writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		URL := request.FormValue("url")
		slave_ID := request.FormValue("rb-id")
		status_code := statusCode(URL)
		sendMaster(URL, slave_ID)
		sendInfo(response_writer, strconv.Itoa(status_code), URL, slave_ID)
	}
}

func main() {
	http.HandleFunc("/", formHandler)
	http.HandleFunc("/form-submit", submitHandler)
	http.ListenAndServe("localhost:4003", nil)
}
