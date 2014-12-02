package receiveAndSendRequestToSlave

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"master/master"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func InitializeTestSlaveMap() (slaveMap map[string]master.Slave) {
	slaveMap = make(map[string]master.Slave)
	slaveMap["slave1"] = master.Slave{URL: "http://10.0.0.122:8080", Heartbeat: time.Now()}
	slaveMap["slave2"] = master.Slave{URL: "http://10.0.1.11:8080", Heartbeat: time.Now()}
	return slaveMap
}

func TestReceiveRequestAndSendToSlave(t *testing.T) {
	testSlaveMap := make(map[string]master.Slave)
	var receivedUrl string
	testMaster := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ReceiveRequestAndSendToSlave(writer, request, testSlaveMap)
	}))

	testSlave := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		receivedUrl = request.PostFormValue("url")
	}))
	testSlaveMap["testSlaveName"] = master.Slave{testSlave.URL, time.Now(), ""}

	m := PostURLRequest{"testSlaveName", "testURL"}
	json_message, _ := json.Marshal(m)
	client := &http.Client{}
	_, err := client.Post(testMaster.URL, "application/json", strings.NewReader(string(json_message)))

	assert.Equal(t, "testURL", receivedUrl)
	assert.Nil(t, err)
}

func TestReceiveRequestAndSendToSlaveWithEmptySlaveAddress(t *testing.T) {
	testSlaveMap := make(map[string]master.Slave)

	testMaster := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ReceiveRequestAndSendToSlave(writer, request, testSlaveMap)
	}))

	testSlaveMap["testSlaveName"] = master.Slave{"", time.Now(), ""}

	m := PostURLRequest{"testSlaveName", "testURL"}
	json_message, _ := json.Marshal(m)
	client := &http.Client{}
	response, err := client.Post(testMaster.URL, "application/json", strings.NewReader(string(json_message)))
	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	receivedResponse := string(body[:])
	assert.Equal(t, "ERROR: Failed to contact slave. Slave has no URL stored.", receivedResponse)
	assert.Nil(t, err)
}

func TestParseJson(t *testing.T) {
	inputJSON := []byte(`{"DestinationSlaveName": "LeftScreen", "URLToLoadInBrowser": "http://google.com"}`)

	parsedJSON, err := parseJSON(inputJSON)
	assert.Equal(t, "LeftScreen", parsedJSON.DestinationSlaveName)
	assert.Equal(t, "http://google.com", parsedJSON.URLToLoadInBrowser)
	assert.Nil(t, err)
}

func TestParseJSONForEmptyInput(t *testing.T) {
	var inputJSON = []byte(`{}`)

	parsedJson, err := parseJSON(inputJSON)

	assert.Equal(t, "", parsedJson.DestinationSlaveName)
	assert.Equal(t, "", parsedJson.URLToLoadInBrowser)
	assert.Nil(t, err)
}

func TestParseJsonForNilInput(t *testing.T) {
	_, err := parseJSON(nil)

	assert.NotNil(t, err)
}

func TestDestinationAddressSlave(t *testing.T) {
	slaveMap := InitializeTestSlaveMap()
	destinationURL := destinationSlaveAddress("slave1", slaveMap)

	assert.Equal(t, "http://10.0.0.122:8080", destinationURL)
}

func TestDestinationAddressSlaveForEmptySlaveMap(t *testing.T) {
	slaveMap := make(map[string]master.Slave)
	destinationURL := destinationSlaveAddress("slave2", slaveMap)

	assert.Equal(t, "", destinationURL)
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
	assert.Nil(t, err)
}

func TestSendUrlValueMessageToSlaveSlaveDoesNotRespond(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
	}))
	testServer.Close()
	err := sendUrlValueMessageToSlave(testServer.URL, "http://index.hu")
	assert.NotNil(t, err)
}
