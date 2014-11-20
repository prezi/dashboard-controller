package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strconv"
	"strings"
	"io/ioutil"
	"network"
)

var MASTER_URL = "http://localhost:5000"
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

func main() {
	fs := http.FileServer(http.Dir(STATIC_PATH))
	_ = network.GetUrl("4003")

	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", formHandler)
	http.HandleFunc("/form-submit", submitHandler)
	http.HandleFunc("/receive_slave", receiveAndMapSlaveAddress)
	http.ListenAndServe(":4003", nil)
}

func formHandler(response_writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		if request.URL.Path != "/" {
			http.Redirect(response_writer, request, "/", 301)
		}
		template, err := template.ParseFiles(path.Join(TEMPLATE_PATH, "form.html"))
		if err != nil {
			fmt.Println(err)
		}
		template.Execute(response_writer, id_list)
	}
}

func submitHandler(response_writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		urlToDisplay := request.FormValue("url")
		slave_ID := request.FormValue("slave-id")
		status_code := checkStatusCode(urlToDisplay)
		sendUrlAndIdToMaster(MASTER_URL, urlToDisplay, slave_ID)
		sendConfirmationMessageToUser(response_writer, strconv.Itoa(status_code), urlToDisplay, slave_ID)
	}
}

func checkStatusCode(link string) int {
	response, err := http.Head(link)
	if err != nil {
		return 0
	} else {
		return response.StatusCode
	}
}

func sendUrlAndIdToMaster(masterUrl, urlToDisplay, id string) error {
	m := Message{id, urlToDisplay}
	json_message, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
		return err
	}
	client := &http.Client{}
	response, err := client.Post(masterUrl, "application/json", strings.NewReader(string(json_message)))
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer response.Body.Close()
	return nil
}

func sendConfirmationMessageToUser(response_writer http.ResponseWriter, status_code, URL, slave_ID string) {
	confirmationMessage := confirmationMessage(URL, status_code, slave_ID)
	header := response_writer.Header()
	header.Set("Content-Type", "application/json")
	response_writer.Write(confirmationMessage)
}

func confirmationMessage(URL, status_code, slave_ID string) []byte {
	t, err := template.ParseFiles(path.Join(TEMPLATE_PATH, "infobox.html"))
	if err != nil {
		fmt.Println(err)
	}

	buf := new(bytes.Buffer)
	t.ExecuteTemplate(buf, "T", StatusMessage{status_code, URL, slave_ID})

	jsonMessage, err := json.Marshal(Reply{HTML: buf.String()})
	if err != nil {
		fmt.Println(err)
	}
	return jsonMessage
}

func receiveAndMapSlaveAddress(_ http.ResponseWriter, request *http.Request) {
	POSTRequestBody, _ := ioutil.ReadAll(request.Body)
	defer request.Body.Close()

	err := json.Unmarshal(POSTRequestBody, &id_list)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("\nNEW SLAVE RECEIVED.\n")
	fmt.Println("Slave Name: ", id_list.Id)
}
