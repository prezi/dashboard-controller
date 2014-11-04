package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
	// "fmt"
)

func TestStatusCode(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
	}))

	responseStatusCode := statusCode(testServer.URL)

	assert.Equal(t, 200, responseStatusCode)
}

func TestSetMimeType(t *testing.T) {
	var headerContentType = ""
	
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		setMimeType(w)
		headerContentType = w.Header().Get("Content-type")
	}))

	sendHeadRequestTo(testServer.URL)

	assert.Equal(t, "text/html", headerContentType)
}

func setMimeType(responseWriter http.ResponseWriter) {
	responseWriter.Header().Set("Content-type", "text/css")
}


func sendHeadRequestTo(url string) {
	client := &http.Client{}

	_,_ = client.Head(url)

}

// func TestReply(t *testing.T) {

// }

// func TestSendInfo(t *testing.T) {

// }
