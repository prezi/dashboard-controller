package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"network"
	"os"
	"path"
	"strings"
	"time"
)

const (
	DEFAULT_MASTER_URL        = "http://localhost:5000"
	DEFAULT_MASTER_IP_ADDRESS = "localhost"
	DEFAULT_MASTER_PORT       = "5000"
	DEFAULT_WEBSERVER_PORT    = "4003"
	IMAGES_PATH               = "src/webserver/assets/images"
	JAVASCRIPTS_PATH          = "src/webserver/assets/javascripts"
	STYLESHEETS_PATH          = "src/webserver/assets/stylesheets"
	VIEWS_PATH                = "src/webserver/views/"
)

type Message struct {
	DestinationSlaveName string
	URLToLoadInBrowser   string
}

type StatusMessage struct {
	StatusMessage string
}

type IdList struct {
	Id []string
}

var id_list = IdList{
	Id: []string{},
}

func main() {
	isMasterAlive := true
	masterIP, masterPort, webserverPort := configFlags()
	masterUrl := network.AddProtocolAndPortToIP(masterIP, masterPort)
	http.Handle("/assets/images/", http.StripPrefix("/assets/images/", http.FileServer(http.Dir(IMAGES_PATH))))
	http.Handle("/assets/javascripts/", http.StripPrefix("/assets/javascripts/", http.FileServer(http.Dir(JAVASCRIPTS_PATH))))
	http.Handle("/assets/stylesheets/", http.StripPrefix("/assets/stylesheets/", http.FileServer(http.Dir(STYLESHEETS_PATH))))
	http.HandleFunc("/", formHandler)
	http.HandleFunc("/form-submit", func(w http.ResponseWriter, r *http.Request) {
		submitHandler(w, r, isMasterAlive)
	})
	http.HandleFunc("/receive_slave", receiveAndMapSlaveAddress)
	go sendInitToMaster(masterUrl, webserverPort, "/webserver_init")
	go startWebserverHeartbeats(5, masterUrl, webserverPort, "/webserver_heartbeat", &isMasterAlive)
	http.ListenAndServe(":"+webserverPort, nil)
}

func configFlags() (masterIP, masterPort, webserverPort string) {
	flag.StringVar(&masterIP, "masterIP", DEFAULT_MASTER_IP_ADDRESS, "master IP address")
	flag.StringVar(&masterPort, "masterPort", DEFAULT_MASTER_PORT, "master port number")
	flag.StringVar(&webserverPort, "webserverPort", DEFAULT_WEBSERVER_PORT, "webserver port number")
	flag.Parse()
	return masterIP, masterPort, webserverPort
}

func formHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		if request.URL.Path != "/" {
			http.Redirect(responseWriter, request, "/", 301)
		}
		parseAndExecuteTemplate(responseWriter)
	}
}

func parseAndExecuteTemplate(responseWriter http.ResponseWriter) {
	template, err := template.ParseFiles(path.Join(VIEWS_PATH, "form.html"))
	handleTemplateParseError(err)
	template.Execute(responseWriter, id_list)
}

func handleTemplateParseError(err error) {
	if err != nil {
		fmt.Println("Html files not found. Please restart from the root folder.")
		os.Exit(1)
	}
}

func submitHandler(response_writer http.ResponseWriter, request *http.Request, isMasterAlive bool) {
	if request.Method == "POST" {
		url := request.FormValue("url")
		name := request.FormValue("slave-id")
		if isMasterAlive == false {
			statusMessage := "Master offline."
			sendConfirmationMessageToUser(response_writer, statusMessage)
		} else if slaveInSlaveList(name, id_list.Id) == false {
			statusMessage := name + " is offline, please refresh your browser to see available screens."
			sendConfirmationMessageToUser(response_writer, statusMessage)
		} else {
			if isUrlValid(url) == true {
				statusMessage := "Success! " + url + " is being displayed on " + name
				sendUrlAndIdToMaster(DEFAULT_MASTER_URL, url, name)
				sendConfirmationMessageToUser(response_writer, statusMessage)
			} else {
				statusMessage := url + " cannot be opened. Try a different one. Sadpanda."
				sendConfirmationMessageToUser(response_writer, statusMessage)
			}
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

func sendConfirmationMessageToUser(response_writer http.ResponseWriter, statusMessage string) {
	confirmationMessage := createConfirmationMessage(statusMessage)

	header := response_writer.Header()
	header.Set("Content-Type", "application/json")
	response_writer.Write(confirmationMessage)
}

func createConfirmationMessage(statusMessage string) []byte {
	t, err := template.ParseFiles(path.Join(VIEWS_PATH, "infobox.html"))
	if err != nil {
		fmt.Println(err)
	}

	buf := new(bytes.Buffer)
	t.ExecuteTemplate(buf, "T", StatusMessage{statusMessage})

	jsonMessage, err := json.Marshal(StatusMessage{StatusMessage: buf.String()})
	if err != nil {
		fmt.Println(err)
	}
	return jsonMessage
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
	fmt.Println("Slave Names: ", id_list.Id)
}

func sendInitToMaster(masterUrl, webserverPort, pattern string) {
	postRequestUrl := masterUrl
	postRequestUrl += pattern
	client := &http.Client{}
	form := network.CreateFormWithInitialValues(map[string]string{"message": "update me!", "webserverPort": webserverPort})
	client.PostForm(postRequestUrl, form)
}

func startWebserverHeartbeats(heartbeatInterval int, masterUrl, webserverPort, pattern string, isMasterAlive *bool) {
	postRequestUrl := masterUrl
	postRequestUrl += pattern
	client := &http.Client{}
	form := network.CreateFormWithInitialValues(map[string]string{"webserverPort": webserverPort})
	beat := time.Tick(time.Duration(heartbeatInterval) * time.Second)
	for _ = range beat {
		_, err := client.PostForm(postRequestUrl, form)
		if err != nil {
			id_list.Id = id_list.Id[:0]
			fmt.Printf("Error communicating with master: %v\n", err)
			*isMasterAlive = false
		} else {
			*isMasterAlive = true
		}
	}
}
