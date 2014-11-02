package main

import (
//	"fmt"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
	"strings"
)

type Greeting struct {
	Text string
}

func TestJsonCanBeParsed(t *testing.T) {
	var jsonBlob = []byte(`{"Text":"Platypus"}`)

	data := ParseGreeting(jsonBlob)

	assert.Equal(t, "Platypus", data.Text)
}

func TestMessageIsSent(t *testing.T) {
	var numberOfMessagesSent = 0
	var message = ""
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		numberOfMessagesSent++
		requestBody, _ := ioutil.ReadAll(r.Body)
		var requestMessage Greeting
		_ = json.Unmarshal(requestBody, &requestMessage)

		message = requestMessage.Text
	}))

	sendMessageToServer(testServer.URL)

	assert.Equal(t, 1, numberOfMessagesSent)
	assert.Equal(t, "hello", message)
}

func sendMessageToServer(url string) {
	client := &http.Client{}
	var greeting Greeting
	greeting.Text = "hello"
	json_message, _ := json.Marshal(greeting)
	_, _ = client.Post(url, "application/json", strings.NewReader(string(json_message)))
}
