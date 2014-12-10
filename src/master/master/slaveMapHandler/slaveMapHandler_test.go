package slaveMapHandler

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"master/master"
	"net/http/httptest"
	"net/http"
	"io/ioutil"
	"time"
)

func TestInitiateNonEmptySlaveMapHandler(t *testing.T) {
	slaveMap := make(map[string]master.Slave)
	slaveMap["slave1"] = master.Slave{URL: "http://10.0.0.122:8080", Heartbeat: time.Now(), PreviouslyDisplayedURL: "http://www.google.com", DisplayedURL: "http://www.prezi.com"}

	testMaster := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		InitiateSlaveMapHandler(slaveMap)
	}))
	client := &http.Client{}
	response, err := client.Get(testMaster.URL)

//	var slaveMapFromJson []string
	slaveNames, _ := ioutil.ReadAll(response.Body)
//	err = json.Unmarshal(slaveNames, &slaveMapFromJson)
//	if err != nil {
//		fmt.Println(err)
//	}

	testMaster.CloseClientConnections()
	testMaster.Close()

	assert.Equal(t, []string{"slave1"}, slaveNames)
	assert.Nil(t, err)
}
