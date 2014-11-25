package network

import (
	"net"
	"net/url"
	"net/http"
	"fmt"
	"os"
	"flag"
	"strings"
	"os/exec"
	"regexp"
)

const (
	DEFAULT_SLAVE_NAME = "SLAVE NAME UNSPECIFIED"
	DEFAULT_MASTER_IP_ADDRESS = "localhost"
	DEFAULT_MASTER_PORT = "5000"
)

func GetLocalIPAddress(port string) (IPAddress string) {
	IPAddressBytes := getIPAddressBytesFromCmdLine()

	return addProtocolAndPortToIp(parseIpAddress(IPAddressBytes), port)
}

func getIPAddressBytesFromCmdLine() (IPAddressWithNoise string) {
	cmd := exec.Command("ifconfig")
	IPAddressBytes, _ := cmd.Output()
	IPAddressWithNoise = string(IPAddressBytes)
	return
}

func parseIpAddress(IPAddress string) string {
	inetAddressRegexpPattern := "inet (addr:)?([0-9]*\\.){3}[0-9]*"
	re := regexp.MustCompile(inetAddressRegexpPattern)
	IPAddress = re.FindAllString(IPAddress, -1)[1]
	IPAddress = strings.Split(IPAddress, " ")[1]
	return IPAddress
}

func addProtocolAndPortToIp(IPAddress string, port string) string {
	hostIPWithPort := net.JoinHostPort(IPAddress, port)
	protocolWithHostIPAndPort := []string{"http://", hostIPWithPort}
	url := strings.Join(protocolWithHostIPAndPort, "")
	return url
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

func AddProtocolAndPortToIP(IPAddress, port string) (url string) {
	hostIPWithPort := net.JoinHostPort(IPAddress, port)
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

func SetMasterUrl() (masterUrl string) {
	masterIP:= DEFAULT_MASTER_IP_ADDRESS
	masterPort:= DEFAULT_MASTER_PORT
	flag.StringVar(&masterIP, "masterIP", DEFAULT_MASTER_IP_ADDRESS, "master IP address")
	flag.StringVar(&masterPort, "masterPort", DEFAULT_MASTER_PORT, "master port number")
	flag.Parse()
	return AddProtocolAndPortToIP(masterIP, masterPort)
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
