package slave

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"net/url"
	"network"
)

const (
	DEFAULT_SLAVE_NAME = "SLAVE NAME UNSPECIFIED"
	DEFAULT_LOCALHOST_PORT = 8080
	DEFAULT_MASTER_IP_ADDRESS = "localhost"
	DEFAULT_MASTER_PORT = 5000
)

var err error

func SetUp() (port int, slaveName, masterURL, OS string) {
	port, slaveName, masterIP, masterPort := configFlags()
	masterURL = network.AddProtocolAndPortToIP(masterIP, masterPort)
	OS = getOS()
	// :0.0 indicates the first screen attached to the first display in localhost
	err = os.Setenv("DISPLAY",":0.0")
	if err != nil {
		fmt.Printf("Error setting DISPLAY environment variable: %v\n", err)
	}

	slaveIPAddress := network.GetLocalIPAddress()
	slaveURL := network.AddProtocolAndPortToIP(slaveIPAddress, port)

	masterURLToReceiveSlave := masterURL + "/receive_slave"
	sendSlaveURLToMaster(slaveName, slaveURL, masterURLToReceiveSlave)

	fmt.Printf("Listening on port: %v\n", port)
	fmt.Println("You can send HTTP POST requests through the command-line with a 'url' parameter to open the url in a browser.")
	fmt.Printf("e.g.: curl localhost:%v -X POST -d \"url=http://www.google.com\"\n", port)

	return port, slaveName, masterURL, OS
}

func configFlags() (port int, slaveName, masterIP string, masterPort int) {
	flag.IntVar(&port, "port", DEFAULT_LOCALHOST_PORT, "the port to listen on for commands")
	flag.StringVar(&slaveName, "slaveName", DEFAULT_SLAVE_NAME, "slave name")
	flag.StringVar(&masterIP, "masterIP", DEFAULT_MASTER_IP_ADDRESS, "master IP address")
	flag.IntVar(&masterPort, "masterPort", DEFAULT_MASTER_PORT, "master port number")
	flag.Parse()
	return port, slaveName, masterIP, masterPort
}

func getOS() (OS string) {
	operatingSystemBytes, err := exec.Command("uname", "-a").Output() // display operating system name...why do we need the -a?
	operatingSystemName := string(operatingSystemBytes)

	var kernel string
	// fmt.Println("cmd", operatingSystemName)

	if network.ErrorHandler(err, "Error encountered while reading kernel: %v\n") {
		kernel = "Unknown"
	} else {
		kernel = strings.Split(operatingSystemName, " ")[0]
	}
	fmt.Println("Kernel detected: ", kernel)

	switch kernel {
	case "Linux":
		OS = "Linux"
	case "Darwin":
		OS = "OS X"
	default:
		OS = "Unknown"
	}

	if (OS == "Unknown") {
		fmt.Println("ERROR: Failed to detect operating system.")
		fmt.Println("Aborting program.")
		os.Exit(1)
	} else {
		fmt.Printf("Operating system detected: %v\n", OS)
	}
	return OS
}

func sendSlaveURLToMaster(slaveName, slaveURL, masterURL string) {
	client := &http.Client{}
	form := url.Values{}
	form.Set("slaveName", slaveName)
	form.Set("slaveURL", slaveURL)
	fmt.Println("slaveURL: ", slaveURL)

	_, err := client.PostForm(masterURL, form)

	network.ErrorHandler(err, "Error communicating with master: %v\n")

	fmt.Printf("Slave mapped to master at %v.\n", masterURL)
	fmt.Printf("Slave Name: %v.\n", slaveName)
	if slaveName == DEFAULT_SLAVE_NAME {
		fmt.Println("TIP: Specify slave name at startup with the flag '-slaveName'")
		fmt.Println("eg. -slaveName=\"Main Lobby\"")
	}
	fmt.Printf("\n\n")
}
