package proxyMonitor

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	TEST_SLAVE_IP = "10.0.214"
)

// TODO: There is a pretty much identical function in slave for SendURLValueMessageToSlave. Refactor into network package.
func TestRequestProxyToAddNewSlaveToIPTables(t *testing.T) {
	var numberOfMessagesSent = 0
	var slaveIP = ""

	testProxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		numberOfMessagesSent++
		slaveIP = request.PostFormValue("IPAddressToAdd")
	}))

	err := RequestProxyToAddNewSlaveToIPTables(testProxy.URL, TEST_SLAVE_IP)

	assert.Equal(t, 1, numberOfMessagesSent)
	assert.Equal(t, TEST_SLAVE_IP, slaveIP)
	assert.Nil(t, err)
}
