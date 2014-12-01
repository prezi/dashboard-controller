package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"network"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var MASTER_URL = "http://localhost:5000"
var WEBSERVER_PORT = "4003"
var VIEWS_PATH = "src/webserver/views/"

const (
	DEFAULT_MASTER_IP_ADDRESS = "localhost"
	DEFAULT_MASTER_PORT       = "5000"
	DEFAULT_WEBSERVER_PORT    = "4003"
	IMAGES_PATH = "src/webserver/assets/images"
	JAVASCRIPTS_PATH = "src/webserver/assets/javascripts"
	STYLESHEETS_PATH = "src/webserver/assets/stylesheets"
)

type Message struct {
	DestinationSlaveName string
	URLToLoadInBrowser   string
}

type StatusMessage struct {
	Code       string
	URL        string
	ID         string
	SlaveError string
}

type Reply struct {
	HTML string
}

type IdList struct {
	Id []string
}

var id_list = IdList{
	Id: []string{},
}

func main() {
	MASTER_URL = setMasterAddress()
	http.Handle("/assets/images/", http.StripPrefix("/assets/images/", http.FileServer(http.Dir(IMAGES_PATH))))
	http.Handle("/assets/javascripts/", http.StripPrefix("/assets/javascripts/", http.FileServer(http.Dir(JAVASCRIPTS_PATH))))
	http.Handle("/assets/stylesheets/", http.StripPrefix("/assets/stylesheets/", http.FileServer(http.Dir(STYLESHEETS_PATH))))
	http.HandleFunc("/", formHandler)
	http.HandleFunc("/form-submit", submitHandler)
	http.HandleFunc("/receive_slave", receiveAndMapSlaveAddress)
	go sendInitToMaster(MASTER_URL, "/webserver_init")
	go startWebserverHeartbeats(5, MASTER_URL, "/webserver_heartbeat")
	http.ListenAndServe(":"+WEBSERVER_PORT, nil)
}

func setMasterAddress() (masterUrl string) {
	masterIP, masterPort := configFlags()
	masterUrl = network.AddProtocolAndPortToIP(masterIP, masterPort)
	return
}

func configFlags() (masterIP, masterPort string) {
	flag.StringVar(&masterIP, "masterIP", DEFAULT_MASTER_IP_ADDRESS, "master IP address")
	flag.StringVar(&masterPort, "masterPort", DEFAULT_MASTER_PORT, "master port number")
	flag.StringVar(&WEBSERVER_PORT, "webserverPort", DEFAULT_WEBSERVER_PORT, "webserver port number")
	flag.Parse()
	return masterIP, masterPort
}

func formHandler(response_writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		if request.URL.Path != "/" {
			http.Redirect(response_writer, request, "/", 301)
		}
		template, err := template.ParseFiles(path.Join(VIEWS_PATH, "form.html"))
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
		slaveError := ""
		url := request.FormValue("url")
		name := request.FormValue("slave-id")
		if slaveInSlaveList(name, id_list.Id) == false {
			slaveError = "This slave does not exist, please refresh your browser."
		}
		sendConfirmationMessageToUser(response_writer, returnStatusMessageFrom(url), url, name, slaveError)
		if isUrlValid(url) == true {
			sendUrlAndIdToMaster(MASTER_URL, url, name)
		}
	}
}

func slaveInSlaveList(slaveName string, slaveIdList []string) bool {
	for _, slaveId := range slaveIdList {
		if slaveId == slaveName {
			return true
		}
	}
	return false
}

func sendConfirmationMessageToUser(response_writer http.ResponseWriter, status_code, URL, slave_ID, slaveError string) {
	confirmationMessage := confirmationMessage(URL, status_code, slave_ID, slaveError)

	header := response_writer.Header()
	header.Set("Content-Type", "application/json")
	response_writer.Write(confirmationMessage)
}

func confirmationMessage(URL, status_code, slave_ID, slaveError string) []byte {
	t, err := template.ParseFiles(path.Join(VIEWS_PATH, "infobox.html"))
	if err != nil {
		fmt.Println(err)
	}

	buf := new(bytes.Buffer)
	t.ExecuteTemplate(buf, "T", StatusMessage{status_code, URL, slave_ID, slaveError})

	jsonMessage, err := json.Marshal(Reply{HTML: buf.String()})
	if err != nil {
		fmt.Println(err)
	}
	return jsonMessage
}

func returnStatusMessageFrom(url string) (statusMessage string) {
	statusCode := checkStatusCode(url)
	if 400 <= statusCode || statusCode == 0 {
		statusMessage = "URL cannot be open :( (HTTP status code " + strconv.Itoa(statusCode) + ")"
	} else {
		statusMessage = "Success!"
	}
	return
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

func isUrlValid(url string) bool {
	if 400 <= checkStatusCode(url) || checkStatusCode(url) == 0 {
		return false
	} else {
		return true
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

func receiveAndMapSlaveAddress(_ http.ResponseWriter, request *http.Request) {
	POSTRequestBody, _ := ioutil.ReadAll(request.Body)
	defer request.Body.Close()

	err := json.Unmarshal(POSTRequestBody, &id_list)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("\nSLAVE LIST UPDATED.\n")
	fmt.Println("Slave Name: ", id_list.Id)
}

func sendInitToMaster(masterUrl, pattern string) {
	postRequestUrl := masterUrl
	postRequestUrl += pattern
	client := &http.Client{}
	form := url.Values{}
	form.Set("message", "update me!")
	form.Set("webserverPort", WEBSERVER_PORT)
	client.PostForm(postRequestUrl, form)
}

func startWebserverHeartbeats(heartbeatInterval int, masterUrl, pattern string) {
	var err error
	postRequestUrl := masterUrl
	postRequestUrl += pattern
	client := &http.Client{}
	form := url.Values{}
	form.Set("webserverPort", WEBSERVER_PORT)
	beat := time.Tick(time.Duration(heartbeatInterval) * time.Second)
	for _ = range beat {
		_, err = client.PostForm(postRequestUrl, form)
		network.ErrorHandler(err, "Error communicating with master: %v\n")
	}
}
