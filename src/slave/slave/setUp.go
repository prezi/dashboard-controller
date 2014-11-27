package slave

import (
	"flag"
	"fmt"
	"network"

	"os/exec"
)

const (
	DEFAULT_SLAVE_NAME        = "SLAVE NAME UNSPECIFIED"
	DEFAULT_LOCALHOST_PORT    = "8080"
	DEFAULT_MASTER_IP_ADDRESS = "localhost"
	DEFAULT_MASTER_PORT       = "5000"
)

var err error

func SetUp() (port, slaveName, masterURL, OS string, BrowserProcess *exec.Cmd) {
	port, slaveName, masterIP, masterPort := configFlags()
	masterURL = network.AddProtocolAndPortToIP(masterIP, masterPort)
	OS = network.GetOS()
	// :0.0 indicates the first screen attached to the first display in localhost
	err = os.Setenv("DISPLAY", ":0.0")
	if err != nil {
		fmt.Printf("Error setting DISPLAY environment variable: %v\n", err)
	}

	slaveIPAddress := network.GetLocalIPAddress(port)
	slaveURL := network.AddProtocolAndPortToIP(slaveIPAddress, port)

	masterURLToReceiveSlave := masterURL + "/receive_slave"
	fmt.Print(slaveName, slaveURL, masterURLToReceiveSlave)

	fmt.Printf("Listening on port: %v\n", port)
	fmt.Println("You can send HTTP POST requests through the command-line with a 'url' parameter to open the url in a browser.")
	fmt.Printf("e.g.: curl localhost:%v -X POST -d \"url=http://www.google.com\"\n", port)
	BrowserProcess = nil

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
