package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
)

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

func TestSetHtmlMimeType(t *testing.T) {
	var path = "file.html"
	assert.Equal(t, "text/html; charset=utf-8", setUpTestServerWithPath(path))
}

func TestSetCssMimeType(t *testing.T) {
	var path = "file.css"
	assert.Equal(t, "text/css; charset=utf-8", setUpTestServerWithPath(path))
}

func TestSetJsMimeType(t *testing.T) {
	var path = "file.js"
	assert.Equal(t, "application/javascript", setUpTestServerWithPath(path))
}

