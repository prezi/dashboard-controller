package master

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
	"time"
	"errors"
	"net/url"
	"net"
	"strings"
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

func TestPrintServerResponseWithError(t *testing.T) {
	err := errors.New("This is error")
	printServerResponse(err,"TestSlaveName")
}

func TestPrintServerConfirmation(t *testing.T) {
	printServerResponse(nil, "HelloClient")
}

func TestSendSlaveListToWebserver(t *testing.T) {
	returnedIds := []string{"slave1", "slave2"}
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		POSTRequestBody, _ := ioutil.ReadAll(request.Body)
		defer request.Body.Close()
		var idlist IdList
		json.Unmarshal(POSTRequestBody, &idlist)
		returnedIds = idlist.Id

	}))
	slaveIPs := InitializeTestSlaveMap()
	sendSlaveListToWebserver(testServer.URL, slaveIPs)
	validIdList := []string{"slave1", "slave2"}
	sort.Strings(validIdList)
	sort.Strings(returnedIds)
	assert.Equal(t, returnedIds, validIdList)
}

func TestGetWebserverAddressWithEmptyRequest(t *testing.T) {
	request := &http.Request{}
	request.RemoteAddr = "Dummy"
	webServerAddress, err := getWebserverAddress(request)
	assert.NotNil(t, err)
	assert.Equal(t, "", webServerAddress)
}

func TestGetWebserverAddressWithEmptyport(t *testing.T) {
	request := &http.Request{}
	request.RemoteAddr = "127.0.0.1:3423"

	webServerAddress, err := getWebserverAddress(request)
	assert.Nil(t, err)
	assert.Equal(t, "http://127.0.0.1", webServerAddress)
}

func TestSendWebserverInit(t *testing.T) {
	testSlaveMap := make(map[string]Slave)

	request := &http.Request{}
	request.RemoteAddr = "127.0.0.1:3423"

	form := url.Values{}
	form.Set("message", "update me!")

	request.Form = form

	SendWebserverInit(request, testSlaveMap)

	assert.Equal(t, "http://127.0.0.1", webServerAddress)
}

func TestSendWebserverInitOnWebsite(t *testing.T) {
	testSlaveMap := InitializeTestSlaveMap()
	slave1exists := false
	slave2exists := false

	testWebServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		POSTRequestBody, _ := ioutil.ReadAll(request.Body)
		defer request.Body.Close()

		stringPostRequestBody := string(POSTRequestBody)
		slave1exists = strings.Contains(stringPostRequestBody,"slave1")
		slave2exists = strings.Contains(stringPostRequestBody,"slave2")
	}))

	testMaster := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		SendWebserverInit(request, testSlaveMap)

	}))

	webServerIp, _, _ := net.SplitHostPort(testWebServer.URL)
	webServerPort := testWebServer.URL[strings.LastIndex(testWebServer.URL,":")+1:]

	request := &http.Request{}
	request.RemoteAddr = webServerIp + ":" + webServerPort

	form := url.Values{}
	form.Set("webserverPort", webServerPort)
	form.Set("message", "update me!")

	client := &http.Client{}
	client.PostForm(testMaster.URL, form)

	assert.Equal(t, webServerIp + ":" + webServerPort, webServerAddress)
	assert.True(t, slave1exists)
	assert.True(t, slave2exists)
}
