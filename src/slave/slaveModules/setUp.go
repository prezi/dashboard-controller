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
const DEFAULT_MASTER_IP_ADDRESS = "http://localhost:5000" // can also receive this from user input
const DEFAULT_SLAVE_NAME = "SLAVE NAME UNSPECIFIED" // will need to receive this back from the master, or can be user-specified name

var port int
var OS string
var slaveName string
var err error

func SetUp() (port int){
	setOS()
	if (OS == "Unknown") {
		fmt.Println("ERROR: Failed to detect operating system.")
		fmt.Println("Aborting program.")
		os.Exit(1)
	} else {
		fmt.Printf("Operating system detected: %v\n", OS)
	}
	// can pass flag argument: $ ./slave -port=8080
	// if flag not specified, will set port=DEFAULT_LOCALHOST_PORT
	flag.IntVar(&port, "port", DEFAULT_LOCALHOST_PORT, "the port to listen on for commands")
	// can pass flag argument: $ ./slave -slaveName="Slave Name"
	// if flag not specified, will set port=DEFAULT_SLAVE_NAME
	flag.StringVar(&slaveName, "slaveName", DEFAULT_SLAVE_NAME, "slave name")
	flag.Parse()

	// :0.0 indicates the first screen attached to the first display in localhost
	err = os.Setenv("DISPLAY",":0.0")
	if err != nil {
		fmt.Printf("Error setting DISPLAY environment variable: %v\n", err)
	}
	sendIPAddressToMaster()

	fmt.Printf("Listening on port: %v\n", port)
	fmt.Println("You can send HTTP POST requests through the command-line with a 'url' parameter to open the url in a browser.")
	fmt.Printf("e.g.: curl localhost:%v -X POST -d \"url=http://www.google.com\"\n", port)

	return port
}

func setOS() {
	// func (c *Cmd) Output() ([]byte, error)
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
}

func getIPAddressFromCmdLine() (IPAddress string){
	cmd := exec.Command("ifconfig")
	IPAddressBytes, _ := cmd.Output()
	IPAddress = string(IPAddressBytes)
	inetAddressRegexpPattern := "inet (addr:)?([0-9]*\\.){3}[0-9]*"
	re := regexp.MustCompile(inetAddressRegexpPattern)
	IPAddress = re.FindAllString(IPAddress, -1)[1]
	IPAddress = strings.Split(IPAddress, " ")[1]

	return IPAddress
}

func sendIPAddressToMaster() {
	client := &http.Client{}
	slaveIPAddress := getIPAddressFromCmdLine()
	form := url.Values{}
	slaveAddressArray := []string{"http://", slaveIPAddress,":", strconv.Itoa(port)}
	slaveIPAddress = strings.Join(slaveAddressArray, "")
	form.Set("slaveName", slaveName)
	form.Set("slaveIPAddress", slaveIPAddress)
	fmt.Println("slaveIPAddress: ", slaveIPAddress)

	masterIPAddressAndExtentionArray := []string{DEFAULT_MASTER_IP_ADDRESS, "/receive_slave"} 
	masterReceiveSlaveAddress := strings.Join(masterIPAddressAndExtentionArray, "")

	_, err := client.PostForm(masterReceiveSlaveAddress, form)

	if err != nil {
		fmt.Printf("Error communicating with master: %v\n", err)
		fmt.Println("Aborting program.")
		// os.Exit(1)
	}

	fmt.Printf("Slave mapped to master at %v.\n", DEFAULT_MASTER_IP_ADDRESS)
	fmt.Printf("Slave Name: %v.\n", slaveName)
	if slaveName == DEFAULT_SLAVE_NAME {
		fmt.Println("TIP: Specify slave name at startup with the flag '-slaveName'") 
		fmt.Println("eg. -slaveName=\"Main Lobby\"")
	}
	fmt.Printf("\n\n")
}