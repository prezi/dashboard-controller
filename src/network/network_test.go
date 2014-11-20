package network

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"regexp"
	"net/http"
	"net/http/httptest"
)

func TestGetUrl(t *testing.T) {
	IPAddress := GetUrl("4003")
	IPAddressRegexpPattern := "([0-9]*\\.){3}[0-9]*:[0-9]*"
	re := regexp.MustCompile(IPAddressRegexpPattern)

	assert.Equal(t, true, re.MatchString(IPAddress))
}

func TestGetIPAddressFromCmdLine(t *testing.T) {
	IPAddress := getIPAddressBytesFromCmdLine()
	temp_string := ""

	assert.IsType(t, temp_string, IPAddress)
}

func TestParseIpAddress(t *testing.T) {
	parsedIPAddress := parseIpAddress(getIPAddressBytesFromCmdLine())
	temp_string := ""

	assert.IsType(t, temp_string, parsedIPAddress)
}

func TestAddProtocolAndPortToIp(t *testing.T) {
	assert.Equal(t, "http://10.0.0.126:1234", addProtocolAndPortToIp("10.0.0.126", "1234"))
}

func TestGetMasterUrl(t *testing.T) {
	masterAddress := GetMasterUrl("localhost:5000")

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
