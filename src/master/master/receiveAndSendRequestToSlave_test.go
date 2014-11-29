package master

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"encoding/json"
	"strings"
	"io/ioutil"
)

func TestReceiveRequestAndSendToSlave(t *testing.T) {
	testSlaveMap := make(map[string]Slave)
	var receivedUrl string
	testMaster:= httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ReceiveRequestAndSendToSlave(writer, request, testSlaveMap)
		}))

	testSlave:= httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		receivedUrl = request.PostFormValue("url")
		}))
	testSlaveMap["testSlaveName"] = Slave{testSlave.URL, time.Now(), ""}

	m := PostURLRequest{"testSlaveName", "testURL"}
	json_message, _ := json.Marshal(m)
	client := &http.Client{}
	_, err := client.Post(testMaster.URL, "application/json", strings.NewReader(string(json_message)))

	assert.Equal(t, "testURL", receivedUrl)
	assert.Nil(t, err)
}

func TestReceiveRequestAndSendToSlaveWithEmptySlaveAddress(t *testing.T) {
	testSlaveMap := make(map[string]Slave)

	testMaster:= httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ReceiveRequestAndSendToSlave(writer, request, testSlaveMap)
	}))

	testSlaveMap["testSlaveName"] = Slave{"", time.Now(), ""}

	m := PostURLRequest{"testSlaveName", "testURL"}
	json_message, _ := json.Marshal(m)
	client := &http.Client{}
	response, err := client.Post(testMaster.URL, "application/json", strings.NewReader(string(json_message)))
	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	receivedResponse := string(body[:])
	assert.Equal(t, "FAILED to send url to slave. Slave URL is empty for some reason.", receivedResponse)
	assert.Nil(t, err)
}

func TestParseJson(t *testing.T) {
	inputJSON := []byte(`{"DestinationSlaveName": "LeftScreen", "URLToLoadInBrowser": "http://google.com"}`)

	parsedJson, err := parseJson(inputJSON)
	assert.Equal(t, "LeftScreen", parsedJson.DestinationSlaveName)
	assert.Equal(t, "http://google.com", parsedJson.URLToLoadInBrowser)
	assert.Nil(t, err)
}

func TestParseJsonForEmptyInput(t *testing.T) {
	var inputJSON = []byte(`{}`)

	parsedJson, err := parseJson(inputJSON)

	assert.Equal(t, "", parsedJson.DestinationSlaveName)
	assert.Equal(t, "", parsedJson.URLToLoadInBrowser)
	assert.Nil(t, err)
}

func TestParseJsonForNilInput(t *testing.T) {
	_, err := parseJson(nil)

	assert.NotNil(t, err)
}

func TestDestinationAddressSlave(t *testing.T) {
	slaveMap := InitializeTestSlaveMap()
	destinationURL := destinationSlaveAddress("slave1", slaveMap)

	assert.Equal(t, "http://10.0.0.122:8080", destinationURL)
}

func TestDestinationAddressSlaveForEmptySlaveMap(t *testing.T) {
	slaveMap :=  make(map[string]Slave)
	destinationURL := destinationSlaveAddress("slave2", slaveMap)

	assert.Equal(t,"", destinationURL)
}

func TestSendUrlValueMessageToSlave(t *testing.T) {
	var numberOfMessagesSent = 0
	var url = ""

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		numberOfMessagesSent++
		url = request.PostFormValue("url")
	}))

	err := sendUrlValueMessageToSlave(testServer.URL, "http://index.hu")

	assert.Equal(t, 1, numberOfMessagesSent)
	assert.Equal(t, "http://index.hu", url)
	assert.Nil(t,err)
}

func TestSendUrlValueMessageToSlaveSlaveDoesNotRespond(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
	}))
	testServer.Close()
	err := sendUrlValueMessageToSlave(testServer.URL, "http://index.hu")
	assert.NotNil(t,err)
}
