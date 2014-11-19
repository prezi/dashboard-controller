package slaveModule

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
	"regexp"
)

func TestSetUp(t *testing.T) {
	port, slaveName, masterIP, OS := SetUp()
	assert.Equal(t, port, 8080)
	assert.Equal(t, slaveName, "SLAVE NAME UNSPECIFIED")
	assert.Equal(t, masterIP, "localhost:5000")
	assert.IsType(t, "Some OS Name", OS)
}

func TestGetIPAddressFromCmdLine(t *testing.T) {
	IPAddress := getIPAddressFromCmdLine(8080)
	IPAddressRegexpPattern := "([0-9]*\\.){3}[0-9]*:[0-9]*"
	re := regexp.MustCompile(IPAddressRegexpPattern)
	assert.Equal(t, true, re.MatchString(IPAddress))
}

func TestGetMasterReceiveSlaveAddress(t *testing.T) {
	masterAddress := getMasterReceiveSlaveAddress("localhost:5000")
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

