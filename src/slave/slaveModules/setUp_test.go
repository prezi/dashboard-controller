package slaveModule

import (
	"net/http"
	"net/http/httptest"
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetIPAddressFromCmdLine(t *testing.T) {
	IPAddress := getIPAddressFromCmdLine(8080)
	fmt.Println(IPAddress)
	assert.Equal(t, "http://10.0.0.34:8080", IPAddress) // this will differ on each user's computer. use regexp.
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


