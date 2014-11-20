package main

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"net/url"
	"regexp"
)

type Slave struct {
	ID  string
	URL string
}

func TestFormHandler(t *testing.T) {
	assert.Equal(t, 200, sendGetToFormHandler("/"))
	assert.Equal(t, 301, sendGetToFormHandler("addfs"))
}

func sendGetToFormHandler(URL string) (int) {
	TEMPLATE_PATH="templates/"
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		request.URL.Path = URL
		formHandler(w,request)
	}))

	client := &http.Client{}
	resp, _ := client.Get(testServer.URL)

	return resp.StatusCode
}

func TestSubmitHandler(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		submitHandler(w,request)
	}))

	client := &http.Client{}
	resp, _ := client.PostForm(testServer.URL, url.Values{"slave-id": {"1"}, "url": {"http://www.google.com"}})

	POSTRequestBody, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	reply := parseJsonReply(POSTRequestBody).HTML

	assert.Equal(t, true, strings.Contains(reply, "http://www.google.com"))
}

func TestCheckStatusCode(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
	}))

	responseStatusCode := checkStatusCode(testServer.URL)
	assert.Equal(t, 200, responseStatusCode)
	responseStatusCode = checkStatusCode("")
	assert.Equal(t, 0, responseStatusCode)
}

func TestSendUrlAndIdToMaster(t *testing.T) {
	var numberOfMessagesSent = 0
	var url = ""
	var id = ""

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		numberOfMessagesSent++
		POSTRequestBody, _ := ioutil.ReadAll(request.Body)
		defer request.Body.Close()
		slave := parseJsonSlave(POSTRequestBody)
		url = slave.URL
		id = slave.ID
	}))
	
	sendUrlAndIdToMaster(testServer.URL, "http://index.hu", "2")
	assert.Equal(t, 1, numberOfMessagesSent)
	assert.Equal(t, "http://index.hu", url)
	assert.Equal(t, "2", id)
}

func parseJsonSlave(input []byte) (slave Slave) {
	err := json.Unmarshal(input, &slave)
	if err != nil {
		fmt.Println("error:", err)
	}
	return
}

func TestSendConfirmationMessageToUser(t *testing.T) {
	TEMPLATE_PATH="templates/"
	var responseHeader string
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		sendConfirmationMessageToUser(w,"aaaa", "bbbb", "cccc")
		responseHeader=w.Header().Get("Content-Type")
	}))
	client := &http.Client{}
	resp, _ := client.Get(testServer.URL)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "application/json", responseHeader)
	POSTRequestBody, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	reply := parseJsonReply(POSTRequestBody).HTML

	assert.Equal(t, true, strings.Contains(reply, "aaaa"))
	assert.Equal(t, true, strings.Contains(reply, "bbbb"))
	assert.Equal(t, true, strings.Contains(reply, "cccc"))
}

func parseJsonReply(input []byte) (reply Reply) {
	err := json.Unmarshal(input, &reply)
	if err != nil {
		fmt.Println("error:", err)
	}
	return
}

func TestConfirmationMessage(t *testing.T) {
	TEMPLATE_PATH="templates/"
	answer_string := parseJsonReply(confirmationMessage("aaaa", "bbbb", "cccc")).HTML
	assert.Equal(t, true, strings.Contains(answer_string, "aaaa"))
	assert.Equal(t, true, strings.Contains(answer_string, "bbbb"))
	assert.Equal(t, true, strings.Contains(answer_string, "cccc"))
}

func TestReceiveAndMapSlaveAddress(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		receiveAndMapSlaveAddress(w,request)
	}))

	client := &http.Client{}
	var testIdList IdList
	testIdList.Id = append(testIdList.Id, "testSlaveId")
	jsonMessage, _ := json.Marshal(testIdList)
	client.Post(testServer.URL, "application/json", strings.NewReader(string(jsonMessage)))

	assert.Equal(t, testIdList, id_list)
}

func TestGetIPAddress(t *testing.T) {
	IPAddress := getIPAddress(4003)
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
	assert.Equal(t, "http://10.0.0.126:1234", addProtocolAndPortToIp("10.0.0.126", 1234))
}
