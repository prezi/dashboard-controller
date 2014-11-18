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
)

type Slave struct {
	ID  string
	URL string
}

// type MockResponseWriter interface {
//         Header() http.Header
//         Write([]byte) (int, error)
//  //       WriteHeader(int)
// }

func setUpTestServerWithPath(path string) (headerContentType string) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		//setMimeType(w, path)
		headerContentType = w.Header().Get("Content-type")
	}))
	sendHeadRequestTo(testServer.URL)
	return
}

func sendHeadRequestTo(url string) {
	client := &http.Client{}
	_, _ = client.Head(url)
}

func parseJsonSlave(input []byte) (slave Slave) {
	err := json.Unmarshal(input, &slave)
	if err != nil {
		fmt.Println("error:", err)
	}
	return
}

func parseJsonReply(input []byte) (reply Reply) {
	err := json.Unmarshal(input, &reply)
	if err != nil {
		fmt.Println("error:", err)
	}
	return
}

func TestStatusCode(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
	}))

	responseStatusCode := statusCode(testServer.URL)

	assert.Equal(t, 200, responseStatusCode)
	responseStatusCode = statusCode("")

	assert.Equal(t, 0, responseStatusCode)
}

func TestSendMaster(t *testing.T) {
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

	returnValue := sendMaster(testServer.URL, "http://index.hu", "2")

	assert.Equal(t, 1, numberOfMessagesSent)
	assert.Equal(t, "http://index.hu", url)
	assert.Equal(t, "2", id)
	assert.Equal(t, 1, returnValue)
}

func TestReply(t *testing.T) {
	TEMPLATE_PATH="templates/"
	answer_string := parseJsonReply(reply("aaaa", "bbbb", "cccc")).HTML
	assert.Equal(t, true, strings.Contains(answer_string, "aaaa"))
	assert.Equal(t, true, strings.Contains(answer_string, "bbbb"))
	assert.Equal(t, true, strings.Contains(answer_string, "cccc"))
}

func TestSendInfo(t *testing.T) {
	//mock_response_writer = MockResponseWriter;
	//var mock_response_writer http.ResponseWriter
	TEMPLATE_PATH="templates/"
	var responseHeader string
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		sendInfo(w,"aaaa", "bbbb", "cccc")
		responseHeader=w.Header().Get("Content-Type")
	}))
	//sendHeadRequestTo(testServer.URL)
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

func TestFormHandler(t *testing.T) {
	assert.Equal(t, 200, sendGetToFormHandler("/"))
	assert.Equal(t, 301, sendGetToFormHandler("addfs"))
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

func TestReceiveAndMapSlaveAddress(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		receiveAndMapSlaveAddress(w,request)
	}))

	client := &http.Client{}
	resp, _ := client.PostForm(testServer.URL, url.Values{"slaveName":{"3"}})
	POSTRequestBody, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	assert.Equal(t, 0, len(POSTRequestBody))
	assert.Equal(t, "3", id_list.Id[len(id_list.Id)-1])
}
