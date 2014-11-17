package slaveModule

import (
	"net/http"
	"net/http/httptest"
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
	"regexp"
)

func TestGetIPAddressFromCmdLine(t *testing.T) {
	IPAddress := getIPAddressFromCmdLine(8080)
	fmt.Println("Current IP Address: ", IPAddress)
	IPAddressRegexpPattern := "([0-9]*\\.){3}[0-9]*:[0-9]*"
	re := regexp.MustCompile(IPAddressRegexpPattern)
	assert.Equal(t, true, re.MatchString(IPAddress))
}

func TestGetMasterReceiveSlaveAddress(t *testing.T) {
	masterAddress := getMasterReceiveSlaveAddress("http://localhost:5000")
	assert.Equal(t, "http://localhost:5000/receive_slave", masterAddress)
}

func TestSendIPAddressToMaster(t *testing.T) {
	var numberOfMessagesSent = 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		numberOfMessagesSent++
	})
	testServer := httptest.NewServer(handler)

	sendIPAddressToMaster("testSlaveName", "http://localhost:8080", testServer.URL)
	assert.Equal(t, 1, numberOfMessagesSent)
}

func TestSendIPAddressToMaster_DEFAULT_SLAVE_NAME(t *testing.T) {
	var numberOfMessagesSent = 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		numberOfMessagesSent++
	})
	testServer := httptest.NewServer(handler)

	sendIPAddressToMaster("DEFAULT_SLAVE_NAME", "http://localhost:8080", testServer.URL)
	assert.Equal(t, 1, numberOfMessagesSent)
}


