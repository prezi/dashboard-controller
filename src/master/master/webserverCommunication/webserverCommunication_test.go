package webserverCommunication

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"master/master"
	"net"
	"net/http"
	"net/http/httptest"
	"network"
	"sort"
	"strings"
	"testing"
	"time"
)

func InitializeTestSlaveMap() (slaveMap map[string]master.Slave) {
	slaveMap = make(map[string]master.Slave)
	slaveMap["slave1"] = master.Slave{URL: "http://10.0.0.122:8080", Heartbeat: time.Now()}
	slaveMap["slave2"] = master.Slave{URL: "http://10.0.1.11:8080", Heartbeat: time.Now()}
	return slaveMap
}

func TestUpdateWebserverAddress(t *testing.T) {
	webServerAddress := "Dummy"
	webServerIP := ""
	webServerPort := "7777"

	testMaster := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		webServerIP, _, _ = net.SplitHostPort(request.RemoteAddr)
		webServerAddress = UpdateWebserverAddress(request, webServerAddress)
	}))

	client := &http.Client{}
	form := network.CreateFormWithInitialValues(map[string]string{"message": "update me!", "webserverPort": webServerPort})
	_, _ = client.PostForm(testMaster.URL, form)
	assert.Equal(t, "http://"+webServerIP+":"+webServerPort, webServerAddress)
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
	SendSlaveListToWebserver(testServer.URL, slaveIPs)
	validIdList := []string{"slave1", "slave2"}
	sort.Strings(validIdList)
	sort.Strings(returnedIds)
	assert.Equal(t, returnedIds, validIdList)
}

func TestGetWebserverAddressWithEmptyRequest(t *testing.T) {
	request := &http.Request{}
	webServerAddress, err := getWebserverAddress(request)
	assert.NotNil(t, err)
	assert.Equal(t, "", webServerAddress)
}

func TestGetWebserverAddressWithEmptyPort(t *testing.T) {
	request := &http.Request{}
	request.RemoteAddr = "127.0.0.1:3423"

	webServerAddress, err := getWebserverAddress(request)
	assert.Nil(t, err)
	assert.Equal(t, "http://127.0.0.1", webServerAddress)
}

func TestSendWebserverInit(t *testing.T) {
	testSlaveMap := make(map[string]master.Slave)

	request := &http.Request{}
	request.RemoteAddr = "127.0.0.1:3423"

	form := network.CreateFormWithInitialValues(map[string]string{"message": "update me!"})

	request.Form = form

	TestWebServerAddress := "http://1as;dlkfjdkls;j"
	TestWebServerAddress = SendWebserverInit(request, testSlaveMap)

	assert.Equal(t, "http://127.0.0.1", TestWebServerAddress)
}

func TestSendWebserverInitOnWebsite(t *testing.T) {
	testSlaveMap := InitializeTestSlaveMap()
	TestWebServerAddress := "http://localhost:4003"
	slave1exists := false
	slave2exists := false

	testWebServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		POSTRequestBody, _ := ioutil.ReadAll(request.Body)
		defer request.Body.Close()

		stringPostRequestBody := string(POSTRequestBody)
		slave1exists = strings.Contains(stringPostRequestBody, "slave1")
		slave2exists = strings.Contains(stringPostRequestBody, "slave2")
	}))

	testMaster := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {

		TestWebServerAddress = SendWebserverInit(request, testSlaveMap)
	}))

	webServerIp, _, _ := net.SplitHostPort(testWebServer.URL)
	webServerPort := testWebServer.URL[strings.LastIndex(testWebServer.URL, ":")+1:]

	request := &http.Request{}
	request.RemoteAddr = webServerIp + ":" + webServerPort

	form := network.CreateFormWithInitialValues(map[string]string{"message": "update me!", "webserverPort": webServerPort})

	client := &http.Client{}
	client.PostForm(testMaster.URL, form)

	assert.Equal(t, webServerIp+":"+webServerPort, TestWebServerAddress)
	assert.True(t, slave1exists)
	assert.True(t, slave2exists)
}
