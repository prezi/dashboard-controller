package master

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func InitializeTestSlaveMap() (slaveMap map[string]Slave) {
	slaveMap = make(map[string]Slave)
	slaveMap["slave1"] = Slave{URL: "http://10.0.0.122:8080", heartbeat: time.Now()}
	slaveMap["slave2"] = Slave{URL: "http://10.0.1.11:8080", heartbeat: time.Now()}
	return slaveMap
}

func TestSetUp(t *testing.T) {
	slaveMap := SetUp()
	assert.Equal(t, 0, len(slaveMap))
}

func TestReceiveSlaveHeartbeat(t *testing.T) {
	testSlaveMap := make(map[string]Slave)
	testSlaveName := "testSlaveName"
	longForm := "Jan 2, 2006 at 3:04pm (MST)"
	beginningOfTime, err := time.Parse(longForm, "Jan 1, 0000 at 01:01am (PST)")
	slavePort := "0000"
	newTime, _ := time.Parse(longForm, "Jan 1, 0000 at 01:01am (PST)")

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	testMaster := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		slaveIP, _, _ := net.SplitHostPort(request.RemoteAddr)
		slaveUrl := "http://" + slaveIP + ":" + slavePort
		testSlaveMap[testSlaveName] = Slave{slaveUrl, beginningOfTime, ""}
		ReceiveSlaveHeartbeat(request, testSlaveMap)
		changedSlave := testSlaveMap[testSlaveName]
		newTime = changedSlave.heartbeat
	}))

	client := &http.Client{}
	form := url.Values{}
	form.Set("slaveName", testSlaveName)
	form.Set("slavePort", slavePort)
	_, err = client.PostForm(testMaster.URL, form)

	assert.NotEqual(t, beginningOfTime, newTime)
}

func TestReceiveSlaveHeartbeatsWithDifferentAddress(t *testing.T) {
	testSlaveMap := make(map[string]Slave)
	testSlaveName := "testSlaveName"
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
		slaveUrl := "not a URL"
		testSlaveMap[testSlaveName] = Slave{slaveUrl, beginningOfTime, ""}
		ReceiveSlaveHeartbeat(request, testSlaveMap)
	}))

	client := &http.Client{}
	form := url.Values{}
	form.Set("slaveName", testSlaveName)
	form.Set("slavePort", slavePort)
	_, _ = client.PostForm(testMaster.URL, form)

	assert.Equal(t, "die", sentMessage)
}

func TestReceiveSlaveHeartbeatsNewSlaveName(t *testing.T) {
	testSlaveMap := make(map[string]Slave)
	testSlaveName := "testSlaveName"
	testSlavePort := "4006"
	exists := false

	testMaster := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		ReceiveSlaveHeartbeat(request, testSlaveMap)
		_, exists = testSlaveMap[testSlaveName]
	}))

	client := &http.Client{}
	form := url.Values{}
	form.Set("slaveName", testSlaveName)
	form.Set("slavePort", testSlavePort) //testSlave.URL[len(testSlave.URL)-5:])
	_, _ = client.PostForm(testMaster.URL, form)

	assert.True(t, exists)
}

func TestProcessRequest(t *testing.T) {
	testSlaveName := "testSlaveName"
	returnedSlaveName := ""
	returnedAddress := ""
	remoteHost := ""

	testServer := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		remoteHost, _, _ = net.SplitHostPort(request.RemoteAddr)
		returnedSlaveName, returnedAddress = processSlaveHeartbeatRequest(request)
	}))
	client := &http.Client{}
	form := url.Values{}
	form.Set("slaveName", testSlaveName)
	form.Set("slavePort", "Test")
	_, _ = client.PostForm(testServer.URL, form)
	assert.Equal(t, returnedSlaveName, testSlaveName)
	assert.Equal(t, "http://"+remoteHost+":Test", returnedAddress)
}

func TestSendKillSignalToSlave(t *testing.T) {
	message := ""
	testSlave := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		message = request.FormValue("message")
	}))
	sendKillSignalToSlave(testSlave.URL)
	assert.Equal(t, "die", message)
}

func TestMonitorSlaves(t *testing.T) {
	test_mode = true

	longForm := "Jan 2, 2006 at 3:04pm (MST)"
	beginningOfTime, _ := time.Parse(longForm, "Jan 1, 0000 at 01:01am (PST)")
	contentLength := 0

	testSlaveName := "slaveName"

	testWebServer := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		POSTRequestBody, _ := ioutil.ReadAll(request.Body)
		defer request.Body.Close()
		contentLength = len(POSTRequestBody)
	}))

	webServerAddress = testWebServer.URL
	testSlaveMap := make(map[string]Slave)
	testSlaveMap[testSlaveName] = Slave{"Dummy", beginningOfTime, ""}
	MonitorSlaves(1, testSlaveMap)

	assert.NotEqual(t, 0, contentLength)
}

func TestRemoveDeadSlaves(t *testing.T) {
	testSlaveMap := make(map[string]Slave)
	longForm := "Jan 2, 2006 at 3:04pm (MST)"
	beginningOfTime, _ := time.Parse(longForm, "Jan 1, 0000 at 01:01am (PST)")
	testSlaveMap["slave1"] = Slave{"-", beginningOfTime, ""}
	testSlaveMap["slave2"] = Slave{"-", beginningOfTime, ""}
	testSlaveMap["slave3"] = Slave{"-", time.Now(), ""}
	testSlaveMap["slave4"] = Slave{"-", time.Now(), ""}
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
	testSlaveMap := make(map[string]Slave)
	longForm := "Jan 2, 2006 at 3:04pm (MST)"
	beginningOfTime, _ := time.Parse(longForm, "Jan 1, 0000 at 01:01am (PST)")
	testSlaveMap["slave1"] = Slave{"-", beginningOfTime, ""}
	testSlaveMap["slave2"] = Slave{"-", beginningOfTime, ""}
	removeDeadSlaves(3, testSlaveMap)
	_, sl1 := testSlaveMap["slave1"]
	_, sl2 := testSlaveMap["slave2"]
	assert.False(t, sl1)
	assert.False(t, sl2)
	assert.Equal(t, 0, len(testSlaveMap))
}

func TestUpdateWebserverAddress(t *testing.T) {
	webServerAddress = "Dummy"
	webServerIP := ""
	webServerPort := "7777"

	testMaster := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		webServerIP, _, _ = net.SplitHostPort(request.RemoteAddr)
		UpdateWebserverAddress(request)
	}))

	client := &http.Client{}
	form := url.Values{}
	form.Set("message", "update me!")
	form.Set("webserverPort", webServerPort)
	_, _ = client.PostForm(testMaster.URL, form)
	assert.Equal(t, "http://"+webServerIP+":"+webServerPort, webServerAddress)
}
