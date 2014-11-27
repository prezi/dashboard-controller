package master

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type PostURLRequest struct {
	DestinationSlaveName string
	URLToLoadInBrowser   string
}

func ReceiveRequestAndSendToSlave(_ http.ResponseWriter, request *http.Request, slaveMap map[string]Slave) {
	POSTRequestBody, _ := ioutil.ReadAll(request.Body)
	defer request.Body.Close()

	slave, _ := parseJson(POSTRequestBody)
	destinationSlaveAddress := destinationSlaveAddress(slave.DestinationSlaveName, slaveMap)
	if destinationSlaveAddress == "" {
		fmt.Println("Abandoning request.")
		return
	}

	fmt.Printf("\nSending %v to %v at %v\n", slave.URLToLoadInBrowser, slave.DestinationSlaveName, destinationSlaveAddress)
	sendUrlValueMessageToSlave(destinationSlaveAddress, slave.URLToLoadInBrowser)
}

func parseJson(input []byte) (request PostURLRequest, err error) {
	err = json.Unmarshal(input, &request)
	if err != nil {
		fmt.Println("error:", err)
	}
	return request, err
}

func destinationSlaveAddress(slaveName string, slaveMap map[string]Slave) (slaveAddress string) {
	if len(slaveMap) == 0 {
		fmt.Println("ERROR: No slaves available.")
		return
	}

	slaveAddress = slaveMap[slaveName].URL
	if slaveAddress == "" {
		fmt.Printf("ERROR: \"%v\" is not a valid slave name.\n", slaveName)
		fmt.Println("Valid slave names are: ", slaveMap)
		return
	}
	return slaveAddress
}

func sendUrlValueMessageToSlave(slaveIPAddress string, urlToDisplay string) {
	client := &http.Client{}

	form := url.Values{}
	form.Set("url", urlToDisplay)

	response, err := client.PostForm(slaveIPAddress, form)
	if err != nil {
		fmt.Printf("Error slave is not responding: %v\n", err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading slave response: %v\n", err)
		return
	}
	defer response.Body.Close()
	fmt.Printf("Slave message: %v\n", (string(body)))
}
