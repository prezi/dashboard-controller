package website

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
	"io/ioutil"
	"master/master/delegateRequestToSlave"
	"network"
)

type PostURLRequest struct {
	DestinationSlaveName string
	URLToLoadInBrowser   string
}

func TestIndexPageHandler(t * testing.T) {
	VIEWS_PATH = network.GetRelativeFilePath("views")
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		IndexPageHandler(w, request)
	}))
	client := &http.Client{}
	resp, _ := client.Get(testServer.URL)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "text/html; charset=utf-8", resp.Header.Get("Content-type"))
}

func TestIndexPageHandlerWithWrongPath(t * testing.T) {
	VIEWS_PATH = network.GetRelativeFilePath("dummy")
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		IndexPageHandler(w, request)
	}))
	client := &http.Client{}
	resp, _ := client.Get(testServer.URL)
	assert.Equal(t, 404, resp.StatusCode)
}

func sendGetToFormHandler(URL string) int {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		request.URL.Path = URL
		testSlaveNames := []string {"a","b"}
		FormHandler(w, request, testSlaveNames)
	}))

	client := &http.Client{}
	resp, _ := client.Get(testServer.URL)

	return resp.StatusCode
}

func TestFormHandler(t *testing.T) {
	assert.Equal(t, 302, sendGetToFormHandler("/"))
}

func TestStatusMessageForAvailableServer(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
	}))
	statusMessage := delegateRequestToSlave.IsURLValid(testServer.URL)

	assert.Equal(t, true, statusMessage)
}

func TestStatusMessageForUnavailableServer(t *testing.T) {
	statusMessage := delegateRequestToSlave.IsURLValid("")

	assert.Equal(t, false, statusMessage)
}

func TestSendConfirmationMessageToUser(t *testing.T) {
	VIEWS_PATH = network.GetRelativeFilePath("views")
	testMessage := "testmessage"
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		sendConfirmationMessageToUser(w, testMessage)
	}))

	client := &http.Client{}
	resp, _ := client.Get(testServer.URL)
	respBodyContents, _:= ioutil.ReadAll(resp.Body)
	respBodyString := string(respBodyContents[:])
	assert.True(t, strings.Contains(respBodyString, testMessage))
	assert.Equal(t,"application/json", resp.Header.Get("Content-type"))
}

func TestCreateConfirmationMessage(t *testing.T) {
	VIEWS_PATH = network.GetRelativeFilePath("views")
	msg := "testmessage"
	confirmationMessageJson, _ := createConfirmationMessage(msg)
	confirmationMessageJsonString := string(confirmationMessageJson[:])
	assert.True(t, strings.Contains(confirmationMessageJsonString, msg))
}

func TestCreateConfirmationMessageWithWrongPath(t *testing.T) {
	VIEWS_PATH = network.GetRelativeFilePath("dummy")
	msg := "testmessage"
	confirmationMessageJson, _ := createConfirmationMessage(msg)
	assert.Equal(t, len(confirmationMessageJson), 0)
}
