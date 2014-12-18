package master

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"net/http"
	"net/http/httptest"
)

func TestGetRelativeFilePath(t *testing.T) {
	filepath := GetRelativeFilePath("assets/images")
	assert.IsType(t, "some/filepath", filepath)
}


func TestCheckStatusCode(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
	}))

	responseStatusCode := checkStatusCode(testServer.URL)
	assert.Equal(t, 200, responseStatusCode)
	responseStatusCode = checkStatusCode("")
	assert.Equal(t, 0, responseStatusCode)
}

func TestCheckStatusCodeWithoutHttp(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
	}))

	responseStatusCode := checkStatusCode(testServer.URL[7:])
	assert.Equal(t, 200, responseStatusCode)
}

func TestIfURLIsValid(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
	}))

	assert.True(t, IsURLValid(testServer.URL))
}

func TestIfURLIsInvalid(t *testing.T) {
	assert.False(t, IsURLValid(""))
}

func TestSetUp(t *testing.T) {
	proxyAddress, proxyPort := SetUpProxy()

	assert.Equal(t, "8080", proxyPort)
	assert.Equal(t, "localhost", proxyAddress)
}
