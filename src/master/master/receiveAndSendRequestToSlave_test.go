package master

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

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

func TestDestinationAddressSlave1(t *testing.T) {
	slaveMap := InitializeTestSlaveMap()
	destinationURL := destinationSlaveAddress("slave1", slaveMap)

	assert.Equal(t, "http://10.0.0.122:8080", destinationURL)
}

func TestDestinationAddressSlave2(t *testing.T) {
	slaveMap := InitializeTestSlaveMap()
	destinationURL := destinationSlaveAddress("slave2", slaveMap)

	assert.Equal(t, "http://10.0.1.11:8080", destinationURL)
}

func TestSendUrlValueMessageToSlave(t *testing.T) {
	var numberOfMessagesSent = 0
	var url = ""

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		numberOfMessagesSent++
		url = request.PostFormValue("url")
	}))

	sendUrlValueMessageToSlave(testServer.URL, "http://index.hu")

	assert.Equal(t, 1, numberOfMessagesSent)
	assert.Equal(t, "http://index.hu", url)
}
