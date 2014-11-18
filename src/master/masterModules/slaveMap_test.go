package masterModule

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestInitializeSlaveIPs(t *testing.T) {
	slaveIPMap := initializeSlaveIPs()

	assert.Equal(t, "http://10.0.0.122:8080", slaveIPMap["1"])
	assert.Equal(t, "http://10.0.1.11:8080", slaveIPMap["2"])
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
