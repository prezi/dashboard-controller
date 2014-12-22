package master

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

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

func TestDefaultProxyPort(t *testing.T) {
	assert.Equal(t, "8080", GetProxyPort())
}
