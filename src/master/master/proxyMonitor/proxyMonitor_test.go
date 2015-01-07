package proxyMonitor

import (
	"github.com/stretchr/testify/assert"
	"master/master"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	TEST_SLAVE_IP   = "10.0.214"
	TEST_SLAVE_URL  = "http://10.0.214:8686"
	TEST_SLAVE_NAME = "yoyo"
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

func TestGetSlaveIPAddresses(t *testing.T) {
	slaveMap := make(map[string]master.Slave)
	slaveMap[TEST_SLAVE_NAME] = master.Slave{URL: TEST_SLAVE_URL, Heartbeat: time.Now(), PreviouslyDisplayedURL: "http://google.com", DisplayedURL: "http://google.com"}
	IPAddresses := getSlaveIPAddresses(slaveMap)

	assert.Equal(t, []string{TEST_SLAVE_IP}, IPAddresses)
}
