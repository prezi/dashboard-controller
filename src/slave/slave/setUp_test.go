package slave

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestSetUp(t *testing.T) {
	port, slaveName, masterURL, OS := SetUp()
	assert.Equal(t, "8080", port)
	assert.Equal(t, "SLAVE NAME UNSPECIFIED", slaveName)
	assert.Equal(t, "http://localhost:5000", masterURL)
	assert.IsType(t, "Some OS Name", OS)
}

func TestSendIPAddressToMaster(t *testing.T) {
	var numberOfMessagesSent = 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		numberOfMessagesSent++
	})
	testServer := httptest.NewServer(handler)

	sendSlaveURLToMaster("testSlaveName", "http://localhost:8080", testServer.URL)
	assert.Equal(t, 1, numberOfMessagesSent)
}

func TestSendIPAddressToMaster_DEFAULT_SLAVE_NAME(t *testing.T) {
	var numberOfMessagesSent = 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		numberOfMessagesSent++
	})
	testServer := httptest.NewServer(handler)

	sendSlaveURLToMaster("DEFAULT_SLAVE_NAME", "http://localhost:8080", testServer.URL)
	assert.Equal(t, 1, numberOfMessagesSent)
}

