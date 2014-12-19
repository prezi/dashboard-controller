package slave

import (
	"flag"
	"fmt"
	"os"
	"net"
	"network"
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
	OS = network.GetOS()
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
