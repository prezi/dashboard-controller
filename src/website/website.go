package website

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"master/master"
	"master/master/delegateRequestToSlave"
	"net/http"
	"network"
	"path"
	"website/session"
)

var (
	IMAGES_PATH      = network.PROJECT_ROOT + "/src/website/assets/images"
	JAVASCRIPTS_PATH = network.PROJECT_ROOT + "/src/website/assets/javascripts"
	STYLESHEETS_PATH = network.PROJECT_ROOT + "/src/website/assets/stylesheets"
	VIEWS_PATH       = network.PROJECT_ROOT + "/src/website/views"
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
		if !master.IsURLValid(URLToDisplay) {
			BadURLStatusMessage := "Sorry, " + URLToDisplay + " cannot be opened. Try a different one. Sadpanda."
			sendConfirmationMessageToUser(response_writer, BadURLStatusMessage)
			return
		}
		if nonExistentSlaves := delegateRequestToSlave.CheckIfRequestedSlavesAreConnected(slaveMap, slaveNamesToUpdate); nonExistentSlaves != "" {
			errorMessage := "Sorry, " + nonExistentSlaves + ` cannot be reached. Please refresh the page
			 to see an updated list.`
			sendConfirmationMessageToUser(response_writer, errorMessage)
			return
		}
		delegateRequestToSlave.ReceiveRequestAndSendToSlave(slaveMap, slaveNamesToUpdate, URLToDisplay)
		statusMessage := "Slaves are updated"
		sendConfirmationMessageToUser(response_writer, statusMessage)
	}
}

func parseFromJSON(requestBody io.ReadCloser) (URLToDisplay string, slaveNames []string, err error) {
	type FormData struct {
		URLToDisplay string
		SlaveNames   []string
	}
	JSONFormData := json.NewDecoder(requestBody)
	var decodedFormData FormData
	err = JSONFormData.Decode(&decodedFormData)
	URLToDisplay = decodedFormData.URLToDisplay
	slaveNames = decodedFormData.SlaveNames
	return
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
