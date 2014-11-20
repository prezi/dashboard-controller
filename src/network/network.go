package network

import (
	"strings"
	"regexp"
	"os/exec"
	"net"
)

func getIPAddress(port string) (IPAddress string) {
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
