package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestInitializeSlaveIPs(t *testing.T) {
	slaveIPMap := initializeSlaveIPs()

	assert.Equal(t, "http://10.0.0.42:8080", slaveIPMap["1"])
	assert.Equal(t, "http://10.0.0.231:8080", slaveIPMap["2"])
}

func TestParseJson(t *testing.T) {
	var inputJson = []byte(`{"ID":"LeftScreen","URL":"http://google.com"}`)

	parsedJson, err := parseJson(inputJson)

	assert.Equal(t, "LeftScreen", parsedJson.ID)
	assert.Equal(t, "http://google.com", parsedJson.URL)
	assert.Nil(t, err)
}

func TestParseJsonForEmptyInput(t *testing.T) {
	var inputJson = []byte(`{}`)

	parsedJson, err := parseJson(inputJson)

	assert.Equal(t, "", parsedJson.ID)
	assert.Equal(t, "", parsedJson.URL)
	assert.Nil(t, err)
}

func TestDestinationUrlSlave1(t *testing.T) {
	slaveIPMap = initializeSlaveIPs()
	destinationURL := destinationSlaveAddress("1")

	assert.Equal(t, "http://10.0.0.42:8080", destinationURL)
}

func TestDestinationUrlSlave2(t *testing.T) {
	slaveIPMap = initializeSlaveIPs()
	destinationURL := destinationSlaveAddress("2")

	assert.Equal(t, "http://10.0.0.231:8080", destinationURL)
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

func TestReceiveAndMapSlaveAddress(t *testing.T) {
	name := ""
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		name = request.PostFormValue("slaveName")
	}))

	error := sendSlaveToWebserver([]string{testServer.URL, "/receive_slave"}, "ApplePie")

	assert.Equal(t, "ApplePie", name)
	assert.Nil(t, error)
}

func TestSendValidSlaveToWebserver(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
	}))

	error := sendSlaveToWebserver([]string{testServer.URL, "/receive_slave"},  "FantasticName")

	assert.Nil(t, error)
}

func TestPrintServerConfirmation(t *testing.T) {
	printServerResponse(nil, "HelloClient")
}
