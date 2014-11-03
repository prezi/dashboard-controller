package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestJsonCanBeParsed(t *testing.T) {
	var inputJson = []byte(`{"ID":"LeftScreen","URL":"http://google.com"}`)

	parsedJson := parseJson(inputJson)

	assert.Equal(t, "LeftScreen", parsedJson.ID)
	assert.Equal(t, "http://google.com", parsedJson.URL)
}

func TestUrlValueMessageIsSent(t *testing.T) {
	var numberOfMessagesSent = 0
	var url = ""

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		numberOfMessagesSent++
		url = request.PostFormValue("url")
	}))

	sendUrlValueMessageToServer(testServer.URL, "http://index.hu")

	assert.Equal(t, 1, numberOfMessagesSent)
	assert.Equal(t, "http://index.hu", url)
}
