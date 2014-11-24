package master

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"encoding/json"
	"sort"
	"net/url"
)

func TestInitializeSlaveIPs(t *testing.T) {
	slaveIPMap := initializeSlaveIPs()

	assert.Equal(t, "http://10.0.0.122:8080", slaveIPMap["1"])
	assert.Equal(t, "http://10.0.1.11:8080", slaveIPMap["2"])
}

func TestDestinationAddressSlave1(t *testing.T) {
	slaveIPMap = SetUp()
	destinationURL := destinationSlaveAddress("1")

	assert.Equal(t, "http://10.0.0.122:8080", destinationURL)
}

func TestDestinationAddressSlave2(t *testing.T) {
	slaveIPMap = SetUp()
	destinationURL := destinationSlaveAddress("2")

	assert.Equal(t, "http://10.0.1.11:8080", destinationURL)
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
	slaveIPs := initializeSlaveIPs()
	sendSlaveToWebserver(testServer.URL, slaveIPs)
	validIdList := []string{"1","2"}
	sort.Strings(validIdList)
	sort.Strings(returnedIds)
	assert.Equal(t, returnedIds, validIdList)
}

func TestWebserverRequestSlaveIds(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(WebserverRequestSlaveIds))

	client := &http.Client{}
	form := url.Values{}
	form.Set("message","send_me_the_list")
	resp, err := client.PostForm(testServer.URL,form)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, nil, err)
	form.Set("message","wrong_message")
	resp, err = client.PostForm(testServer.URL,form)
	assert.Equal(t, 500 , resp.StatusCode)
}