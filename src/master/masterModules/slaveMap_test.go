package masterModule

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"encoding/json"
	"sort"
)


func TestInitializeSlaveIPs(t *testing.T) {
	slaveIPMap := initializeSlaveIPs()

	assert.Equal(t, "http://10.0.0.122:8080", slaveIPMap["1"])
	assert.Equal(t, "http://10.0.1.11:8080", slaveIPMap["2"])
}

func TestReceiveAndMapSlaveAddress(t *testing.T) {
	// name := ""
	// testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
	// 	name = request.PostFormValue("slaveName")
	// }))

	// error := sendSlaveToWebserver([]string{testServer.URL, "/receive_slave"}, "ApplePie")

	// assert.Equal(t, "ApplePie", name)
	// assert.Nil(t, error)
}

func TestSendValidSlaveToWebserver(t *testing.T) {
	// testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
	// }))

	// error := sendSlaveToWebserver([]string{testServer.URL, "/receive_slave"},  "FantasticName")

	// assert.Nil(t, error)
}

func TestPrintServerConfirmation(t *testing.T) {
	printServerResponse(nil, "HelloClient")
}

func TestSendSlaveToWebserver(t *testing.T) {
	returnedIds := []string{"1","2"}
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		POSTRequestBody,_:=ioutil.ReadAll(request.Body)
		defer request.Body.Close()
		var idlist IdList;
		json.Unmarshal(POSTRequestBody, &idlist)
		returnedIds = idlist.Id

	}))
	slaveIPs := make(map[string]string)
	slaveIPs["1"] = "http://10.0.0.122:8080"
	slaveIPs["2"] = "http://10.0.1.11:8080"
	sendSlaveToWebserver([]string{testServer.URL}, slaveIPs)
	validIdList := []string{"1","2"}
	sort.Strings(validIdList)
	sort.Strings(returnedIds)
	assert.Equal(t,returnedIds,validIdList)
}