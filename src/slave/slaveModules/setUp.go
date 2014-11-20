package slaveModule

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"net/url"
	"net"
)

const DEFAULT_LOCALHOST_PORT = 8080
const DEFAULT_MASTER_IP_ADDRESS = "localhost:5000" 
const DEFAULT_SLAVE_NAME = "SLAVE NAME UNSPECIFIED"

var err error

func SetUp() (port int, slaveName, masterIP, OS string) {
	port, slaveName, masterIP = configFlags()
	OS = getOS()
	// :0.0 indicates the first screen attached to the first display in localhost
	err = os.Setenv("DISPLAY",":0.0")
	if err != nil {
		fmt.Printf("Error setting DISPLAY environment variable: %v\n", err)
	}

	slaveIPAddress := getIPAddress(port)
	masterIPAddressToReceiveSlave := getMasterReceiveSlaveAddress(masterIP)
	sendIPAddressToMaster(slaveName, slaveIPAddress, masterIPAddressToReceiveSlave)

	fmt.Printf("Listening on port: %v\n", port)
	fmt.Println("You can send HTTP POST requests through the command-line with a 'url' parameter to open the url in a browser.")
	fmt.Printf("e.g.: curl localhost:%v -X POST -d \"url=http://www.google.com\"\n", port)

	return port, slaveName, masterIP, OS
}

func configFlags() (port int, slaveName, masterIP string) {
	flag.IntVar(&port, "port", DEFAULT_LOCALHOST_PORT, "the port to listen on for commands")
	flag.StringVar(&slaveName, "slaveName", DEFAULT_SLAVE_NAME, "slave name")
	flag.StringVar(&masterIP, "masterIP", DEFAULT_MASTER_IP_ADDRESS, "master IP address")
	flag.Parse()
	return port, slaveName, masterIP
}

func getOS() (OS string) {
	operatingSystemBytes, err := exec.Command("uname", "-a").Output() // display operating system name...why do we need the -a?
	operatingSystemName := string(operatingSystemBytes)

	var kernel string
	// fmt.Println("cmd", operatingSystemName)

	if err != nil {
		fmt.Printf("Error encountered while reading kernel: %v\n", err)
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

func getIPAddress(port int) (IPAddress string) {
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

	IPAddress = IPAddressArray[0]
	slaveAddressArray := []string{"http://", IPAddress,":", strconv.Itoa(port)}
	IPAddress = strings.Join(slaveAddressArray, "")
	
	return IPAddress
}

func getMasterReceiveSlaveAddress(masterIPAddress string) (masterAddress string) {
	masterIPAddressAndExtentionArray := []string{"http://", masterIPAddress, "/receive_slave"} 
	return strings.Join(masterIPAddressAndExtentionArray, "")
}

func sendIPAddressToMaster(slaveName string, slaveIPAddress string, masterAddress string) {
	client := &http.Client{}
	form := url.Values{}
	form.Set("slaveName", slaveName)
	form.Set("slaveIPAddress", slaveIPAddress)
	// fmt.Println("slaveIPAddress: ", slaveIPAddress)

	_, err := client.PostForm(masterAddress, form)

	if err != nil {
		fmt.Printf("Error communicating with master: %v\n", err)
		fmt.Println("Aborting program.")
		// os.Exit(1)
	}

	fmt.Printf("Slave mapped to master at %v.\n", masterAddress)
	fmt.Printf("Slave Name: %v.\n", slaveName)
	if slaveName == DEFAULT_SLAVE_NAME {
		fmt.Println("TIP: Specify slave name at startup with the flag '-slaveName'") 
		fmt.Println("eg. -slaveName=\"Main Lobby\"")
	}
	fmt.Printf("\n\n")
}