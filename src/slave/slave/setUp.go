package slave

import (
	"flag"
	"fmt"
	"network"
	"os"
	"net"
	"os/exec"
	"strings"
)

const (
	DEFAULT_SLAVE_NAME        = "SLAVE NAME UNSPECIFIED"
	DEFAULT_LOCALHOST_PORT    = "8080"
	DEFAULT_MASTER_IP_ADDRESS = "localhost"
	DEFAULT_MASTER_PORT       = "5000"
)

var err error

func SetUp() (port, slaveName, masterURL, OS string) {
	port, slaveName, masterIP, masterPort := configFlags()
	masterURL = addProtocolAndPortToIP(masterIP, masterPort)
	OS = getOS()
	// :0.0 indicates the first screen attached to the first display in localhost
	err = os.Setenv("DISPLAY", ":0.0")
	if err != nil {
		fmt.Printf("Error setting DISPLAY environment variable: %v\n", err)
	}

	fmt.Printf("\nListening on port: %v\n", port)
	fmt.Println("You can send HTTP POST requests through the command-line with a 'url' parameter to open the url in a browser.")
	fmt.Printf("e.g.: curl localhost:%v -X POST -d \"url=http://www.google.com\"\n", port)

	return
}

func configFlags() (port, slaveName, masterIP, masterPort string) {
	flag.StringVar(&port, "port", DEFAULT_LOCALHOST_PORT, "the port to listen on for commands")
	flag.StringVar(&slaveName, "slaveName", DEFAULT_SLAVE_NAME, "slave name")
	flag.StringVar(&masterIP, "masterIP", DEFAULT_MASTER_IP_ADDRESS, "master IP address")
	flag.StringVar(&masterPort, "masterPort", DEFAULT_MASTER_PORT, "master port number")
	flag.Parse()
	return port, slaveName, masterIP, masterPort
}

func addProtocolAndPortToIP(IPAddress, port string) (url string) {
	hostIPWithPort := net.JoinHostPort(IPAddress, port)
	return "http://" + hostIPWithPort
}

func getOS() (OS string) {
	operatingSystemBytes, err := exec.Command("uname", "-a").Output() // display operating system name...why do we need the -a?
	operatingSystemName := string(operatingSystemBytes)

	var kernel string

	if network.ErrorHandler(err, "Error encountered while reading kernel: %v\n") {
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
