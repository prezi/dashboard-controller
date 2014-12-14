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
	"website/session"
	"io"
)

var (
	IMAGES_PATH      = master.GetRelativeFilePath("assets/images")
	JAVASCRIPTS_PATH = master.GetRelativeFilePath("assets/javascripts")
	STYLESHEETS_PATH = master.GetRelativeFilePath("assets/stylesheets")
	VIEWS_PATH       = master.GetRelativeFilePath("views")
)

type StatusMessage struct {
	StatusMessage string
}

func IndexPageHandler(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles(path.Join(VIEWS_PATH, "login.html"))
	if network.ErrorHandler(err, "Error encountered while parsing website template files: %v.") {
		w.WriteHeader(404)
		return
	}
	template.Execute(w, "Login Error Message Here.")
}

func FormHandler(w http.ResponseWriter, r *http.Request, slaveNames []string) {
	if r.Method == "GET" {
		userName := session.GetUsername(r)
		if userName != "" {
			displayFormPage(w, slaveNames, userName)
		} else {
			http.Redirect(w, r, "/", 302)
		}
	}
}

func displayFormPage(responseWriter http.ResponseWriter, slaveNames []string, userName string) {
	type HTMLData struct {
		UserName   string
		SlaveNames []string
	}

	template, err := template.ParseFiles(path.Join(VIEWS_PATH, "form.html"))
	network.ErrorHandler(err, "Error encountered while parsing website template files: %v.")
	dataForTemplate := HTMLData{UserName: userName, SlaveNames: slaveNames}

	template.Execute(responseWriter, dataForTemplate)
}

func SubmitHandler(response_writer http.ResponseWriter, request *http.Request, slaveMap map[string]master.Slave) {
	if request.Method == "POST" {
		URLToDisplay, slaveNamesToUpdate, err := parseFromJSON(request.Body)
		if err != nil {
			fmt.Println("Error parsing JSON request discarded")
			sendConfirmationMessageToUser(response_writer, "Failed to parse JSON ")
			return
		}
		if !isURLValid(URLToDisplay) {
			BadURLStatusMessage := "Sorry, " + URLToDisplay + " cannot be opened. Try a different one. Sadpanda."
			sendConfirmationMessageToUser(response_writer, BadURLStatusMessage)
			return
		}
		if nonExistentSlave := allSlavesAreConnected(slaveMap, slaveNamesToUpdate); nonExistentSlave != ""{
			errorMessage := "Sorry, " + nonExistentSlave + ` cannot be reached. Please refresh the page
			 to see an updated list.`
			sendConfirmationMessageToUser(response_writer, errorMessage)
			return
		}
		sendURLToSlaves(slaveMap, slaveNamesToUpdate, URLToDisplay)
		statusMessage := "Slaves are updated"
		sendConfirmationMessageToUser(response_writer, statusMessage)
	}
}

func parseFromJSON(requestBody io.ReadCloser) (URLToDisplay string, slaveNames []string, err error) {
	type FormData struct {
		URLToDisplay   string
		SlaveNames []string
	}
	JSONFormData := json.NewDecoder(requestBody)
	var decodedFormData FormData
	err = JSONFormData.Decode(&decodedFormData)
	URLToDisplay = decodedFormData.URLToDisplay
	slaveNames = decodedFormData.SlaveNames
	return
}

func sendURLToSlaves(slaveMap map[string]master.Slave, slaveNames []string, URLToDisplay string) {
	for _, slaveName := range slaveNames {
		receiveAndSendRequestToSlave.ReceiveRequestAndSendToSlave(slaveMap, slaveName, URLToDisplay)
	}
}

func sendConfirmationMessageToUser(response_writer http.ResponseWriter, statusMessage string) {
	confirmationMessage, _ := createConfirmationMessage(statusMessage)

	header := response_writer.Header()
	header.Set("Content-Type", "application/json")
	response_writer.Write(confirmationMessage)
}

func createConfirmationMessage(statusMessage string) (jsonMessage []byte, err error) {
	t, err := template.ParseFiles(path.Join(VIEWS_PATH, "infobox.html"))
	if err != nil {
		fmt.Println(err)
		return
	}

	buf := new(bytes.Buffer)
	t.ExecuteTemplate(buf, "T", StatusMessage{statusMessage})

	jsonMessage, err = json.Marshal(StatusMessage{StatusMessage: buf.String()})
	if err != nil {
		fmt.Println(err)
	}
	return
}

func checkStatusCode(urlToDisplay string) int {
	if (len(urlToDisplay) <= 6) {
		urlToDisplay = "http://" + urlToDisplay
	} else if (string(urlToDisplay[0:6]) != "http:/" && string(urlToDisplay[0:6]) != "https:") {
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

func allSlavesAreConnected(slaveMap map[string]master.Slave, slaveNamesToUpdate []string) (nonExistentSlave string) {
	for _, nonExistentSlave = range slaveNamesToUpdate {
		if _, isExists := slaveMap[nonExistentSlave]; !isExists {
			return nonExistentSlave
		}
	}
	return ""
}
