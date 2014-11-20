package network

import (
	"strings"
	"regexp"
	"os/exec"
	"net"
	"net/url"
	"net/http"
	"fmt"
)

const DEFAULT_SLAVE_NAME = "SLAVE NAME UNSPECIFIED"

func GetUrl(port string) string {
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
	url := "http://"
	url += hostIPWithPort
	return url
}

func GetMasterUrl(masterIPAddress string) (masterAddress string) {
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
