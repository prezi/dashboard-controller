package slaveMapHandler

import (
	"github.com/stretchr/testify/assert"
	"master/master"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gorilla/mux"
	"time"
)

func TestInitiateEmptySlaveMapHandler(t *testing.T) {
	router := mux.NewRouter()
	responseRecorder := httptest.NewRecorder()

	slaveMap := make(map[string]master.Slave)

	InitiateSlaveMapHandler(router, slaveMap)

	request, _ := http.NewRequest("GET", "/slavemap", nil)

	router.ServeHTTP(responseRecorder, request)

	assert.Equal(t, 200, responseRecorder.Code)
	assert.Equal(t, "null", responseRecorder.Body.String())
}

func TestInitiateNonEmptySlaveMapHandler(t *testing.T) {
	router := mux.NewRouter()
	responseRecorder := httptest.NewRecorder()

	slaveMap := make(map[string]master.Slave)
	slaveMap["slave1"] = master.Slave{URL: "http://10.0.0.122:8080", Heartbeat: time.Now(), PreviouslyDisplayedURL: "http://www.google.com", DisplayedURL: "http://www.prezi.com"}

	InitiateSlaveMapHandler(router, slaveMap)

	request, _ := http.NewRequest("GET", "/slavemap", nil)

	router.ServeHTTP(responseRecorder, request)

	assert.Equal(t, 200, responseRecorder.Code)
	assert.Equal(t, "[\"slave1\"]", responseRecorder.Body.String())
}
