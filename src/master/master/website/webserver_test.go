package website

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type PostURLRequest struct {
	DestinationSlaveName string
	URLToLoadInBrowser   string
}

func parseJsonSlave(input []byte) (slave PostURLRequest) {
	err := json.Unmarshal(input, &slave)
	if err != nil {
		fmt.Println("error:", err)
	}
	return
}

func parseJsonReply(input []byte) (reply StatusMessage) {
	err := json.Unmarshal(input, &reply)
	if err != nil {
		fmt.Println("error:", err)
	}
	return
}

func sendGetToFormHandler(URL string) int {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		request.URL.Path = URL
	}))

	client := &http.Client{}
	resp, _ := client.Get(testServer.URL)

	return resp.StatusCode
}

func TestFormHandler(t *testing.T) {
	assert.Equal(t, 200, sendGetToFormHandler("/"))
}

func TestStatusMessageForAvailableServer(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
	}))
	statusMessage := isUrlValid(testServer.URL)

	assert.Equal(t, true, statusMessage)
}

func TestStatusMessageForUnavailableServer(t *testing.T) {
	statusMessage := isUrlValid("")

	assert.Equal(t, false, statusMessage)
}

func TestCheckStatusCode(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
	}))

	responseStatusCode := checkStatusCode(testServer.URL)
	assert.Equal(t, 200, responseStatusCode)
	responseStatusCode = checkStatusCode("")
	assert.Equal(t, 0, responseStatusCode)
}

func TestIfUrlIsValid(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
	}))

	assert.True(t, isUrlValid(testServer.URL))
}

func TestIfUrlIsInvalid(t *testing.T) {
	assert.False(t, isUrlValid(""))
}
