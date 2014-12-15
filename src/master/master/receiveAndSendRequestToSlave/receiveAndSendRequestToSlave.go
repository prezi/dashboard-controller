package receiveAndSendRequestToSlave

import (
	"fmt"
	"io/ioutil"
	"master/master"
	"net/http"
	"network"
	"strings"
)

type PostURLRequest struct {
	DestinationSlaveName string
	URLToLoadInBrowser   string
}

func ReceiveRequestAndSendToSlave(slaveMap map[string]master.Slave, slaveNames []string, urlToLoadInBrowser string) {
	for _, slaveName := range slaveNames {
		destinationSlaveAddress := getDestinationSlaveAddress(slaveName, slaveMap)
		if destinationSlaveAddress == "" {
			fmt.Println("Abandoning request.")
			return
		}

		fmt.Printf("\nSending %v to %v at %v\n", urlToLoadInBrowser, slaveName, destinationSlaveAddress)
		sendURLValueMessageToSlave(destinationSlaveAddress, urlToLoadInBrowser)
		updateSlaveDisplayedURL(slaveMap, slaveName, urlToLoadInBrowser)
	}
}

func updateSlaveDisplayedURL(slaveMap map[string]master.Slave, slaveName, urlToLoadInBrowser string) {
	slaveInstance := slaveMap[slaveName]
	slaveInstance.PreviouslyDisplayedURL = slaveInstance.DisplayedURL
	slaveInstance.DisplayedURL = urlToLoadInBrowser
	slaveMap[slaveName] = slaveInstance
	fmt.Println(slaveMap[slaveName])
}

func getDestinationSlaveAddress(slaveName string, slaveMap map[string]master.Slave) (slaveAddress string) {
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

func CheckIfRequestedSlavesAreConnected(slaveMap map[string]master.Slave, slaveNamesToUpdate []string) (string) {
	offlineSlaveList := make([]string, 0 ,len(slaveNamesToUpdate))
	for _, nonExistentSlave := range slaveNamesToUpdate {
		if _, isExists := slaveMap[nonExistentSlave]; !isExists {
			offlineSlaveList = append(offlineSlaveList, nonExistentSlave)
		}
	}
	return strings.Join(offlineSlaveList,", ")
}
