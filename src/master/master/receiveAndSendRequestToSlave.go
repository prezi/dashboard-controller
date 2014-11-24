package master

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"net/url"
)

type Slave struct {
	ID string
	URL string
}

func ReceiveRequestAndSendToSlave(_ http.ResponseWriter, request *http.Request) {
	POSTRequestBody, _ := ioutil.ReadAll(request.Body)
	defer request.Body.Close()

	slave, _ := parseJson(POSTRequestBody)
	destinationSlaveAddress := destinationSlaveAddress(slave.ID)
	if destinationSlaveAddress == "" {
		fmt.Println("Abandoning request.")
		return
	}

	fmt.Printf("\nSending %v to %v at %v", slave.URL, slave.ID, destinationSlaveAddress)
	sendUrlValueMessageToSlave(destinationSlaveAddress, slave.URL)
}

func parseJson(input []byte) (slave Slave, err error) {
	err = json.Unmarshal(input, &slave)
	if err != nil {
		fmt.Println("error:", err)
	}
	return slave, err
}

func destinationSlaveAddress(slaveID string) (slaveAddress string) {
	if len(slaveIPMap) == 0 {
		fmt.Println("ERROR: No slaves available.")
		return
	}

	slaveAddress = slaveIPMap[slaveID]
	if slaveAddress ==  "" {
		fmt.Printf("ERROR: \"%v\" is not a valid slave ID.\n", slaveID)
		fmt.Println("Valid slave IDs are: ", slaveIPMap)
		return
	}
	return slaveAddress
}

func sendUrlValueMessageToSlave(slaveIPAddress string, urlToDisplay string) {
	client := &http.Client{}

	form := url.Values{}
	form.Set("url", urlToDisplay)

	_,_ = client.PostForm(slaveIPAddress, form)
}