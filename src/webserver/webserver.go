package main

import (
	"strings"
	"net/http"
	"log"
	"mime"
	"path"
	"path/filepath"
	"html/template"
	"encoding/json"
	"strconv"
	"bytes"
	"fmt"
)

var MASTER_URL ="http://localhost:5000"
var TEMPLATE_PATH = "src/webserver/templates/"
var STATIC_PATH = "src/webserver/static"

type Message struct {
	ID  string
	URL string
}

type StatusMessage struct {
	Code string
	URL  string
	ID   string
}

type Reply struct {
	HTML string
}

type IdList struct {
	Id []string
}

var id_list = IdList{
	Id: []string{"1", "2"},
}

func statusCode(link string) (int) {
	response, err := http.Head(link)
	if (err != nil) {
		return 0
	} else {
		return response.StatusCode
	}
}

func sendMaster(masterUrl,urlToDisplay, id string) {
	m := Message{id, urlToDisplay}
	json_message, err := json.Marshal(m)
	if (err != nil) {
		log.Fatal(err)
	}

	client := &http.Client{}
	response, err := client.Post(masterUrl, "application/json", strings.NewReader(string(json_message)))
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
}

func reply(URL, status_code, slave_ID string) ([]byte) {
	// reply := Reply{status_code, URL, slave_ID}
	
	t, err := template.ParseFiles(path.Join(TEMPLATE_PATH,"infobox.html"))		
	if (err != nil) {
			log.Fatal(err)
	} 

	buf := new(bytes.Buffer)
	t.ExecuteTemplate(buf, "T", StatusMessage{status_code, URL, slave_ID})
	
	json_message, err := json.Marshal(Reply{HTML:buf.String()})
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

func setMimeType(responseWriter http.ResponseWriter, path string) {
	mime_type := mime.TypeByExtension(filepath.Ext(path))
	responseWriter.Header().Set("Content-type", mime_type)
}

func formHandler(response_writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		template, err := template.ParseFiles(path.Join(TEMPLATE_PATH,"form.html"))
		if (err != nil) {
    		http.Error(response_writer, http.StatusText(500), 500)
			log.Fatal(err)
		} 
		template.Execute(response_writer, id_list)
	}
}

func submitHandler(response_writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		urlToDisplay := request.FormValue("url")
		slave_ID := request.FormValue("rb-id")
		status_code := statusCode(urlToDisplay)
	    sendMaster(MASTER_URL,urlToDisplay, slave_ID)
		sendInfo(response_writer, strconv.Itoa(status_code), urlToDisplay, slave_ID)
	}
}

func receiveAndMapSlaveAddress(_ http.ResponseWriter, request *http.Request) {
	slaveName := request.PostFormValue("slaveName")
	fmt.Printf("\nNEW SLAVE RECEIVED.\n")
	fmt.Println("Slave Name: ", slaveName)
	id_list.Id = append(id_list.Id, slaveName)
}

func main() {
	fs := http.FileServer(http.Dir(STATIC_PATH))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", formHandler)
	http.HandleFunc("/form-submit", submitHandler)
	http.HandleFunc("/receive_slave", receiveAndMapSlaveAddress)
	http.ListenAndServe("localhost:4003", nil)
}
