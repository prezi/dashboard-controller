package website

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
)

const (
	IMAGES_PATH               = "src/master/master/website/assets/images"
	JAVASCRIPTS_PATH          = "src/master/master/website/assets/javascripts"
	STYLESHEETS_PATH          = "src/master/master/website/assets/stylesheets"
	VIEWS_PATH                = "src/master/master/website/views/"
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


func FormHandler(responseWriter http.ResponseWriter, request *http.Request, slaveNames []string) {
	if request.Method == "GET" {
		if request.URL.Path != "/" {
			http.Redirect(responseWriter, request, "/", 301)
		}
		parseAndExecuteTemplate(responseWriter, slaveNames)
	}
}

func parseAndExecuteTemplate(responseWriter http.ResponseWriter, slaveNames []string) {
	template, err := template.ParseFiles(path.Join(VIEWS_PATH, "form.html"))
	handleTemplateParseError(err)
	id_list := IdList{Id: slaveNames}
	template.Execute(responseWriter, id_list)
}

func handleTemplateParseError(err error) {
	if err != nil {
		fmt.Println("Html files not found. Please restart from the root folder.")
		os.Exit(1)
	}
}

func SubmitHandler(response_writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		url := request.FormValue("url")
		name := request.FormValue("slave-id")
		if slaveInSlaveList(name, id_list.Id) == false {
			statusMessage := name + " is offline, please refresh your browser to see available screens."
			sendConfirmationMessageToUser(response_writer, statusMessage)
		} else {
			if isUrlValid(url) == true {
				statusMessage := "Success! " + url + " is being displayed on " + name
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
