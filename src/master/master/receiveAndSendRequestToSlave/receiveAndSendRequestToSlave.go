package receiveAndSendRequestToSlave

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"master/master"
	"net/http"
	"network"
)

type PostURLRequest struct {
	DestinationSlaveName string
	URLToLoadInBrowser   string
}

func ReceiveRequestAndSendToSlave(writer http.ResponseWriter, request *http.Request, slaveMap map[string]master.Slave) {
	POSTRequestBody, _ := ioutil.ReadAll(request.Body)
	defer request.Body.Close()

	incomingRequest, _ := parseJSON(POSTRequestBody)
	destinationSlaveAddress := destinationSlaveAddress(incomingRequest.DestinationSlaveName, slaveMap)
	if destinationSlaveAddress == "" {
		fmt.Println("Abandoning request.")
		fmt.Fprintf(writer, "ERROR: Failed to contact slave. Slave has no URL stored.")
		return
	}

	fmt.Printf("\nSending %v to %v at %v\n", incomingRequest.URLToLoadInBrowser, incomingRequest.DestinationSlaveName, destinationSlaveAddress)
	sendURLValueMessageToSlave(destinationSlaveAddress, incomingRequest.URLToLoadInBrowser)
}

func parseJSON(input []byte) (request PostURLRequest, err error) {
	err = json.Unmarshal(input, &request)
	if err != nil {
		fmt.Println("error:", err)
	}
	return request, err
}

func destinationSlaveAddress(slaveName string, slaveMap map[string]master.Slave) (slaveAddress string) {
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

func sendURLValueMessageToSlave(slaveIPAddress string, urlToDisplay string) (err error) {
	client := &http.Client{}
	form := network.CreateFormWithInitialValues(map[string]string{"url": urlToDisplay})
	response, err := client.PostForm(slaveIPAddress, form)
	if err != nil {
		fmt.Printf("Error slave is not responding: %v\n", err)
		return
	}
	body, _ := ioutil.ReadAll(response.Body)

	defer response.Body.Close()
	fmt.Printf("Slave message: %v\n", (string(body[:])))
	return
}
