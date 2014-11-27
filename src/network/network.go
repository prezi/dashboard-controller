package network

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

const (
	DEFAULT_SLAVE_NAME        = "SLAVE NAME UNSPECIFIED"
	DEFAULT_MASTER_IP_ADDRESS = "localhost"
	DEFAULT_MASTER_PORT       = "5000"
)

func AddProtocolAndPortToIP(IPAddress, port string) (url string) {
	hostIPWithPort := net.JoinHostPort(IPAddress, port)
	return "http://" + hostIPWithPort
}

func GetOS() (OS string) {
	operatingSystemBytes, err := exec.Command("uname", "-a").Output() // display operating system name...why do we need the -a?
	operatingSystemName := string(operatingSystemBytes)

	var kernel string
	// fmt.Println("cmd", operatingSystemName)

	if ErrorHandler(err, "Error encountered while reading kernel: %v\n") {
		kernel = "Unknown"
	} else {
		kernel = strings.Split(operatingSystemName, " ")[0]
	}
	fmt.Println("\nKernel detected: ", kernel)

	switch kernel {
	case "Linux":
		OS = "Linux"
	case "Darwin":
		OS = "OS X"
	default:
		OS = "Unknown"
	}

	if OS == "Unknown" {
		fmt.Println("ERROR: Failed to detect operating system.")
		fmt.Println("Aborting program.")
		os.Exit(1)
	} else {
		fmt.Printf("Operating system detected: %v\n", OS)
	}
	return OS
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
