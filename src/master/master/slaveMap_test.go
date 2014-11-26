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
)

func TestInitializeSlaveMap(t *testing.T) {
	slaveMap := initializeSlaveMap()

	assert.Equal(t, "http://10.0.0.122:8080", slaveMap["slave1"].URL)
	assert.Equal(t, "http://10.0.1.11:8080", slaveMap["slave2"].URL)
}

func TestPrintServerConfirmation(t *testing.T) {
	printServerResponse(nil, "HelloClient")
}

func TestSendSlaveToWebserver(t *testing.T) {
	returnedIds := []string{"slave1", "slave2"}
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		POSTRequestBody, _ := ioutil.ReadAll(request.Body)
		defer request.Body.Close()
		var idlist IdList
		json.Unmarshal(POSTRequestBody, &idlist)
		returnedIds = idlist.Id

	}))
	slaveIPs := initializeSlaveMap()
	sendSlaveToWebserver(testServer.URL, slaveIPs)
	validIdList := []string{"slave1", "slave2"}
	sort.Strings(validIdList)
	sort.Strings(returnedIds)
	assert.Equal(t, returnedIds, validIdList)
}

func TestWebserverRequestSlaveIds(t *testing.T) {
	slaveMap := initializeSlaveMap()
	WebserverRequestSlaveIdsHandler := func(w http.ResponseWriter, r *http.Request) {
		WebserverRequestSlaveIds(w, r, slaveMap)
	}

	testServer := httptest.NewServer(http.HandlerFunc(WebserverRequestSlaveIdsHandler))

	client := &http.Client{}
	form := url.Values{}
	form.Set("message", "send_me_the_list")
	resp, err := client.PostForm(testServer.URL, form)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, nil, err)
	form.Set("message", "wrong_message")
	resp, err = client.PostForm(testServer.URL, form)
	assert.Equal(t, 500, resp.StatusCode)
}
