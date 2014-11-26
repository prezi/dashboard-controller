package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"
	"io/ioutil"
	"errors"
	"net/url"
	"os"
	"network"
	"flag"
	"strconv"
)

var MASTER_URL = "http://localhost:5000"
var TEMPLATE_PATH = "src/webserver/templates/"
var STATIC_PATH = "src/webserver/static"

const (
	DEFAULT_MASTER_IP_ADDRESS = "localhost"
	DEFAULT_MASTER_PORT = "5000"
	DEFAULT_WEBSERVER_PORT = "4003"
)

type Message struct {
	DestinationSlaveName string
	URLToLoadInBrowser string
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
	Id: []string{"slave1", "slave2"},
}

func main() {
	fs := http.FileServer(http.Dir(STATIC_PATH))
	MASTER_URL = setMasterAddress()
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", formHandler)
	http.HandleFunc("/form-submit", submitHandler)
	http.HandleFunc("/receive_slave", receiveAndMapSlaveAddress)
	go requestSlaveIdsOnStart(MASTER_URL,"/webserver_init")
	http.ListenAndServe(":" + DEFAULT_WEBSERVER_PORT, nil)
}

func setMasterAddress() (masterUrl string) {
	masterIP, masterPort := configFlags()
	masterUrl = network.AddProtocolAndPortToIP(masterIP, masterPort)
	return
}

func configFlags() (masterIP, masterPort string) {
	flag.StringVar(&masterIP, "masterIP", DEFAULT_MASTER_IP_ADDRESS, "master IP address")
	flag.StringVar(&masterPort, "masterPort", DEFAULT_MASTER_PORT, "master port number")
	flag.Parse()
	return masterIP, masterPort
}

func requestSlaveIdsOnStart(masterUrl,pattern string) (err error) {
	err = nil
	postRequestUrl := masterUrl
	postRequestUrl += pattern
	client := &http.Client{}
	form := url.Values{}
	form.Set("message","send_me_the_list")
	form.Set("webserverPort", DEFAULT_WEBSERVER_PORT)
	resp, err := client.PostForm(postRequestUrl,form)
	if err != nil {
		fmt.Println(err)
		return
	}
	if resp.StatusCode != 200 {
		err = errors.New("Master is not available")
	}
	return
}

func formHandler(response_writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		if request.URL.Path != "/" {
			http.Redirect(response_writer, request, "/", 301)
		}
		template, err := template.ParseFiles(path.Join(TEMPLATE_PATH, "form.html"))
		if err != nil {
			fmt.Println("Html files not found. Please restart from the root folder.")
			os.Exit(1)
		} else {
			template.Execute(response_writer, id_list)
		}
	}
}

func submitHandler(response_writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		urlToDisplay := request.FormValue("url")
		slave_ID := request.FormValue("slave-id")
		status_code := checkStatusCode(urlToDisplay)
		statusMessage := ""
		if 400 <= status_code || status_code == 0 {
			statusMessage = "URL cannot be open :( (HTTP status code " + strconv.Itoa(status_code) + ")" 
		} else {
			sendUrlAndIdToMaster(MASTER_URL, urlToDisplay, slave_ID)
			statusMessage = "Success!" 
		}
		sendConfirmationMessageToUser(response_writer, statusMessage, urlToDisplay, slave_ID)
	}
}

func checkStatusCode(urlToDisplay string) int {
	if len(urlToDisplay) < 4 || string(urlToDisplay[0:4]) != "http" {
		urlToDisplay = "http://" + urlToDisplay
	}

    response, err := http.Head(urlToDisplay)
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
