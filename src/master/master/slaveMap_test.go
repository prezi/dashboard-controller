package master

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"testing"
	"time"
)

func InitializeTestSlaveMap() (slaveMap map[string]Slave) {
	slaveMap = make(map[string]Slave)
	slaveMap["slave1"] = Slave{URL: "http://10.0.0.122:8080", heartbeat: time.Now()}
	slaveMap["slave2"] = Slave{URL: "http://10.0.1.11:8080", heartbeat: time.Now()}
	return slaveMap
}

func TestInitializeSlaveMap(t *testing.T) {
	slaveMap := InitializeTestSlaveMap()

	assert.Equal(t, "http://10.0.0.122:8080", slaveMap["slave1"].URL)
	assert.Equal(t, "http://10.0.1.11:8080", slaveMap["slave2"].URL)
}

func TestPrintServerConfirmation(t *testing.T) {
	printServerResponse(nil, "HelloClient")
}

func TestSendSlaveListToWebserver(t *testing.T) {
	returnedIds := []string{"slave1", "slave2"}
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		POSTRequestBody, _ := ioutil.ReadAll(request.Body)
		defer request.Body.Close()
		var idlist IdList
		json.Unmarshal(POSTRequestBody, &idlist)
		returnedIds = idlist.Id

	}))
	slaveIPs := InitializeTestSlaveMap()
	sendSlaveListToWebserver(testServer.URL, slaveIPs)
	validIdList := []string{"slave1", "slave2"}
	sort.Strings(validIdList)
	sort.Strings(returnedIds)
	assert.Equal(t, returnedIds, validIdList)
}

func TestWebserverRequestSlaveListAtStartUp(t *testing.T) {
	slaveMap := initializeSlaveMap()
	WebserverRequestSlaveIdsHandler := func(w http.ResponseWriter, r *http.Request) {
		WebserverRequestSlaveListAtStartUp(w, r, slaveMap)
	}

	testServer := httptest.NewServer(http.HandlerFunc(WebserverRequestSlaveIdsHandler))

	client := &http.Client{}
	form := url.Values{}
	form.Set("webserverPort", "5000")
	resp, err := client.PostForm(testServer.URL, form)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, nil, err)
	form.Set("webserverPort", "")
	resp, err = client.PostForm(testServer.URL, form)
	assert.Equal(t, 500, resp.StatusCode)
}
