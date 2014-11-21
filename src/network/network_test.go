package network

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"regexp"
	"net/http"
	"net/http/httptest"
	"errors"
)

func TestGetLocalIPAddress(t *testing.T) {
	IPAddress := GetLocalIPAddress()
	IPAddressRegexpPattern := "([0-9]*\\.){3}[0-9]*"
	re := regexp.MustCompile(IPAddressRegexpPattern)
	assert.Equal(t, true, re.MatchString(IPAddress))
}

func TestAddProtocolAndPortToIP(t *testing.T) {
	assert.Equal(t, "http://10.0.0.126:1234", AddProtocolAndPortToIP("10.0.0.126", 1234))
}

func TestErrorHandler(t *testing.T) {
	assert.Equal(t, false, ErrorHandler(nil, "This is an error message."))
	err := errors.New("This is an error message.")
	assert.Equal(t, true, ErrorHandler(err, "%v"))
}

func TestSetMasterUrl(t *testing.T) {
	masterURL := SetMasterUrl()
	assert.Equal(t, "http://localhost:5000", masterURL )
}

func TestSendSlaveURLToMaster(t *testing.T) {
	var numberOfMessagesSent = 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		numberOfMessagesSent++
	})
	testServer := httptest.NewServer(handler)

	sendSlaveURLToMaster("testSlaveName", "http://localhost:8080", testServer.URL)
	assert.Equal(t, 1, numberOfMessagesSent)
}

func TestSendSlaveURLToMaster_DEFAULT_SLAVE_NAME(t *testing.T) {
	var numberOfMessagesSent = 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		numberOfMessagesSent++
	})
	testServer := httptest.NewServer(handler)

	sendSlaveURLToMaster("DEFAULT_SLAVE_NAME", "http://localhost:8080", testServer.URL)
	assert.Equal(t, 1, numberOfMessagesSent)
}
