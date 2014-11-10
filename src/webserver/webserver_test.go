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

	sendMaster(testServer.URL, "http://index.hu", "2")

	assert.Equal(t, 1, numberOfMessagesSent)
	assert.Equal(t, "http://index.hu", url)
	assert.Equal(t, "2", id)
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
