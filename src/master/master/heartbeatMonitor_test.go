package master

import (
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
	"testing"
)

func TestMonitorSlaveHeartbeats(t *testing.T) {
	testSlaveMap := make(map[string]Slave)
	testSlaveName := "testSlaveName"
	longForm := "Jan 2, 2006 at 3:04pm (MST)"
	beginningOfTime, err:= time.Parse(longForm, "Jan 1, 0000 at 01:01am (PST)")
	slavePort := "0000"
	newTime, _:= time.Parse(longForm, "Jan 1, 0000 at 01:01am (PST)")

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	testMaster:= httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		slaveIP, _, _ := net.SplitHostPort(request.RemoteAddr)
		slaveUrl := "http://" + slaveIP + ":" + slavePort
		testSlaveMap[testSlaveName] = Slave{slaveUrl, beginningOfTime, ""}
		MonitorSlaveHeartbeats(request,testSlaveMap)
		changedSlave := testSlaveMap[testSlaveName]
		newTime = changedSlave.heartbeat
	}))

	client := &http.Client{}
	form := url.Values{}
	form.Set("slaveName", testSlaveName)
	form.Set("slavePort", slavePort)
	_, err = client.PostForm(testMaster.URL, form)

	assert.NotEqual(t,beginningOfTime,newTime)
}

func TestMonitorSlaveHeartbeatsWithDifferentAddress(t *testing.T) {
	testSlaveMap := make(map[string]Slave)
	testSlaveName := "testSlaveName"
	longForm := "Jan 2, 2006 at 3:04pm (MST)"
	beginningOfTime, _ := time.Parse(longForm, "Jan 1, 0000 at 01:01am (PST)")
	sentMessage := ""

	testSlave:= httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		sentMessage = request.FormValue("message")
	}))
	slaveURL, _ := url.Parse(testSlave.URL)
	slaveIP, slavePort, _ := net.SplitHostPort(slaveURL.Host)
	testMaster:= httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		request.URL.Host = slaveIP
		slaveUrl := "not a URL"
		testSlaveMap[testSlaveName] = Slave{slaveUrl, beginningOfTime, ""}
		MonitorSlaveHeartbeats(request,testSlaveMap)
	}))

	client := &http.Client{}
	form := url.Values{}
	form.Set("slaveName", testSlaveName)
	form.Set("slavePort", slavePort)
	_, _ = client.PostForm(testMaster.URL, form)

	assert.Equal(t,"die",sentMessage)
}

func TestMonitorSlaveHeartbeatsNewSlaveName(t *testing.T) {
	testSlaveMap := make(map[string]Slave)
	testSlaveName := "testSlaveName"
	testSlavePort := "4006"
	exists := false

	testMaster := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		MonitorSlaveHeartbeats(request,testSlaveMap)
		_, exists = testSlaveMap[testSlaveName]
	}))

	client := &http.Client{}
	form := url.Values{}
	form.Set("slaveName", testSlaveName)
	form.Set("slavePort", testSlavePort)//testSlave.URL[len(testSlave.URL)-5:])
	_, _ = client.PostForm(testMaster.URL, form)

	assert.True(t,exists)
}

func TestProcessRequest(t *testing.T) {
	testSlaveName := "testSlaveName"
	returnedSlaveName := ""
	returnedAddress := ""
	remoteHost := ""


	testServer := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		remoteHost, _, _ = net.SplitHostPort(request.RemoteAddr)
		fmt.Println("###################",remoteHost)
		returnedSlaveName, returnedAddress = processRequest(request)
	}))
	client := &http.Client{}
	form := url.Values{}
	form.Set("slaveName", testSlaveName)
	form.Set("slavePort", "Test")
	_, _ = client.PostForm(testServer.URL, form)
	assert.Equal(t, returnedSlaveName, testSlaveName)
	assert.Equal(t, "http://" + remoteHost + ":Test", returnedAddress)
}
