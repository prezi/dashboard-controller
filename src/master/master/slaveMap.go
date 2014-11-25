package master

import (
	"net/http"
	"fmt"
	"strings"
	"time"
	"encoding/json"
)

var webserverAddress = "http://localhost:4003"// TODO: make dynamic webserver address


type Slave struct {
	URL string
	heartbeat time.Time
	displayedURL string // TODO: store currently displayed URL for each slave
}

type IdList struct {
	Id []string
}

// TODO: Allow user to input names of currently running slaves at master startup.
// Alternatively, allow to manually add names of currently running slave while the master is running.
func SetUp() (slaveMap map[string]Slave) {
	return initializeSlaveMap()  
}

func initializeSlaveMap() (slaveMap map[string]Slave) {
	slaveMap = make(map[string]Slave)
	slaveMap["slave1"] = Slave{URL: "http://10.0.0.122:8080", heartbeat: time.Now()}
	slaveMap["slave2"] = Slave{URL: "http://10.0.1.11:8080", heartbeat: time.Now()}
	return slaveMap
}

func ReceiveAndMapSlaveAddress(_ http.ResponseWriter, request *http.Request, slaveMap map[string]Slave) {
	slaveName := request.PostFormValue("slaveName")
	slaveURL := request.PostFormValue("slaveURL")
	fmt.Printf("\nNEW SLAVE RECEIVED.\n")
	fmt.Println("Slave Name: ", slaveName)
	fmt.Println("Slave URL: ", slaveURL)

	if returnedSlave, existsInMap := slaveMap[slaveName]; existsInMap {
		fmt.Printf("WARNING: Slave with name \"%v\" already exists with the IP address: %v. \nUpdating %v's IP address to %v.\n", slaveName, returnedSlave.URL, slaveName, slaveURL)
	}
	slaveMap[slaveName] = Slave{URL: slaveURL, heartbeat: time.Now()}
	err := sendSlaveToWebserver(webserverAddress, slaveMap)
	printServerResponse(err, slaveName)

	fmt.Printf("Mapped \"%v\" to %v.\n", slaveName, slaveURL)
	fmt.Println("Valid slave IDs are: ")
	for slaveName, _ := range slaveMap {
		fmt.Println(slaveName)
	}
}

func printServerResponse(error error, slaveName string) {
	if error != nil {
		fmt.Printf("Error communicating with webserver: %v\n", error)
		fmt.Printf("%v not updated on webserver.\n", slaveName)
	} else {
		fmt.Printf("Added \"%v\" to webserver slave list.\n", slaveName)
	}
}

func sendSlaveToWebserver(webserverAddress string, slaveMap map[string]Slave) (err error) {
	err = nil
	client := &http.Client{}
	webserverAddress = webserverAddress + "/receive_slave"
	var idList IdList
	for slaveName := range slaveMap {
        idList.Id = append(idList.Id, slaveName)
    }
	jsonMessage, err := json.Marshal(idList)
	_,err = client.Post(webserverAddress, "application/json", strings.NewReader(string(jsonMessage)))
	return err
}

func WebserverRequestSlaveIds(writer http.ResponseWriter, request *http.Request, slaveMap map[string]Slave) {
	message := request.PostFormValue("message")
	if message == "send_me_the_list" {
		sendSlaveToWebserver(webserverAddress, slaveMap)
		writer.WriteHeader(200)
	} else {
		writer.WriteHeader(500)
	}
}
