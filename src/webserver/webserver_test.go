package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

type Slave struct {
	ID string
	URL string
}

func parseJson(input []byte) (slave Slave) {
	err := json.Unmarshal(input, &slave)
	if err != nil {
		fmt.Println("error:", err)
	}
	return
}

func setUpTestServerWithPath(path string) (headerContentType string) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		setMimeType(w,path)
		headerContentType = w.Header().Get("Content-type")
	}))
	sendHeadRequestTo(testServer.URL)
	return
}

func sendHeadRequestTo(url string) {
	client := &http.Client{}
	_,_ = client.Head(url)
}

func TestStatusCode(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
	}))

	responseStatusCode := statusCode(testServer.URL)

	assert.Equal(t, 200, responseStatusCode)
}

func TestSendMaster(t *testing.T) {
	var numberOfMessagesSent = 0
	var url = ""
	var id = ""

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		numberOfMessagesSent++
		POSTRequestBody, _ := ioutil.ReadAll(request.Body)
		defer request.Body.Close()
		slave := parseJson(POSTRequestBody)
		url = slave.URL
		id = slave.ID

	}))

	sendMaster(testServer.URL, "http://index.hu","2")

	assert.Equal(t, 1, numberOfMessagesSent)
	assert.Equal(t, "http://index.hu", url)
	assert.Equal(t, "2", id)
}

