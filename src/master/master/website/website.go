package website

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"master/master"
	"master/master/receiveAndSendRequestToSlave"
	"net/http"
	"network"
	"path"
	"runtime"
)

var (
	IMAGES_PATH      = getRelativeFilePath("assets/images")
	JAVASCRIPTS_PATH = getRelativeFilePath("assets/javascripts")
	STYLESHEETS_PATH = getRelativeFilePath("assets/stylesheets")
	VIEWS_PATH       = getRelativeFilePath("views")
)

type StatusMessage struct {
	StatusMessage string
}

func getRelativeFilePath(relativeFileName string) (filePath string) {
	_, filename, _, _ := runtime.Caller(1)
	filePath = path.Join(path.Dir(filename), relativeFileName)
	return
}

func LoginHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		template, err := template.ParseFiles(path.Join(VIEWS_PATH, "login.html"))
		network.ErrorHandler(err, "Error encountered while parsing website template files: %v.")
		template.Execute(responseWriter, "future errormessage when login failed")
	}
	if request.Method == "POST" {
		username := request.FormValue("username")
		password := request.FormValue("password")
		fmt.Println(username, password)
		http.Redirect(responseWriter, request, "/", 301)
	}
}

func FormHandler(responseWriter http.ResponseWriter, request *http.Request, slaveNames []string) {
	if request.Method == "GET" {
		// if request.URL.Path != "/login" && request.URL.Path != "/" {
		// 	http.Redirect(responseWriter, request, "/login", 301)
		// }
		fmt.Println("EXECUTING TEMPLATE YO")
		parseAndExecuteTemplate(responseWriter, slaveNames)
	}
}

// func internalPageHandler(w http.ResponseWriter, r *http.Request) {
// 	userName := getUserName(r)
// 	if userName != "" {
// 		fmt.Fprintf(w, internalPage, userName)
// 	} else {
// 		http.Redirect(w, r, "/", 302)
// 	}
// }

func parseAndExecuteTemplate(responseWriter http.ResponseWriter, slaveNames []string) {
	type SlaveNameList struct {
		SlaveNames []string
	}

	template, err := template.ParseFiles(path.Join(VIEWS_PATH, "form.html"))
	network.ErrorHandler(err, "Error encountered while parsing website template files: %v.")
	slaveNameList := SlaveNameList{SlaveNames: slaveNames}
	template.Execute(responseWriter, slaveNameList)
}

func SubmitHandler(response_writer http.ResponseWriter, request *http.Request, slaveMap map[string]master.Slave) {
	if request.Method == "POST" {
		url := request.FormValue("url")
		slaveName := request.FormValue("slave-id")

		if slaveExistsInSlaveMap(slaveName, slaveMap) == false {
			statusMessage := slaveName + " is offline. Please refresh your browser to see available destinations screens."
			sendConfirmationMessageToUser(response_writer, statusMessage)
		} else {
			if isURLValid(url) == true {
				statusMessage := "Success! " + url + " is being displayed on " + slaveName + "."
				sendConfirmationMessageToUser(response_writer, statusMessage)
				receiveAndSendRequestToSlave.ReceiveRequestAndSendToSlave(slaveMap, slaveName, url)
			} else {
				statusMessage := "Sorry, " + url + " cannot be opened. Try a different one. Sadpanda."
				sendConfirmationMessageToUser(response_writer, statusMessage)
			}
		}
	}
}

func slaveExistsInSlaveMap(slaveName string, slaveMap map[string]master.Slave) bool {
	for slaveNameInMap, _ := range slaveMap {
		if slaveName == slaveNameInMap {
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

func isURLValid(url string) bool {
	if 400 <= checkStatusCode(url) || checkStatusCode(url) == 0 {
		return false
	} else {
		return true
	}
}
