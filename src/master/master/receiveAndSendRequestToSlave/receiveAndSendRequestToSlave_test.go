package receiveAndSendRequestToSlave

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"master/master"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const (
	TEST_SLAVE_NAME = "test slave"
	TEST_URL_1      = "http://google.com"
	TEST_URL_2      = "http://placekitten.com"
)

func InitializeTestSlaveMap() (slaveMap map[string]master.Slave) {
	slaveMap = make(map[string]master.Slave)
	slaveMap["slave1"] = master.Slave{URL: "http://10.0.0.122:8080", Heartbeat: time.Now()}
	slaveMap["slave2"] = master.Slave{URL: "http://10.0.1.11:8080", Heartbeat: time.Now()}
	return slaveMap
}

func TestReceiveRequestAndSendToSlave(t *testing.T) {
	testSlaveMap := make(map[string]master.Slave)
	var receivedUrl string
	testMaster := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ReceiveRequestAndSendToSlave(testSlaveMap, TEST_SLAVE_NAME, TEST_URL_1)
	}))

	testSlave := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		receivedUrl = request.PostFormValue("url")
	}))
	testSlaveMap[TEST_SLAVE_NAME] = master.Slave{testSlave.URL, time.Now(), TEST_URL_1, TEST_URL_1}

	m := PostURLRequest{TEST_SLAVE_NAME, TEST_URL_1}
	json_message, _ := json.Marshal(m)
	client := &http.Client{}
	_, err := client.Post(testMaster.URL, "application/json", strings.NewReader(string(json_message)))

	assert.Equal(t, TEST_URL_1, receivedUrl)
	assert.Nil(t, err)
}

func TestUpdateSlaveDisplayedURL(t *testing.T) {
	testSlaveMap := InitializeTestSlaveMap()
	updateSlaveDisplayedURL(testSlaveMap, "slave1", TEST_URL_2)
	assert.Equal(t, testSlaveMap["slave1"].DisplayedURL, TEST_URL_2)
	assert.Equal(t, testSlaveMap["slave1"].PreviouslyDisplayedURL, "")
	assert.Equal(t, testSlaveMap["slave2"].DisplayedURL, "")
	assert.Equal(t, testSlaveMap["slave2"].PreviouslyDisplayedURL, "")
}

func TestReceiveRequestAndSendToSlaveWithEmptySlaveAddress(t *testing.T) {
	testSlaveMap := make(map[string]master.Slave)

	testMaster := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ReceiveRequestAndSendToSlave(testSlaveMap, "testSlaveName", "someurl")
	}))

	testSlaveMap[TEST_SLAVE_NAME] = master.Slave{"", time.Now(), TEST_URL_1, TEST_URL_1}

	m := PostURLRequest{TEST_SLAVE_NAME, TEST_URL_1}
	json_message, _ := json.Marshal(m)
	client := &http.Client{}
	_, err := client.Post(testMaster.URL, "application/json", strings.NewReader(string(json_message)))
	// body, err := ioutil.ReadAll(response.Body)
	// defer response.Body.Close()
	// receivedResponse := string(body[:])
	// assert.Equal(t, "ERROR: Failed to contact slave. Slave has no URL stored.", receivedResponse)
	assert.Nil(t, err)
}

func TestGetDestinationAddressSlave(t *testing.T) {
	slaveMap := InitializeTestSlaveMap()
	destinationURL := getDestinationSlaveAddress("slave1", slaveMap)

	assert.Equal(t, "http://10.0.0.122:8080", destinationURL)
}

func TestDestinationAddressSlaveForEmptySlaveMap(t *testing.T) {
	slaveMap := make(map[string]master.Slave)
	destinationURL := getDestinationSlaveAddress("slave2", slaveMap)

	assert.Equal(t, "", destinationURL)
}

func TestSendURLValueMessageToSlave(t *testing.T) {
	var numberOfMessagesSent = 0
	var url = ""

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		numberOfMessagesSent++
		url = request.PostFormValue("url")
	}))

	err := sendURLValueMessageToSlave(testServer.URL, "http://index.hu")

	assert.Equal(t, 1, numberOfMessagesSent)
	assert.Equal(t, "http://index.hu", url)
	assert.Nil(t, err)
}

func TestSendURLValueMessageToSlaveSlaveDoesNotRespond(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
	}))
	testServer.Close()
	err := sendURLValueMessageToSlave(testServer.URL, "http://index.hu")
	assert.NotNil(t, err)
}

func TestCheckIfRequestedSlavesAreConnected(t *testing.T) {
	slaveList := []string{"a","b","c"}
	slaveMap := map[string]master.Slave{
		"a":master.Slave{},
		"b":master.Slave{},
		"c":master.Slave{},
	}
	returnValue := CheckIfRequestedSlavesAreConnected(slaveMap, slaveList)
	assert.Equal(t, returnValue, "")
}

func TestCheckIfRequestedSlavesAreConnectedWithOneMissingSlave(t *testing.T) {
	slaveList := []string{"a","b","c","d"}
	slaveMap := map[string]master.Slave{
		"a":master.Slave{},
		"b":master.Slave{},
		"c":master.Slave{},
	}
	returnValue := CheckIfRequestedSlavesAreConnected(slaveMap, slaveList)
	assert.Equal(t, returnValue, "d")
}
func TestCheckIfRequestedSlavesAreConnectedWithMultipleMissingSlaves(t *testing.T) {
	slaveList := []string{"a","b","c","d"}
	slaveMap := map[string]master.Slave{
		"a":master.Slave{},
		"b":master.Slave{},
	}
	returnValue := CheckIfRequestedSlavesAreConnected(slaveMap, slaveList)
	assert.Equal(t, returnValue, "c, d")
}
