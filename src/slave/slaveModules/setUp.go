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
	"regexp"
)

const DEFAULT_LOCALHOST_PORT = 8080
const DEFAULT_MASTER_IP_ADDRESS = "localhost:5000" // TODO: Fix TCP connection for local IP. For now, can only communicate with master if same localhost as slave.
const DEFAULT_SLAVE_NAME = "SLAVE NAME UNSPECIFIED"

var err error

func SetUp() (port int) {
	var slaveName string
	var masterIP string
	flag.IntVar(&port, "port", DEFAULT_LOCALHOST_PORT, "the port to listen on for commands")
	flag.StringVar(&slaveName, "slaveName", DEFAULT_SLAVE_NAME, "slave name")
	flag.StringVar(&masterIP, "masterIP", DEFAULT_MASTER_IP_ADDRESS, "master IP address")
	flag.Parse()

	// :0.0 indicates the first screen attached to the first display in localhost
	err = os.Setenv("DISPLAY",":0.0")
	if err != nil {
		fmt.Printf("Error setting DISPLAY environment variable: %v\n", err)
	}

	slaveIPAddress := getIPAddressFromCmdLine(port)
	masterIPAddress := getMasterReceiveSlaveAddress(masterIP) // TODO: make this dynamic
	fmt.Println("THIS IS THE MASTER IP ADDRESS", masterIPAddress)
	sendIPAddressToMaster(slaveName, slaveIPAddress, masterIPAddress)

	fmt.Printf("Listening on port: %v\n", port)
	fmt.Println("You can send HTTP POST requests through the command-line with a 'url' parameter to open the url in a browser.")
	fmt.Printf("e.g.: curl localhost:%v -X POST -d \"url=http://www.google.com\"\n", port)

	return port
}

func GetOS() (OS string) {
	operatingSystemBytes, err := exec.Command("uname", "-a").Output() // display operating system name...why do we need the -a?
	operatingSystemName := string(operatingSystemBytes)

	var kernel string
	// fmt.Println("cmd", operatingSystemName)

	if err != nil {
		fmt.Printf("Error encountered while reading kernal: %v\n", err)
		kernel = "Unknown"
	} else {
		kernel = strings.Split(operatingSystemName, " ")[0]
	}
	fmt.Println("Kernal detected: ", kernel)

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

func getIPAddressFromCmdLine(port int) (IPAddress string) {
	cmd := exec.Command("ifconfig")
	IPAddressBytes, _ := cmd.Output()
	IPAddress = string(IPAddressBytes)
	inetAddressRegexpPattern := "inet (addr:)?([0-9]*\\.){3}[0-9]*"
	re := regexp.MustCompile(inetAddressRegexpPattern)
	IPAddress = re.FindAllString(IPAddress, -1)[1]
	IPAddress = strings.Split(IPAddress, " ")[1]

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
	fmt.Println("slaveIPAddress: ", slaveIPAddress)

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