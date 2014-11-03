package main

import (
	"io/ioutil"
	"encoding/json"
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

func TestMessageIsSent(t *testing.T) {
	var numberOfMessagesSent = 0
	var messageID = ""
	var messageURL = ""

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		numberOfMessagesSent++
		requestBody, _ := ioutil.ReadAll(r.Body)
		var requestMessage Slave
		_ = json.Unmarshal(requestBody, &requestMessage)

		messageID = requestMessage.ID
		messageURL = requestMessage.URL
	}))

	sendMessageToServer(testServer.URL)

	assert.Equal(t, 1, numberOfMessagesSent)
	assert.Equal(t, "LeftScreen", messageID)
	assert.Equal(t, "http://index.hu", messageURL)
}
