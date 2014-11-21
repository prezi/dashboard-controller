package network

import (
	"net"
	"net/url"
	"net/http"
	"fmt"
	"os"
	"strconv"
)

const DEFAULT_SLAVE_NAME = "SLAVE NAME UNSPECIFIED"

func GetLocalIPAddress() (IPAddress string) {
	name, err := os.Hostname()
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}

	IPAddressArray, err := net.LookupHost(name)

	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}

	return IPAddressArray[0]
}

func AddProtocolAndPortToIP(IPAddress string, port int) (url string) {
	hostIPWithPort := net.JoinHostPort(IPAddress, strconv.Itoa(port))
	return "http://" + hostIPWithPort
}

func ErrorHandler(err error, message string) (errorOccurred bool) {
	if err != nil {
		fmt.Printf(message, err)
		fmt.Println("Aborting program.")
		// os.Exit(1)
		return true
	}	
	return false
}

func SetMasterIP() (url string) {
	return ""
}

func sendSlaveURLToMaster(slaveName, slaveURL, masterURL string) {
	client := &http.Client{}
	form := url.Values{}
	form.Set("slaveName", slaveName)
	form.Set("slaveURL", slaveURL)
	fmt.Println("slaveURL: ", slaveURL)

	_, err := client.PostForm(masterURL, form)

	ErrorHandler(err, "Error communicating with master: %v\n")

	fmt.Printf("Slave mapped to master at %v.\n", masterURL)
	fmt.Printf("Slave Name: %v.\n", slaveName)
	if slaveName == DEFAULT_SLAVE_NAME {
		fmt.Println("TIP: Specify slave name at startup with the flag '-slaveName'") 
		fmt.Println("eg. -slaveName=\"Main Lobby\"")
	}
	fmt.Printf("\n\n")
}
