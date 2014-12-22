package slaveMonitor

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"master/master"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"network"
	"testing"
	"time"
)

const (
	TEST_SLAVE_NAME = "testSlaveName"
	TEST_SLAVE_PORT = "0000"
)

func TestSetUp(t *testing.T) {
	slaveMap := master.GetSlaveMap()
	assert.Equal(t, 0, len(slaveMap))
}

func TestReceiveSlaveHeartbeat(t *testing.T) {
	testSlaveMap := make(map[string]master.Slave)
	longForm := "Jan 2, 2006 at 3:04pm (MST)"
	beginningOfTime, err := time.Parse(longForm, "Jan 1, 0000 at 01:01am (PST)")
	newTime, _ := time.Parse(longForm, "Jan 1, 0000 at 01:01am (PST)")

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	testMaster := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		slaveIP, _, _ := net.SplitHostPort(request.RemoteAddr)
		slaveURL := "http://" + slaveIP + ":" + TEST_SLAVE_PORT
		testSlaveMap[TEST_SLAVE_NAME] = master.Slave{URL: slaveURL, Heartbeat: beginningOfTime}
		ReceiveSlaveHeartbeat(request, testSlaveMap)
		changedSlave := testSlaveMap[TEST_SLAVE_NAME]
		newTime = changedSlave.Heartbeat
	}))

	client := &http.Client{}
	testForm := network.CreateFormWithInitialValues(map[string]string{"slaveName": TEST_SLAVE_NAME, "slavePort": TEST_SLAVE_PORT})
	_, err = client.PostForm(testMaster.URL, testForm)

	assert.NotEqual(t, beginningOfTime, newTime)
}

func TestReceiveSlaveHeartbeatsWithDifferentAddress(t *testing.T) {
	testSlaveMap := make(map[string]master.Slave)
	longForm := "Jan 2, 2006 at 3:04pm (MST)"
	beginningOfTime, _ := time.Parse(longForm, "Jan 1, 0000 at 01:01am (PST)")
	sentMessage := ""

	testSlave := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		sentMessage = request.FormValue("message")
	}))
	slaveURL, _ := url.Parse(testSlave.URL)
	slaveIP, slavePort, _ := net.SplitHostPort(slaveURL.Host)
	testMaster := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		request.URL.Host = slaveIP
		slaveURL := "not a URL"
		testSlaveMap[TEST_SLAVE_NAME] = master.Slave{URL: slaveURL, Heartbeat: beginningOfTime}
		ReceiveSlaveHeartbeat(request, testSlaveMap)
	}))

	client := &http.Client{}
	testForm := network.CreateFormWithInitialValues(map[string]string{"slaveName": TEST_SLAVE_NAME, "slavePort": slavePort})
	_, _ = client.PostForm(testMaster.URL, testForm)

	assert.Equal(t, "die", sentMessage)
}

func TestReceiveSlaveHeartbeatsNewSlaveName(t *testing.T) {
	testSlaveMap := make(map[string]master.Slave)
	exists := false

	testMaster := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		ReceiveSlaveHeartbeat(request, testSlaveMap)
		_, exists = testSlaveMap[TEST_SLAVE_NAME]
	}))

	client := &http.Client{}
	testForm := network.CreateFormWithInitialValues(map[string]string{"slaveName": TEST_SLAVE_NAME, "slavePort": TEST_SLAVE_PORT})
	_, _ = client.PostForm(testMaster.URL, testForm)

	assert.True(t, exists)
}

func TestProcessRequest(t *testing.T) {
	returnedSlaveName := ""
	returnedAddress := ""
	remoteHost := ""

	testServer := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		remoteHost, _, _ = net.SplitHostPort(request.RemoteAddr)
		returnedSlaveName, returnedAddress = processSlaveHeartbeatRequest(request)
	}))
	client := &http.Client{}
	testForm := network.CreateFormWithInitialValues(map[string]string{"slaveName": TEST_SLAVE_NAME, "slavePort": TEST_SLAVE_PORT})
	_, _ = client.PostForm(testServer.URL, testForm)
	assert.Equal(t, returnedSlaveName, TEST_SLAVE_NAME)
	assert.Equal(t, "http://"+remoteHost+":0000", returnedAddress)
}

func TestSendKillSignalToSlave(t *testing.T) {
	message := ""
	testSlave := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		message = request.FormValue("message")
	}))
	sendKillSignalToSlave(testSlave.URL)
	assert.Equal(t, "die", message)
}

func TestRemoveDeadSlaves(t *testing.T) {
	testSlaveMap := make(map[string]master.Slave)
	longForm := "Jan 2, 2006 at 3:04pm (MST)"
	beginningOfTime, _ := time.Parse(longForm, "Jan 1, 0000 at 01:01am (PST)")
	testSlaveMap["slave1"] = master.Slave{URL: "-", Heartbeat: beginningOfTime}
	testSlaveMap["slave2"] = master.Slave{URL: "-", Heartbeat: beginningOfTime}
	testSlaveMap["slave3"] = master.Slave{URL: "-", Heartbeat: time.Now()}
	testSlaveMap["slave4"] = master.Slave{URL: "-", Heartbeat: time.Now()}
	removeDeadSlaves(3, testSlaveMap)
	_, sl1 := testSlaveMap["slave1"]
	_, sl2 := testSlaveMap["slave2"]
	_, sl3 := testSlaveMap["slave3"]
	_, sl4 := testSlaveMap["slave4"]
	assert.False(t, sl1)
	assert.False(t, sl2)
	assert.True(t, sl3)
	assert.True(t, sl4)
}

func TestRemoveDeadSlavesRemoveAll(t *testing.T) {
	testSlaveMap := make(map[string]master.Slave)
	longForm := "Jan 2, 2006 at 3:04pm (MST)"
	beginningOfTime, _ := time.Parse(longForm, "Jan 1, 0000 at 01:01am (PST)")
	testSlaveMap["slave1"] = master.Slave{URL: "-", Heartbeat: beginningOfTime}
	testSlaveMap["slave2"] = master.Slave{URL: "-", Heartbeat: beginningOfTime}
	removeDeadSlaves(3, testSlaveMap)
	_, sl1 := testSlaveMap["slave1"]
	_, sl2 := testSlaveMap["slave2"]
	assert.False(t, sl1)
	assert.False(t, sl2)
	assert.Equal(t, 0, len(testSlaveMap))
}
